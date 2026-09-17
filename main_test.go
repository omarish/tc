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
