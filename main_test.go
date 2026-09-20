package main

import (
	"bytes"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestRunStdin(t *testing.T) {
	var out, errBuf bytes.Buffer
	code := run(nil, strings.NewReader("hello"), &out, &errBuf)
	if code != 0 {
		t.Fatalf("exit %d, stderr %q", code, errBuf.String())
	}
	if got := strings.TrimSpace(out.String()); got != "1" {
		t.Fatalf("stdout = %q, want %q", got, "1")
	}
}

func TestRunOneFile(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "a.txt")
	if err := os.WriteFile(path, []byte("hello"), 0o644); err != nil {
		t.Fatal(err)
	}
	var out, errBuf bytes.Buffer
	code := run([]string{path}, strings.NewReader(""), &out, &errBuf)
	if code != 0 {
		t.Fatalf("exit %d, stderr %q", code, errBuf.String())
	}
	want := "1 " + path + "\n"
	if out.String() != want {
		t.Fatalf("stdout = %q, want %q", out.String(), want)
	}
}

func TestRunMultiFileTotal(t *testing.T) {
	dir := t.TempDir()
	a := filepath.Join(dir, "a.txt")
	b := filepath.Join(dir, "b.txt")
	if err := os.WriteFile(a, []byte("hello"), 0o644); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(b, []byte("hello"), 0o644); err != nil {
		t.Fatal(err)
	}
	var out, errBuf bytes.Buffer
	code := run([]string{a, b}, strings.NewReader(""), &out, &errBuf)
	if code != 0 {
		t.Fatalf("exit %d, stderr %q", code, errBuf.String())
	}
	lines := strings.Split(strings.TrimSuffix(out.String(), "\n"), "\n")
	if len(lines) != 3 {
		t.Fatalf("got %d lines: %q", len(lines), out.String())
	}
	if lines[0] != "1 "+a {
		t.Fatalf("line1 = %q", lines[0])
	}
	if lines[1] != "1 "+b {
		t.Fatalf("line2 = %q", lines[1])
	}
	if lines[2] != "2 total" {
		t.Fatalf("total line = %q, want %q", lines[2], "2 total")
	}
}

func TestRunUnreadableContinues(t *testing.T) {
	dir := t.TempDir()
	ok := filepath.Join(dir, "ok.txt")
	missing := filepath.Join(dir, "missing.txt")
	if err := os.WriteFile(ok, []byte("hello"), 0o644); err != nil {
		t.Fatal(err)
	}
	var out, errBuf bytes.Buffer
	code := run([]string{missing, ok}, strings.NewReader(""), &out, &errBuf)
	if code == 0 {
		t.Fatal("expected non-zero exit")
	}
	if !strings.Contains(errBuf.String(), missing) {
		t.Fatalf("stderr should mention missing file: %q", errBuf.String())
	}
	if !strings.Contains(out.String(), "1 "+ok) {
		t.Fatalf("stdout should still count ok file: %q", out.String())
	}
}

func TestRunHelp(t *testing.T) {
	for _, arg := range []string{"-h", "--help"} {
		var out, errBuf bytes.Buffer
		code := run([]string{arg}, strings.NewReader(""), &out, &errBuf)
		if code != 0 {
			t.Fatalf("%s: exit %d, stderr %q", arg, code, errBuf.String())
		}
		if !strings.Contains(out.String(), "Usage: tc") {
			t.Fatalf("%s: stdout missing usage: %q", arg, out.String())
		}
		if !strings.Contains(out.String(), "-v") {
			t.Fatalf("%s: stdout should mention -v: %q", arg, out.String())
		}
		if errBuf.Len() != 0 {
			t.Fatalf("%s: unexpected stderr %q", arg, errBuf.String())
		}
	}
}

func TestRunVersion(t *testing.T) {
	for _, arg := range []string{"-v", "--version"} {
		var out, errBuf bytes.Buffer
		code := run([]string{arg}, strings.NewReader(""), &out, &errBuf)
		if code != 0 {
			t.Fatalf("%s: exit %d, stderr %q", arg, code, errBuf.String())
		}
		got := out.String()
		if !strings.HasPrefix(got, "tc ") {
			t.Fatalf("%s: stdout = %q, want prefix %q", arg, got, "tc ")
		}
		if !strings.Contains(got, version) {
			t.Fatalf("%s: stdout = %q, want to contain version %q", arg, got, version)
		}
		if !strings.HasSuffix(got, "\n") {
			t.Fatalf("%s: stdout missing trailing newline: %q", arg, got)
		}
		if errBuf.Len() != 0 {
			t.Fatalf("%s: unexpected stderr %q", arg, errBuf.String())
		}
	}
}

func TestRunUnknownOption(t *testing.T) {
	var out, errBuf bytes.Buffer
	code := run([]string{"-e"}, strings.NewReader(""), &out, &errBuf)
	if code == 0 {
		t.Fatal("expected non-zero exit for unknown option")
	}
	if !strings.Contains(errBuf.String(), "unknown option") {
		t.Fatalf("stderr = %q", errBuf.String())
	}
	if !strings.Contains(errBuf.String(), "tc -h") {
		t.Fatalf("stderr should hint -h: %q", errBuf.String())
	}
}

func TestRunDoubleDashFile(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "-h")
	if err := os.WriteFile(path, []byte("hello"), 0o644); err != nil {
		t.Fatal(err)
	}
	var out, errBuf bytes.Buffer
	code := run([]string{"--", path}, strings.NewReader(""), &out, &errBuf)
	if code != 0 {
		t.Fatalf("exit %d, stderr %q", code, errBuf.String())
	}
	want := "1 " + path + "\n"
	if out.String() != want {
		t.Fatalf("stdout = %q, want %q", out.String(), want)
	}
}

// invalidUTF8 is "ok " followed by a truncated 4-byte emoji.
var invalidUTF8 = []byte("ok \xf0\x9f\x91")

func TestRunWarnsOnInvalidUTF8(t *testing.T) {
	var out, errBuf bytes.Buffer
	code := run(nil, bytes.NewReader(invalidUTF8), &out, &errBuf)
	if code != 0 {
		t.Fatalf("exit %d, want 0 (a warning must not fail the run)", code)
	}
	// stdout stays machine-readable: the count and nothing else.
	if got := out.String(); got != "2\n" {
		t.Fatalf("stdout = %q, want %q", got, "2\n")
	}
	e := errBuf.String()
	if !strings.Contains(e, "warning") || !strings.Contains(e, "not valid UTF-8") {
		t.Fatalf("stderr = %q, want a UTF-8 warning", e)
	}
	if !strings.Contains(e, stdinName) {
		t.Fatalf("stderr should name stdin: %q", e)
	}
}

func TestRunNoWarningOnValidUTF8(t *testing.T) {
	var out, errBuf bytes.Buffer
	code := run(nil, strings.NewReader("hello 日本語 👍"), &out, &errBuf)
	if code != 0 {
		t.Fatalf("exit %d, stderr %q", code, errBuf.String())
	}
	if errBuf.Len() != 0 {
		t.Fatalf("unexpected stderr for valid input: %q", errBuf.String())
	}
}

func TestRunStrictRejectsInvalidUTF8(t *testing.T) {
	var out, errBuf bytes.Buffer
	code := run([]string{"--strict"}, bytes.NewReader(invalidUTF8), &out, &errBuf)
	if code == 0 {
		t.Fatal("expected non-zero exit under --strict")
	}
	if out.Len() != 0 {
		t.Fatalf("--strict must not print a count: %q", out.String())
	}
	if !strings.Contains(errBuf.String(), "not valid UTF-8") {
		t.Fatalf("stderr = %q", errBuf.String())
	}
	if strings.Contains(errBuf.String(), "warning") {
		t.Fatalf("--strict should report an error, not a warning: %q", errBuf.String())
	}
}

func TestRunStrictAllowsValidUTF8(t *testing.T) {
	var out, errBuf bytes.Buffer
	code := run([]string{"--strict"}, strings.NewReader("hello"), &out, &errBuf)
	if code != 0 {
		t.Fatalf("exit %d, stderr %q", code, errBuf.String())
	}
	if got := strings.TrimSpace(out.String()); got != "1" {
		t.Fatalf("stdout = %q, want 1", got)
	}
}

// Under --strict a bad file behaves like an unreadable one: report it, skip it,
// keep going, exit non-zero.
func TestRunStrictSkipsBadFileContinues(t *testing.T) {
	dir := t.TempDir()
	good := filepath.Join(dir, "good.txt")
	bad := filepath.Join(dir, "bad.txt")
	if err := os.WriteFile(good, []byte("hello"), 0o644); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(bad, invalidUTF8, 0o644); err != nil {
		t.Fatal(err)
	}
	var out, errBuf bytes.Buffer
	code := run([]string{"--strict", bad, good}, strings.NewReader(""), &out, &errBuf)
	if code == 0 {
		t.Fatal("expected non-zero exit")
	}
	if strings.Contains(out.String(), bad) {
		t.Fatalf("bad file should not be counted: %q", out.String())
	}
	if !strings.Contains(out.String(), "1 "+good) {
		t.Fatalf("good file should still be counted: %q", out.String())
	}
	// Total reflects only the files that were counted.
	if !strings.Contains(out.String(), "1 total") {
		t.Fatalf("total should be 1: %q", out.String())
	}
}

// Without --strict the bad file is still counted, and the total includes it.
func TestRunWarnsPerFileAndStillTotals(t *testing.T) {
	dir := t.TempDir()
	good := filepath.Join(dir, "good.txt")
	bad := filepath.Join(dir, "bad.txt")
	if err := os.WriteFile(good, []byte("hello"), 0o644); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(bad, invalidUTF8, 0o644); err != nil {
		t.Fatal(err)
	}
	var out, errBuf bytes.Buffer
	code := run([]string{bad, good}, strings.NewReader(""), &out, &errBuf)
	if code != 0 {
		t.Fatalf("exit %d, stderr %q", code, errBuf.String())
	}
	if !strings.Contains(errBuf.String(), bad) {
		t.Fatalf("warning should name the bad file: %q", errBuf.String())
	}
	if strings.Contains(errBuf.String(), good) {
		t.Fatalf("no warning expected for the good file: %q", errBuf.String())
	}
	if !strings.Contains(out.String(), "3 total") {
		t.Fatalf("total should be 2+1=3: %q", out.String())
	}
}

func TestParseArgsStrict(t *testing.T) {
	opts, err := parseArgs([]string{"--strict", "a.txt"})
	if err != nil {
		t.Fatal(err)
	}
	if !opts.strict {
		t.Error("--strict not set")
	}
	if len(opts.files) != 1 || opts.files[0] != "a.txt" {
		t.Errorf("files = %v", opts.files)
	}
}
