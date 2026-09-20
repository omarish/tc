package main

import (
	"fmt"
	"sync"

	"github.com/pkoukk/tiktoken-go"
	tiktokenloader "github.com/pkoukk/tiktoken-go-loader"
)

// defaultEncodingName is the encoding used by tc in v1.
// A future -e/--encoding flag will plug in here; for now we always return o200k_base.
const defaultEncodingName = "o200k_base"

var (
	encOnce sync.Once
	encVal  *tiktoken.Tiktoken
	encErr  error
)

// encoding returns the tokenizer for this run of tc.
// Isolated so a future CLI flag can choose among encodings without touching call sites.
func encoding() (*tiktoken.Tiktoken, error) {
	encOnce.Do(func() {
		// The offline loader keeps the vocabulary embedded in the binary, so
		// there are no network calls and no cache directory at runtime.
		tiktoken.SetBpeLoader(tiktokenloader.NewOfflineLoader())
		// v1: always o200k_base. Future: map -e/--encoding (or model name) here.
		encVal, encErr = tiktoken.GetEncoding(defaultEncodingName)
		if encErr != nil {
			encErr = fmt.Errorf("encoding %s: %w", defaultEncodingName, encErr)
		}
	})
	return encVal, encErr
}

// countTokens returns the number of o200k_base tokens in text.
func countTokens(text string) (int, error) {
	enc, err := encoding()
	if err != nil {
		return 0, err
	}
	// EncodeOrdinary treats "<|endoftext|>" and friends as the literal text the
	// user actually typed, which is what tc is being asked to count.
	return len(enc.EncodeOrdinary(text)), nil
}
