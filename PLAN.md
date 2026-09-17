# Plan

Locked v1 decisions (implemented):

1. Count tokens in files or stdin only.
2. Always OpenAI `o200k_base`.
3. Only `-h` / `--help` as a CLI option; no encoding flags yet.
4. Encoding selection isolated in `encoding()` for a future `-e/--encoding`.
5. Module path `github.com/omarish/tc`.
6. `wc`-like interface (stdin / per-file / total / continue on errors).
