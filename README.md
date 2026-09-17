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
tc [file ...]
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

## Encoding

- **v1:** always `o200k_base`. Options: `-h` / `--help`, `-v` / `--version`.
- **Future:** `-e` / `--encoding` will select among encodings; see `encoding.go`.

No network calls at runtime, no config files, no model APIs. The tokenizer
vocabulary is embedded in the binary via
[`github.com/tiktoken-go/tokenizer`](https://github.com/tiktoken-go/tokenizer).

## License

MIT — see [LICENSE](LICENSE).
