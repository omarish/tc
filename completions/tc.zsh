#compdef tc
# zsh completion for tc

_arguments -s -S \
  '(-h --help)'{-h,--help}'[show help]' \
  '(-v --version)'{-v,--version}'[print version]' \
  '*:file:_files'
