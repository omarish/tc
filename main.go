package main

import (
	"fmt"
	"io"
	"io/fs"
	"os"
	"strings"
)

// version is set at link time via -ldflags "-X main.version=...".
var version = "dev"

func main() {
	os.Exit(run(os.Args[1:], os.Stdin, os.Stdout, os.Stderr))
}

const usage = `Usage: tc [options] [file ...]
       tc -h | --help
       tc -v | --version

Like wc, but for tokens. Count tokens using OpenAI o200k_base
(GPT-4o and related models).

With no files, read standard input and print the token count.
With one or more files, print the count and filename for each;
with more than one file, also print a total.

Options:
  -h, --help      show this help
  -v, --version   print version and exit
      --strict    fail on input that is not valid UTF-8

Input must be UTF-8. Input that is not valid UTF-8 is counted the way an
API client would send it: each maximal invalid byte sequence becomes one
U+FFFD. tc warns on stderr when this happens; stdout is unaffected, so
pipelines keep working. Use --strict to reject such input instead.

Examples:
  echo -n "hello" | tc
  tc README.md
  tc a.txt b.txt
  cat notes.txt | tc
  tc --strict *.md

Encoding is always o200k_base in v1. A future -e/--encoding flag may
select other encodings; for now there are no counting options.
`

// stdinName is how standard input is named in diagnostics.
const stdinName = "(standard input)"

type options struct {
	files   []string
	help    bool
	version bool
	strict  bool
}

// run implements the wc-like CLI. It is separated from main for testing.
func run(args []string, stdin io.Reader, stdout, stderr io.Writer) int {
	opts, err := parseArgs(args)
	if err != nil {
		fmt.Fprintf(stderr, "tc: %v\n", err)
		return 1
	}
	if opts.help {
		fmt.Fprint(stdout, usage)
		return 0
	}
	if opts.version {
		fmt.Fprintf(stdout, "tc %s\n", version)
		return 0
	}

	if len(opts.files) == 0 {
		n, valid, err := countReader(stdin)
		if err != nil {
			fmt.Fprintf(stderr, "tc: %v\n", err)
			return 1
		}
		if !valid {
			reportInvalid(stderr, stdinName, opts.strict)
			if opts.strict {
				return 1
			}
		}
		fmt.Fprintf(stdout, "%d\n", n)
		return 0
	}

	var total int
	var failed bool
	for _, name := range opts.files {
		n, valid, err := countFile(name)
		if err != nil {
			fmt.Fprintf(stderr, "tc: %s: %s\n", name, errString(err))
			failed = true
			continue
		}
		if !valid {
			reportInvalid(stderr, name, opts.strict)
			if opts.strict {
				failed = true
				continue
			}
		}
		fmt.Fprintf(stdout, "%d %s\n", n, name)
		total += n
	}

	if len(opts.files) > 1 {
		fmt.Fprintf(stdout, "%d total\n", total)
	}

	if failed {
		return 1
	}
	return 0
}

// reportInvalid writes the diagnostic for input that is not valid UTF-8.
// Under --strict it is an error; otherwise it is a warning and counting
// proceeds with U+FFFD substitution.
func reportInvalid(stderr io.Writer, name string, strict bool) {
	if strict {
		fmt.Fprintf(stderr, "tc: %s: not valid UTF-8\n", name)
		return
	}
	fmt.Fprintf(stderr, "tc: warning: %s: not valid UTF-8; counted with U+FFFD substitution\n", name)
}

// parseArgs splits CLI args into options and file paths.
// -h / --help requests usage. -v / --version requests version.
// --strict rejects input that is not valid UTF-8.
// -- ends option parsing.
// Any other dash-led token is an error (so typos are not treated as filenames).
func parseArgs(args []string) (options, error) {
	var opts options
	for i := 0; i < len(args); i++ {
		a := args[i]
		if a == "--" {
			opts.files = append(opts.files, args[i+1:]...)
			return opts, nil
		}
		if a == "-h" || a == "--help" {
			opts.help = true
			return opts, nil
		}
		if a == "-v" || a == "--version" {
			opts.version = true
			return opts, nil
		}
		if a == "--strict" {
			opts.strict = true
			continue
		}
		if strings.HasPrefix(a, "-") && a != "-" {
			return options{}, fmt.Errorf("unknown option %s\nTry 'tc -h' for help.", a)
		}
		opts.files = append(opts.files, a)
	}
	return opts, nil
}

func errString(err error) string {
	if pe, ok := err.(*fs.PathError); ok {
		return pe.Err.Error()
	}
	return err.Error()
}

func countFile(name string) (n int, valid bool, err error) {
	f, err := os.Open(name)
	if err != nil {
		return 0, false, err
	}
	defer f.Close()
	return countReader(f)
}

// countReader counts tokens in r, reporting whether the input was valid UTF-8.
func countReader(r io.Reader) (n int, valid bool, err error) {
	data, err := io.ReadAll(r)
	if err != nil {
		return 0, false, err
	}
	text, valid := decodeUTF8(data)
	n, err = countTokens(text)
	return n, valid, err
}
