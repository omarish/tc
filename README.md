# tc

Unix-style CLI to count tokens (like `wc`, for tokens).

[![CI](https://github.com/omarish/tc/actions/workflows/ci.yml/badge.svg)](https://github.com/omarish/tc/actions/workflows/ci.yml)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](LICENSE)

## What it does

`tc` counts tokens in files or stdin using OpenAI's **`o200k_base`** encoding
(the encoding used by GPT-4o and related models).

v1 has one job: plain `tc` always means `o200k_base`.
The only option is `-h` / `--help` for usage.

A future `-e` / `--encoding` flag is planned; encoding selection is already
isolated in one place (`encoding()` in `encoding.go`) so that flag can plug in
without rewriting the CLI.

## Install

### Prebuilt binaries

Download the latest release for your OS/arch from:

**[github.com/omarish/tc/releases/latest](https://github.com/omarish/tc/releases/latest)**

| Platform | Asset |
|---|---|
| macOS Apple Silicon | [`tc_darwin_arm64`](https://github.com/omarish/tc/releases/latest/download/tc_darwin_arm64) |
| macOS Intel | [`tc_darwin_amd64`](https://github.com/omarish/tc/releases/latest/download/tc_darwin_amd64) |
| Linux x86_64 | [`tc_linux_amd64`](https://github.com/omarish/tc/releases/latest/download/tc_linux_amd64) |
| Linux arm64 | [`tc_linux_arm64`](https://github.com/omarish/tc/releases/latest/download/tc_linux_arm64) |
| Windows x86_64 | [`tc_windows_amd64.exe`](https://github.com/omarish/tc/releases/latest/download/tc_windows_amd64.exe) |
| Windows arm64 | [`tc_windows_arm64.exe`](https://github.com/omarish/tc/releases/latest/download/tc_windows_arm64.exe) |

Example (macOS Apple Silicon):

```bash
curl -L -o tc https://github.com/omarish/tc/releases/latest/download/tc_darwin_arm64
chmod +x tc
sudo mv tc /usr/local/bin/tc
```

> Links resolve after the first tagged release (`v0.1.0`). Until then, build from source below.

### From source

Requires Go 1.21+.

```bash
go install github.com/omarish/tc@latest
```

Or:

```bash
git clone https://github.com/omarish/tc.git
cd tc
go build -o tc .
```

## Usage

```text
tc [file ...]
tc -h
```

Behavior mirrors `wc`:

| Invocation | Output |
|---|---|
| `tc` (no files) | Read stdin; print just the token count |
| `tc file` | Print `COUNT filename` |
| `tc file1 file2 ...` | One `COUNT filename` line per file, then `COUNT total` |
| Unreadable file | Error to stderr, continue others; exit non-zero if any failed |
| `tc -h` / `tc --help` | Print usage and exit 0 |
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
```

## Encoding

- **v1:** always `o200k_base`. Only `-h` / `--help` is accepted as an option.
- **Future:** `-e` / `--encoding` will select among encodings; see `encoding.go`.

No network calls at runtime, no config files, no model APIs. The tokenizer
vocabulary is embedded in the binary via
[`github.com/tiktoken-go/tokenizer`](https://github.com/tiktoken-go/tokenizer).

## License

MIT — see [LICENSE](LICENSE).
