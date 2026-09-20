package main

import (
	"encoding/hex"
	"encoding/json"
	"os"
	"testing"
)

// The corpus and the expected answers live in testdata/. expected.json is
// generated from the reference Python tiktoken by scripts/gen_expected.py, and
// CI reruns that script with --check so the committed answers cannot silently
// drift. These tests therefore pin tc to real tiktoken without needing Python
// on the machine running `go test`.

type corpusCase struct {
	Name string `json:"name"`
	Note string `json:"note"`
	Text string `json:"text"`
	Hex  string `json:"hex"`
}

type corpusFile struct {
	Encoding string       `json:"encoding"`
	Cases    []corpusCase `json:"cases"`
}

type expectedCase struct {
	ValidUTF8  bool   `json:"valid_utf8"`
	DecodedHex string `json:"decoded_hex"`
	Tokens     int    `json:"tokens"`
	IDs        []int  `json:"ids"`
}

type expectedFile struct {
	Encoding        string                  `json:"encoding"`
	TiktokenVersion string                  `json:"tiktoken_version"`
	Cases           map[string]expectedCase `json:"cases"`
}

// bytes returns the raw input for a case: readable cases carry text, byte-level
// cases carry hex.
func (c corpusCase) bytes(t *testing.T) []byte {
	t.Helper()
	if (c.Text == "") == (c.Hex == "") && c.Name != "empty" {
		t.Fatalf("case %q: need exactly one of text or hex", c.Name)
	}
	if c.Hex != "" {
		b, err := hex.DecodeString(c.Hex)
		if err != nil {
			t.Fatalf("case %q: bad hex: %v", c.Name, err)
		}
		return b
	}
	return []byte(c.Text)
}

func loadCorpus(t *testing.T) (corpusFile, expectedFile) {
	t.Helper()
	var corpus corpusFile
	var expected expectedFile
	for _, l := range []struct {
		path string
		into any
	}{
		{"testdata/corpus.json", &corpus},
		{"testdata/expected.json", &expected},
	} {
		data, err := os.ReadFile(l.path)
		if err != nil {
			t.Fatalf("read %s: %v", l.path, err)
		}
		if err := json.Unmarshal(data, l.into); err != nil {
			t.Fatalf("parse %s: %v", l.path, err)
		}
	}
	if corpus.Encoding != expected.Encoding {
		t.Fatalf("encoding mismatch: corpus %q, expected %q", corpus.Encoding, expected.Encoding)
	}
	if len(corpus.Cases) != len(expected.Cases) {
		t.Fatalf("corpus has %d cases but expected.json has %d; rerun scripts/gen_expected.py",
			len(corpus.Cases), len(expected.Cases))
	}
	return corpus, expected
}

// TestConformanceTokenCounts is the contract: for every case in the corpus, tc
// must report the same token count the OpenAI tokenizer does.
func TestConformanceTokenCounts(t *testing.T) {
	corpus, expected := loadCorpus(t)
	for _, c := range corpus.Cases {
		c := c
		t.Run(c.Name, func(t *testing.T) {
			want, ok := expected.Cases[c.Name]
			if !ok {
				t.Fatalf("no expected entry; rerun scripts/gen_expected.py")
			}
			raw := c.bytes(t)

			text, valid := decodeUTF8(raw)
			if valid != want.ValidUTF8 {
				t.Errorf("valid UTF-8 = %v, want %v (%s)", valid, want.ValidUTF8, c.Note)
			}
			// Lock the decoder itself to Python's replacement behaviour, not
			// just the resulting count: a wrong number of U+FFFD can still
			// land on the right count by luck.
			if got := hex.EncodeToString([]byte(text)); got != want.DecodedHex {
				t.Errorf("decoded bytes mismatch (%s)\n got: %s\nwant: %s", c.Note, got, want.DecodedHex)
			}

			n, err := countTokens(text)
			if err != nil {
				t.Fatalf("countTokens: %v", err)
			}
			if n != want.Tokens {
				t.Errorf("tokens = %d, want %d (%s)", n, want.Tokens, c.Note)
			}
		})
	}
}

// TestConformanceViaCLI runs the same corpus through the real CLI entry point,
// so the wiring between decoding, counting and output is covered too.
func TestConformanceViaCLI(t *testing.T) {
	corpus, expected := loadCorpus(t)
	dir := t.TempDir()
	for _, c := range corpus.Cases {
		c := c
		t.Run(c.Name, func(t *testing.T) {
			want := expected.Cases[c.Name]
			path := dir + "/" + c.Name
			if err := os.WriteFile(path, c.bytes(t), 0o644); err != nil {
				t.Fatal(err)
			}
			n, valid, err := countFile(path)
			if err != nil {
				t.Fatalf("countFile: %v", err)
			}
			if n != want.Tokens {
				t.Errorf("tokens = %d, want %d (%s)", n, want.Tokens, c.Note)
			}
			if valid != want.ValidUTF8 {
				t.Errorf("valid = %v, want %v", valid, want.ValidUTF8)
			}
		})
	}
}
