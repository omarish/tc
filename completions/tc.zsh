#compdef tc
# zsh completion for tc

_arguments -s -S \
  '(-h --help)'{-h,--help}'[show help]' \
  '(-v --version)'{-v,--version}'[print version]' \
  '--strict[fail on input that is not valid UTF-8]' \
  '*:file:_files'
