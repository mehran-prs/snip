#--------------------------------
# FZF autocompletion integration
#--------------------------------

#__snip_fzf_debug()
#{
#    local file="$BASH_COMP_DEBUG_FILE"
#    if [[ -n ${file} ]]; then
#        echo "$*" >> "${file}"
#    fi
#}

### fzf Autocomplete functions
_fzf_complete_snip() {
  # local tokens=(${(z)1}) # convert the command to array. e.g., snip -c abc => tokens=["snip","-c","abc"]
  _fzf_path_completion $(snip dir $prefix) "$1"
}
