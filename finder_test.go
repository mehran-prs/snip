package main

import (
	"fmt"
	"os"
	"path"
	"path/filepath"
	"testing"
)

func makeTree(t *testing.T, basePath string, files ...string) {
	for _, f := range files {
		appendLineSeparator := f[len(f)-1] == os.PathSeparator
		f = path.Join(basePath, f)
		if appendLineSeparator { // if file name ends with "/", add to res too.
			f = f + string(os.PathSeparator)
		}

		assertEqual(t, os.MkdirAll(filepath.Dir(f), 0755), nil)
		if f[len(f)-1] != os.PathSeparator {
			assertEqual(t, os.WriteFile(f, []byte(fmt.Sprintf("The %s file", f)), 0644), nil)
		}
	}
}

func TestBasename(t *testing.T) {
	assertEqual(t, baseName("abc"), "abc")
	assertEqual(t, baseName("abc.txt"), "abc.txt")
	assertEqual(t, baseName(path.Join("123", "abc.yaml")), "abc.yaml")
	assertEqual(t, baseName("abc"+string(os.PathSeparator)), "")
	assertEqual(t, baseName(""), "")
	assertEqual(t, baseName("."), "")
}

func TestFindFiles(t *testing.T) {
	searchDir := t.TempDir()
	Cfg = &Config{Dir: searchDir}

	paths := []string{
		"check.txt",
		"check.md",
		"/check/hi.md",
		"/check/hi.yaml",
		"/check2/b/cart.yaml",
		"/abc/def/cart.yaml",
	}

	makeTree(t, searchDir, paths...)

	cases := []struct {
		tag       string
		searchDir string
		search    string
		exclude   []string
		prepend   string
		res       []string
	}{
		{
			tag: "t1",
			res: []string{
				"abc/",
				"check/",
				"check",
				"check.txt",
				"check2/",
			}},
		{
			tag:     "t2",
			exclude: []string{"abc", "check2", "check.txt"},
			res: []string{
				"check/",
				"check",
			}},
		{
			tag:     "t3",
			exclude: []string{"abc", "check2", "check.txt"},
			prepend: "abc",
			res: []string{
				"abc/check/",
				"abc/check",
			}},
		{
			tag:    "t4",
			search: "abc",
			res: []string{
				"abc/",
			},
		},
		{
			tag:    "t5",
			search: "check",
			res: []string{
				"check/",
				"check",
				"check.txt",
				"check2/",
			},
		},
		{
			tag:       "t6",
			searchDir: path.Join(searchDir, "abc"),
			res: []string{
				"def/",
			},
		},
		{
			tag:       "t7",
			searchDir: path.Join(searchDir, "abc/def"),
			search:    "ca",
			res:       []string{"cart.yaml"},
		},
		{
			tag:       "t8",
			searchDir: path.Join(searchDir, "abc/def"),
			search:    "check",
		},
	}

	// Search
	for _, c := range cases {
		t.Run(c.tag, func(t *testing.T) {
			root := c.searchDir
			if root == "" {
				root = searchDir
			}
			res, err := findFiles(root, c.search, c.exclude, c.prepend)
			assertEqual(t, err, nil)
			assertEqualSlice(t, res, c.res)
		})
	}
}
