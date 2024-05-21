package main

import "fmt"

func genFzfZshCompletion(appName string) string {
	return fmt.Sprintf(`
### fzf Autocomplete function
_fzf_complete_%[1]s() {
  _fzf_path_completion $(%[1]s dir $prefix) "$@"
}
`, appName)
}

func genFzfBashCompletion(appName string) string {
	return fmt.Sprintf(`
### fzf Autocomplete function
_fzf_complete_%[1]s() {
  _fzf_path_completion $(%[1]s dir $prefix) "$@"
}

type __fzf_defc &>/dev/null && __fzf_defc %[1]s _fzf_complete_%[1]s "-o default -o bashdefault"
`, appName)
}
