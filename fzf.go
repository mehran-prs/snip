package main

import "fmt"

func genFzfZshCompletion(appName string) string {
	return fmt.Sprintf(`
### fzf Autocomplete function
_fzf_complete_snip() {
  _fzf_path_completion $(%[1]s dir $prefix) "$@"
}
`, appName)
}

func genFzfBashCompletion(appName string) string {
	return fmt.Sprintf(`
### fzf Autocomplete function
_fzf_complete_snip() {
  _fzf_path_completion $(%[1]s dir $prefix) "$@"
}
`, appName)
}
