package main

import (
	"os"
	"path"
	"path/filepath"
	"strings"
	"sync"
	"testing"
)

func setEnv(t *testing.T, key, value string) {
	t.Helper()

	oldVal := os.Getenv(key)
	t.Setenv(key, value)
	t.Cleanup(func() {
		t.Setenv(key, oldVal)
	})
}

func resetConfig() {
	cfgOnce = sync.Once{}
}

func TestLoadDefaultConfig(t *testing.T) {
	defer resetConfig()
	setEnv(t, "EDITOR", "abc")
	// Default Values
	homeDir, err := os.UserHomeDir()
	assertEqual(t, err, nil)

	assertEqual(t, loadConfig("TEST", "TEST"), nil)
	assertEqual(t, Cfg.Dir, path.Join(homeDir, "snippets"))
	assertEqualSlice(t, Cfg.FileViewerCMD, []string{"cat"})
	assertEqualSlice(t, Cfg.MarkdownViewerCMD, []string{"cat"})
	assertEqual(t, Cfg.Editor, "abc")
	assertEqual(t, Cfg.Git, "git")
	assertEqualSlice(t, Cfg.Exclude, []string{".git", ".idea"})
	assertEqual(t, Cfg.Verbose, false)
	assertEqual(t, Cfg.LogTmpFileName, "")
}

func TestLoadConfig(t *testing.T) {
	defer resetConfig()
	// Default Values
	setEnv(t, "TEST_DIR", "/ab/c")
	setEnv(t, "TEST_FILE_VIEWER_CMD", "touch a")
	setEnv(t, "TEST_MARKDOWN_VIEWER_CMD", "touch b")
	setEnv(t, "TEST_EDITOR", "vi")
	setEnv(t, "TEST_GIT", "abc")
	setEnv(t, "TEST_EXCLUDE", ".a,.b")
	setEnv(t, "TEST_VERBOSE", "TRUE")
	setEnv(t, "TEST_LOG_TMP_FILENAME", "abc.log")

	assertEqual(t, loadConfig("TEST", "TEST"), nil)

	assertEqual(t, Cfg.Dir, "/ab/c")
	assertEqualSlice(t, Cfg.FileViewerCMD, []string{"touch", "a"})
	assertEqualSlice(t, Cfg.MarkdownViewerCMD, []string{"touch", "b"})
	assertEqual(t, Cfg.Editor, "vi")
	assertEqual(t, Cfg.Git, "abc")
	assertEqualSlice(t, Cfg.Exclude, []string{".a", ".b"})
	assertTrue(t, Cfg.Verbose)
	assertEqual(t, Cfg.LogTmpFileName, "abc.log")
}

func TestLoadConfigInheritance(t *testing.T) {
	defer resetConfig()
	setEnv(t, "TEST_DIR", "/ab/c")
	setEnv(t, "SNIP_DIR", "/ab/d")
	setEnv(t, "SNIP_GIT", "abc")
	assertEqual(t, loadConfig("TEST", "TEST"), nil)

	assertEqual(t, Cfg.Dir, "/ab/c")
	assertEqual(t, Cfg.Dir, "/ab/c")
}

func TestConfig_ViewerCmd(t *testing.T) {
	Cfg = &Config{
		MarkdownViewerCMD: []string{"abc", "def"},
		FileViewerCMD:     []string{"123"},
	}

	cmd := Cfg.ViewerCmd("abc.md")
	assertEqualSlice(t, cmd.Args, []string{"abc", "def", "abc.md"})
	cmd = Cfg.ViewerCmd("abc.yaml")
	assertEqualSlice(t, cmd.Args, []string{"123", "abc.yaml"})
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
