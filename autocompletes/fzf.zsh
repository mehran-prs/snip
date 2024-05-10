#--------------------------------
# FZF autocompletion integration
#--------------------------------

### Helper functions
_resolve_alias_to_cmd(){
   local tokens=(${(z)1})
   local cmd=$tokens[1]
   local alias_definition=$(alias $cmd)
   local resolved_cmd=${alias_definition#*=}
   echo ${1/$cmd/${(Q)resolved_cmd}}
}


### fzf Autocomplete functions
_fzf_complete_snip() {
  local tokens=(${(z)1}) # convert the command to array. e.g., snip -c abc => tokens=["snip","-c","abc"]
  # The value of the -c flag:
  local cfg_path=""
  for ((i = 1; i <= $#tokens; i++)); do
    if [[ ${tokens[i]} == "-c" || ${tokens[i]} == "--config" ]]; then
      # The value of the -c flag will be the next element in the array
      cfg_path=${tokens[i + 1]}
      break
    fi
  done
  if [ cfg_path!="" ]; then
    cfg_path="-c $cfg_path"
  fi

  _fzf_path_completion $(snip $cfg_path dir $prefix) "$1"
}

_snip_aliases_fzf_complete(){
  # Store all parameters in an array
  local params=("$@")
  params[1]= $(_resolve_alias_to_cmd $1)
  _fzf_complete_snip "${params[@]}"
}

# Register teh fzf completion for all aliases:
SNIP_ALIASES=${SNIP_ALIASES:-()}
for snip_alias in $SNIP_ALIASES
do
   _fzf_complete_${snip_alias}() {
     _snip_aliases_fzf_complete
  }
done


