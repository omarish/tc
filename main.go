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

const usage = `Usage: tc [file ...]
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

Examples:
  echo -n "hello" | tc
  tc README.md
  tc a.txt b.txt
  cat notes.txt | tc

Encoding is always o200k_base in v1. A future -e/--encoding flag may
select other encodings; for now there are no counting options.
`

// run implements the wc-like CLI. It is separated from main for testing.
func run(args []string, stdin io.Reader, stdout, stderr io.Writer) int {
	files, help, showVersion, err := parseArgs(args)
	if err != nil {
		fmt.Fprintf(stderr, "tc: %v\n", err)
		return 1
	}
	if help {
		fmt.Fprint(stdout, usage)
		return 0
	}
	if showVersion {
		fmt.Fprintf(stdout, "tc %s\n", version)
		return 0
	}

	if len(files) == 0 {
		n, err := countReader(stdin)
		if err != nil {
			fmt.Fprintf(stderr, "tc: %v\n", err)
			return 1
		}
		fmt.Fprintf(stdout, "%d\n", n)
		return 0
	}

	var total int
	var failed bool
	for _, name := range files {
		n, err := countFile(name)
		if err != nil {
			fmt.Fprintf(stderr, "tc: %s: %s\n", name, errString(err))
			failed = true
			continue
		}
		fmt.Fprintf(stdout, "%d %s\n", n, name)
		total += n
	}

	if len(files) > 1 {
		fmt.Fprintf(stdout, "%d total\n", total)
	}

	if failed {
		return 1
	}
	return 0
}

// parseArgs splits CLI args into file paths.
// -h / --help requests usage. -v / --version requests version.
// -- ends option parsing.
// Any other dash-led token is an error (so typos are not treated as filenames).
func parseArgs(args []string) (files []string, help, showVersion bool, err error) {
	for i := 0; i < len(args); i++ {
		a := args[i]
		if a == "--" {
			return append(files, args[i+1:]...), false, false, nil
		}
		if a == "-h" || a == "--help" {
			return nil, true, false, nil
		}
		if a == "-v" || a == "--version" {
			return nil, false, true, nil
		}
		if strings.HasPrefix(a, "-") && a != "-" {
			return nil, false, false, fmt.Errorf("unknown option %s\nTry 'tc -h' for help.", a)
		}
		files = append(files, a)
	}
	return files, false, false, nil
}

func errString(err error) string {
	if pe, ok := err.(*fs.PathError); ok {
		return pe.Err.Error()
	}
	return err.Error()
}

func countFile(name string) (int, error) {
	f, err := os.Open(name)
	if err != nil {
		return 0, err
	}
	defer f.Close()
	return countReader(f)
}

func countReader(r io.Reader) (int, error) {
	data, err := io.ReadAll(r)
	if err != nil {
		return 0, err
	}
	return countTokens(string(data))
}
