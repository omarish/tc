# tc

Like `wc`, but for tokens.

[![CI](https://github.com/omarish/tc/actions/workflows/ci.yml/badge.svg)](https://github.com/omarish/tc/actions/workflows/ci.yml)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](LICENSE)

## What it does

`tc` counts tokens in files or stdin using OpenAI's **`o200k_base`** encoding
(the encoding used by GPT-4o and related models).

v1 has one job: plain `tc` always means `o200k_base`.
Options: `-h` / `--help` for usage, `-v` / `--version` for the version string.

A future `-e` / `--encoding` flag is planned; encoding selection is already
isolated in one place (`encoding()` in `encoding.go`) so that flag can plug in
without rewriting the CLI.

## Install

### Homebrew

```bash
brew install omarish/tap/tc
```

### Prebuilt binaries

Download the latest release for your OS/arch from:

**[github.com/omarish/tc/releases/latest](https://github.com/omarish/tc/releases/latest)**

Each release includes:

- Raw binaries (`tc_<os>_<arch>` / `.exe`) for direct download
- Archives (`tc_<version>_<os>_<arch>.tar.gz` or `.zip`) for packaging / Homebrew
- A `SHA256SUMS` file covering every asset (verify with `sha256sum -c SHA256SUMS`)

| Platform | Binary | Archive |
|---|---|---|
| macOS Apple Silicon | [`tc_darwin_arm64`](https://github.com/omarish/tc/releases/latest/download/tc_darwin_arm64) | [`tc_*_darwin_arm64.tar.gz`](https://github.com/omarish/tc/releases/latest) |
| macOS Intel | [`tc_darwin_amd64`](https://github.com/omarish/tc/releases/latest/download/tc_darwin_amd64) | [`tc_*_darwin_amd64.tar.gz`](https://github.com/omarish/tc/releases/latest) |
| Linux x86_64 | [`tc_linux_amd64`](https://github.com/omarish/tc/releases/latest/download/tc_linux_amd64) | [`tc_*_linux_amd64.tar.gz`](https://github.com/omarish/tc/releases/latest) |
| Linux arm64 | [`tc_linux_arm64`](https://github.com/omarish/tc/releases/latest/download/tc_linux_arm64) | [`tc_*_linux_arm64.tar.gz`](https://github.com/omarish/tc/releases/latest) |
| Windows x86_64 | [`tc_windows_amd64.exe`](https://github.com/omarish/tc/releases/latest/download/tc_windows_amd64.exe) | [`tc_*_windows_amd64.zip`](https://github.com/omarish/tc/releases/latest) |
| Windows arm64 | [`tc_windows_arm64.exe`](https://github.com/omarish/tc/releases/latest/download/tc_windows_arm64.exe) | [`tc_*_windows_arm64.zip`](https://github.com/omarish/tc/releases/latest) |

Example (macOS Apple Silicon):

```bash
curl -L -o tc https://github.com/omarish/tc/releases/latest/download/tc_darwin_arm64
chmod +x tc
sudo mv tc /usr/local/bin/tc
```

Or from the archive (Homebrew-style):

```bash
curl -L -o tc.tar.gz https://github.com/omarish/tc/releases/latest/download/tc_0.1.0_darwin_arm64.tar.gz
tar -xzf tc.tar.gz
sudo mv tc /usr/local/bin/tc
```

> Links resolve after the first tagged release (`v0.1.0`). Until then, build from source below.
> Replace `0.1.0` with the release version; checksums are in `SHA256SUMS` on that release.

### From source

Requires Go 1.21+.

```bash
go install github.com/omarish/tc@latest
```

Or:

```bash
git clone https://github.com/omarish/tc.git
cd tc
make          # embeds VERSION=dev by default
# VERSION=0.1.0 make
```

### Man page

After installing the binary (and optionally the man page from the repo):

```bash
# from a clone
sudo mkdir -p /usr/local/share/man/man1
sudo cp man/tc.1 /usr/local/share/man/man1/
man tc
```

### Shell completions

From a clone of this repo:

**bash** — source or install into bash-completion:

```bash
# session
source completions/tc.bash

# or install (Debian/Ubuntu-style)
sudo cp completions/tc.bash /etc/bash_completion.d/tc
```

**zsh** — add the completions dir to `fpath` or copy the file:

```bash
# e.g. oh-my-zsh custom completions
mkdir -p ~/.oh-my-zsh/custom/completions
cp completions/tc.zsh ~/.oh-my-zsh/custom/completions/_tc
# ensure fpath includes that dir, then: compinit
```

**fish**:

```bash
mkdir -p ~/.config/fish/completions
cp completions/tc.fish ~/.config/fish/completions/tc.fish
```

## Usage

```text
tc [--strict] [file ...]
tc -h | --help
tc -v | --version
```

Behavior mirrors `wc`:

| Invocation | Output |
|---|---|
| `tc` (no files) | Read stdin; print just the token count |
| `tc file` | Print `COUNT filename` |
| `tc file1 file2 ...` | One `COUNT filename` line per file, then `COUNT total` |
| Unreadable file | Error to stderr, continue others; exit non-zero if any failed |
| `tc -h` / `tc --help` | Print usage and exit 0 |
| `tc -v` / `tc --version` | Print `tc <version>` and exit 0 |
| Input that is not valid UTF-8 | Warn to stderr, still count on stdout; exit 0 |
| `tc --strict` on such input | Error to stderr, no count, exit non-zero |
| Unknown option (e.g. `-e`) | Error to stderr with a hint to `-h`; exit non-zero |

### Examples

```bash
echo -n "hello" | tc
# 1

tc README.md
# <n> README.md

tc a.txt b.txt
# <n> a.txt
# <m> b.txt
# <n+m> total

cat notes.txt | tc
```

## Unicode

Valid UTF-8 is counted exactly, including the awkward cases. A ZWJ emoji
sequence really does cost 11 tokens:

```bash
printf '\U0001F468\u200D\U0001F469\u200D\U0001F467\u200D\U0001F466' | tc
# 11
```

Input is expected to be UTF-8. Input that **isn't** — a Latin-1 file, a UTF-16
file, a truncated write, a stray binary — still gets counted, because that is
what `wc` would do, but `tc` tells you:

```bash
$ printf 'ok \xf0\x9f\x91' | tc          # "ok " + a truncated emoji
tc: warning: (standard input): not valid UTF-8; counted with U+FFFD substitution
2
```

The count goes to stdout and the warning to stderr, so pipelines keep working
and `2>/dev/null` silences it. Use `--strict` to reject such input instead:

```bash
$ printf 'ok \xf0\x9f\x91' | tc --strict
tc: (standard input): not valid UTF-8
$ echo $?
1
```

**Why the substitution matters.** The API only ever sees text a JSON encoder
produced, so the count that matters is the one you get *after* replacement.
`tc` applies the WHATWG "maximal subpart" rule — each maximal invalid byte
sequence becomes exactly one U+FFFD — which is what Python's
`bytes.decode(errors="replace")` does, and therefore what your HTTP client
does. Substituting per *byte* instead (the easy mistake) turns a truncated
4-byte emoji into three U+FFFD rather than one, and since BPE merges runs of
U+FFFD, that changes the count.

## Conformance

`tc` is tested against the reference Python
[`tiktoken`](https://github.com/openai/tiktoken) across a corpus of unicode
edge cases — emoji and ZWJ sequences, CJK, RTL scripts, combining marks,
zero-width and bidi controls, BOMs, overlong and surrogate encodings,
truncated sequences, Latin-1 and UTF-16 mistaken for UTF-8, and binary junk.

- [`testdata/corpus.json`](testdata/corpus.json) — the inputs
- [`testdata/expected.json`](testdata/expected.json) — token IDs and counts,
  **generated from real `tiktoken`**, never edited by hand
- [`scripts/gen_expected.py`](scripts/gen_expected.py) — regenerates it

`go test` needs no Python: it asserts against the committed answers. CI runs
the generator with `--check` on every push, so if `tiktoken` ever changes, or
someone edits the golden file by hand, the build fails with a diff.

```bash
go test ./...                              # offline
pip install tiktoken
python scripts/gen_expected.py --check     # what CI does
```

### Known upstream limitation

Three corpus cases are skipped, and the skips are tripwires that fail if the
cases start passing. `tiktoken-go` ships a code-generated regexp2 engine for
the `o200k_base` split pattern, and that engine mishandles the `\s*[\r\n]+`
alternative: a **blank line containing whitespace** splits into two pieces
where the reference tokenizer produces one.

```
input        reference   here
"a\n \nb"            3      4
"a\n\t\nb"           3      4
```

It costs one extra token per whitespace-only blank line, so it compounds on
text that has many of them. Ordinary blank lines, CRLF line endings, trailing
spaces and markdown hard breaks are all unaffected — none of the 35 real files
in these two repos hits it.

The bug is in the generated engine, not the pattern: interpreting the identical
pattern with `regexp2` directly gives the correct split. It is present in every
`tiktoken-go` release through v0.8.1. `pkoukk/tiktoken-go` (with its offline
loader, so still no network calls) tokenizes all three cases correctly and is
the likely fix.

To add a case, append it to `corpus.json` (`text` for readable input, `hex` for
raw bytes) and rerun the generator.

## Encoding

- **v1:** always `o200k_base`. Options: `-h` / `--help`, `-v` / `--version`,
  `--strict`.
- **Future:** `-e` / `--encoding` will select among encodings; see `encoding.go`.

No network calls at runtime, no config files, no model APIs. The tokenizer
vocabulary is embedded in the binary via
[`github.com/tiktoken-go/tokenizer`](https://github.com/tiktoken-go/tokenizer).

## License

MIT — see [LICENSE](LICENSE).
