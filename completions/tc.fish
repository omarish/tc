# fish completion for tc
complete -c tc -s h -l help -d 'Show help'
complete -c tc -s v -l version -d 'Print version'
complete -c tc -l strict -d 'Fail on input that is not valid UTF-8'
# remaining args: default path/file completion
