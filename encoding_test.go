package main

import "testing"

func TestEncodingIsO200k(t *testing.T) {
	enc, err := encoding()
	if err != nil {
		t.Fatal(err)
	}
	if got := enc.GetName(); got != "o200k_base" {
		t.Fatalf("encoding name = %q, want o200k_base", got)
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
