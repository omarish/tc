package main

import (
	"fmt"

	"github.com/tiktoken-go/tokenizer"
)

// defaultEncodingName is the encoding used by tc in v1.
// A future -e/--encoding flag will plug in here; for now we always return o200k_base.
const defaultEncodingName = "o200k_base"

// encoding returns the tokenizer codec for this run of tc.
// Isolated so a future CLI flag can choose among encodings without touching call sites.
func encoding() (tokenizer.Codec, error) {
	// v1: always o200k_base. Future: map -e/--encoding (or model name) to tokenizer.Encoding here.
	enc, err := tokenizer.Get(tokenizer.O200kBase)
	if err != nil {
		return nil, fmt.Errorf("encoding %s: %w", defaultEncodingName, err)
	}
	return enc, nil
}

// countTokens returns the number of o200k_base tokens in text.
func countTokens(text string) (int, error) {
	enc, err := encoding()
	if err != nil {
		return 0, err
	}
	ids, _, err := enc.Encode(text)
	if err != nil {
		return 0, err
	}
	return len(ids), nil
}
