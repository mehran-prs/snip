package main

import "testing"

func TestGenFzfCompletion(t *testing.T) {
	res := `
### fzf Autocomplete function
_fzf_complete_abc() {
  _fzf_path_completion $(abc dir $prefix) "$@"
}
`
	assertEqual(t, genFzfZshCompletion("abc"), res)
	assertEqual(t, genFzfBashCompletion("abc"), res)
}
