package main

import "testing"

// The tokenizer library does not expose the encoding name, so pin the encoding
// by behaviour instead: these IDs are o200k_base and nothing else.
func TestEncodingIsO200k(t *testing.T) {
	enc, err := encoding()
	if err != nil {
		t.Fatal(err)
	}
	for _, c := range []struct {
		text string
		want []int
	}{
		{"Hello world", []int{13225, 2375}},
		{"hello", []int{24912}},
		// A special-token marker must count as the literal text it is.
		{"<|endoftext|>", []int{27, 91, 419, 1440, 919, 91, 29}},
	} {
		got := enc.EncodeOrdinary(c.text)
		if len(got) != len(c.want) {
			t.Errorf("EncodeOrdinary(%q) = %v, want %v", c.text, got, c.want)
			continue
		}
		for i := range got {
			if got[i] != c.want[i] {
				t.Errorf("EncodeOrdinary(%q) = %v, want %v", c.text, got, c.want)
				break
			}
		}
	}
}

func TestCountTokensHello(t *testing.T) {
	// "hello" is a single token under o200k_base.
	n, err := countTokens("hello")
	if err != nil {
		t.Fatal(err)
	}
	if n != 1 {
		t.Fatalf("countTokens(%q) = %d, want 1", "hello", n)
	}
}

func TestCountTokensEmpty(t *testing.T) {
	n, err := countTokens("")
	if err != nil {
		t.Fatal(err)
	}
	if n != 0 {
		t.Fatalf("countTokens(%q) = %d, want 0", "", n)
	}
}
