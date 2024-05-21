package main

import "testing"

func TestGenFzfCompletion(t *testing.T) {
	zshRes := `
### fzf Autocomplete function
_fzf_complete_abc() {
  _fzf_path_completion $(abc dir $prefix) "$@"
}
`
	bashRes := `
### fzf Autocomplete function
_fzf_complete_abc() {
  _fzf_path_completion $(abc dir $prefix) "$@"
}

type __fzf_defc &>/dev/null && __fzf_defc abc _fzf_complete_abc "-o default -o bashdefault"
`
	assertEqual(t, genFzfZshCompletion("abc"), zshRes)
	assertEqual(t, genFzfBashCompletion("abc"), bashRes)
}
