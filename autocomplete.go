package main

import (
	"os"
	"path"
	"path/filepath"
	"slices"
	"strings"

	"github.com/labstack/gommon/color"
	"github.com/spf13/cobra"
)

func AutoCompleteFileName(dir string, exclude []string, toComplete string) ([]string, cobra.ShellCompDirective) {
	searchDir := filepath.Dir(toComplete)
	root := path.Join(dir, searchDir)
	res, err := findFiles(root, baseName(toComplete), exclude, searchDir)
	if err != nil {
		return nil, cobra.ShellCompDirectiveError
	}

	return res, cobra.ShellCompDirectiveNoSpace | cobra.ShellCompDirectiveNoFileComp
}

// baseName is like filepath.Base but returns empty string in the following cases:
// if name is just line-separator or ends with line-separator
// if name is an empty string
func baseName(name string) string {
	base := filepath.Base(name)
	if len(name) == 0 || base == "/" || name[len(name)-1] == filepath.Separator {
		return ""
	}
	return base
}

func findFiles(root string, search string, exclude []string, prepend string) ([]string, error) {
	root = strings.TrimSuffix(root, string(os.PathSeparator)) + string(os.PathSeparator)
	var result []string
	walkFn := func(path string, info os.FileInfo, err error) error {
		path = strings.TrimPrefix(path, root)
		if err != nil {
			return err
		}

		if slices.Contains(exclude, path) {
			if info.IsDir() {
				return filepath.SkipDir
			}
			return nil
		}

		if path == "" { // the root itself should be ignored
			return nil
		}

		if strings.Contains(path, search) || search == "" {
			res := filepath.Join(prepend, strings.TrimSuffix(path, ".md")) // Remove .md from end of markdown files.
			if info.IsDir() {
				res = color.Bold(color.Cyan(res) + "/") // change style of directories (colorize and append / to them)
			}

			result = append(result, res)
			if info.IsDir() {
				return filepath.SkipDir
			}
		}

		return nil
	}

	err := filepath.Walk(root, walkFn)
	if err != nil {
		return nil, err
	}

	return result, nil
}
