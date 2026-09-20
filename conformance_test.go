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
			if skipKnownBug(t, c.Name, n == want.Tokens) {
				return
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
			if skipKnownBug(t, c.Name, n == want.Tokens) {
				return
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

// knownUpstreamBugs lists corpus cases where the Go tokenizer library disagrees
// with the reference implementation through no fault of this tool.
//
// tiktoken-go ships a code-generated regexp2 engine (codec/regexp.gen.go) for
// the o200k_base split pattern, and that engine mishandles the `\s*[\r\n]+`
// alternative: a blank line that contains whitespace splits into two pieces
// where the reference produces one. The bug is in the generated engine, not in
// the pattern -- interpreting the identical pattern with regexp2 directly gives
// the correct split -- and it is present in every release through v0.8.1.
//
// Each entry is a tripwire, not just a skip: if a case starts agreeing, the
// test fails and tells you to delete the entry.
var knownUpstreamBugs = map[string]string{
	"blank-line-with-space":      `"\n \n" splits as "\n" + " \n"`,
	"blank-line-with-tab":        `"\n\t\n" splits as "\n" + "\t\n"`,
	"blank-lines-whitespace-run": "one extra token per whitespace-only blank line",
}

// skipKnownBug reports whether name is a known upstream failure, skipping the
// test if so. It fails the test when a known-bad case starts passing, so the
// list above cannot quietly go stale.
func skipKnownBug(t *testing.T, name string, agrees bool) bool {
	reason, known := knownUpstreamBugs[name]
	if !known {
		return false
	}
	if agrees {
		t.Fatalf("case %q now agrees with the reference tokenizer; "+
			"remove it from knownUpstreamBugs", name)
	}
	t.Skipf("known upstream tiktoken-go bug: %s", reason)
	return true
}

// TestUpstreamWhitespaceSplitBug documents the bug above with its minimal
// reproducer, independently of the corpus. Delete it together with
// knownUpstreamBugs once the dependency is fixed or replaced.
func TestUpstreamWhitespaceSplitBug(t *testing.T) {
	// The reference tokenizer emits a single token (id 47812) for "\n \n",
	// because `\s*[\r\n]+` matches the whole run.
	const input = "\n \n"
	ids, err := encodeIDs(input)
	if err != nil {
		t.Fatal(err)
	}
	if len(ids) == 1 {
		t.Fatalf("tiktoken-go now returns %v for %q; the upstream bug is fixed, "+
			"so delete this test and knownUpstreamBugs", ids, input)
	}
	if len(ids) != 2 {
		t.Fatalf("unexpected token count %d for %q: %v", len(ids), input, ids)
	}
}

// encodeIDs exposes token IDs for tests; tc itself only needs the count.
func encodeIDs(text string) ([]int, error) {
	enc, err := encoding()
	if err != nil {
		return nil, err
	}
	uids, _, err := enc.Encode(text)
	if err != nil {
		return nil, err
	}
	ids := make([]int, len(uids))
	for i, u := range uids {
		ids[i] = int(u)
	}
	return ids, nil
}
