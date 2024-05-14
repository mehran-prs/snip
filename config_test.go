package main

import (
	"os"
	"path"
	"path/filepath"
	"strings"
	"testing"
)

func TestLoadDefaultConfig(t *testing.T) {
	os.Setenv("EDITOR", "abc")
	// Default Values
	homeDir, err := os.UserHomeDir()
	assertEqual(t, err, nil)

	assertEqual(t, loadConfig("TEST_"), nil)
	assertEqual(t, Cfg.Dir, path.Join(homeDir, "snippets"))
	assertEqualSlice(t, Cfg.FileViewerCMD, []string{"bat", "--style", "plain", "--paging", "never"})
	assertEqualSlice(t, Cfg.MarkdownViewerCMD, []string{"glow"})
	assertEqual(t, Cfg.Editor, "abc")
	assertEqual(t, Cfg.Git, "git")
	assertEqualSlice(t, Cfg.Exclude, []string{".git", ".idea"})
	assertEqual(t, Cfg.Verbose, false)
	assertEqual(t, Cfg.LogTmpFileName, "")
}

func TestLoadDConfig(t *testing.T) {
	// Default Values
	assertEqual(t, os.Setenv("TEST_DIR", "/ab/c"), nil)
	assertEqual(t, os.Setenv("TEST_FILE_VIEWER_CMD", "touch a"), nil)
	assertEqual(t, os.Setenv("TEST_MARKDOWN_VIEWER_CMD", "touch b"), nil)
	assertEqual(t, os.Setenv("TEST_EDITOR", "vi"), nil)
	assertEqual(t, os.Setenv("TEST_GIT", "abc"), nil)
	assertEqual(t, os.Setenv("TEST_EXCLUDE", ".a,.b"), nil)
	assertEqual(t, os.Setenv("TEST_VERBOSE", "TRUE"), nil)
	assertEqual(t, os.Setenv("TEST_LOG_TMP_FILENAME", "abc.log"), nil)

	assertEqual(t, loadConfig("TEST_"), nil)

	assertEqual(t, Cfg.Dir, "/ab/c")
	assertEqualSlice(t, Cfg.FileViewerCMD, []string{"touch", "a"})
	assertEqualSlice(t, Cfg.MarkdownViewerCMD, []string{"touch", "b"})
	assertEqual(t, Cfg.Editor, "vi")
	assertEqual(t, Cfg.Git, "abc")
	assertEqualSlice(t, Cfg.Exclude, []string{".a", ".b"})
	assertEqual(t, Cfg.Verbose, true)
	assertEqual(t, Cfg.LogTmpFileName, "abc.log")
}

func TestConfig_ViewerCmd(t *testing.T) {
	Cfg = &Config{
		MarkdownViewerCMD: []string{"abc", "def"},
		FileViewerCMD:     []string{"123", "456"},
	}

	cmd := Cfg.ViewerCmd("abc.md")
	assertEqualSlice(t, cmd.Args, []string{"abc", "def", "abc.md"})
	cmd = Cfg.ViewerCmd("abc.yaml")
	assertEqualSlice(t, cmd.Args, []string{"123", "456", "abc.yaml"})
}

func TestConfig_SnippetPath(t *testing.T) {
	f, err := os.CreateTemp("", "*abc.yaml")
	assertTrue(t, err == nil)
	defer func() {
		_ = f.Close()
		_ = os.Remove(f.Name())
	}()

	filedir := filepath.Dir(f.Name())
	filename := filepath.Base(f.Name())
	filenameNoExtension := strings.TrimSuffix(filename, ".yaml")

	Cfg = &Config{
		Dir: filedir,
	}

	assertEqual(t, Cfg.SnippetPath(f.Name()), f.Name())
	assertEqual(t, Cfg.SnippetPath(filename), f.Name())
	assertEqual(t, Cfg.SnippetPath(filenameNoExtension), path.Join(filedir, filenameNoExtension+".md"))
	assertEqual(t, Cfg.SnippetPath("abc"), Cfg.Dir+"/abc.md")
}
