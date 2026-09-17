# Plan

Locked v1 decisions (implemented):

1. Count tokens in files or stdin only.
2. Always OpenAI `o200k_base`.
3. CLI options: `-h` / `--help`, `-v` / `--version`; no encoding flags yet.
4. Encoding selection isolated in `encoding()` for a future `-e/--encoding`.
5. Module path `github.com/omarish/tc`.
6. `wc`-like interface (stdin / per-file / total / continue on errors).

Ergonomics landed:

- Version via `main.version` + ldflags (`-v` / `--version`); Makefile + release workflow.
- Richer `-h` with Examples; man page `man/tc.1`.
- Shell completions (bash / zsh / fish).
- Release assets: raw binaries + `tc_<ver>_<os>_<arch>.tar.gz`/`.zip` + `SHA256SUMS`.
- README: Homebrew tap `omarish/tap/tc`, checksums, man, completions.
