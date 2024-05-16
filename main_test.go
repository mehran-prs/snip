package main

import (
	"os"
	"path"
	"testing"

	"github.com/spf13/cobra"
)

func TestBoot(t *testing.T) {
	defer resetConfig()

	cmd := &cobra.Command{}
	setEnv(t, prefix+"DIR", "/a/b/c")
	assertEqual(t, boot(cmd, nil), nil)
	assertEqual(t, Cfg.Dir, "/a/b/c")
}

func TestBootAndShutdown(t *testing.T) {
	defer resetConfig()
	f, err := os.CreateTemp("", "abc")
	assertEqual(t, err, nil)
	defer func() {
		_ = f.Close()
		_ = os.Remove(f.Name())
	}()

	cmd := &cobra.Command{Use: "test"}
	setEnv(t, "TEST_DIR", "/a/b/c")
	setEnv(t, "TEST_LOG_TMP_FILENAME", f.Name())
	setEnv(t, prefix+"DIR", "/a/b/d")
	setEnv(t, prefix+"GIT", "abc")
	assertEqual(t, boot(cmd, nil), nil)
	assertEqual(t, Cfg.Dir, "/a/b/c")
	assertEqual(t, Cfg.Git, "abc")

	// Make sure logger is tuned.
	assertTrue(t, logFile != nil)
	assertEqual(t, logFile.Name(), f.Name())

	assertEqual(t, shutdown(cmd, nil), nil)
	assertEqual(t, logFile, nil)
}

func TestCobraAutocompleteFilename(t *testing.T) {
	defer resetConfig()

	searchDir := t.TempDir()
	Cfg = &Config{Dir: searchDir, Exclude: []string{".git"}}

	paths := []string{
		"check.txt",
		"check.md",
		"/check/hi.md",
		"/check/hi.yaml",
		"/check2/b/cart.yaml",
		"/abc/def/cart.yaml",
		".git/a/b.txt",
	}

	makeTree(t, searchDir, paths...)

	cases := []struct {
		tag        string
		toComplete string
		res        []string
	}{
		{
			tag: "t1",
			res: []string{
				"abc/",
				"check/",
				"check",
				"check.txt",
				"check2/",
			},
		},
		{
			tag:        "t4",
			toComplete: "abc",
			res: []string{
				"abc/",
			},
		},
		{
			tag:        "t5",
			toComplete: "check",
			res: []string{
				"check/",
				"check",
				"check.txt",
				"check2/",
			},
		},
		{
			tag:        "t6",
			toComplete: "abc/",
			res: []string{
				"abc/def/",
			},
		},
		{
			tag:        "t7",
			toComplete: "abc/def/ca",
			res:        []string{"abc/def/cart.yaml"},
		},
		{
			tag:        "t8",
			toComplete: "abc/def/check",
		},
	}

	// Search
	for _, c := range cases {
		t.Run(c.tag, func(t *testing.T) {
			res, _ := cobraAutoCompleteFileName(nil, nil, c.toComplete)
			assertEqualSlice(t, res, c.res)
		})
	}
}

func TestCmdViewSnippet(t *testing.T) {
	tmpDir := t.TempDir()
	Cfg = &Config{
		Dir:               tmpDir,
		MarkdownViewerCMD: []string{"touch"},
		FileViewerCMD:     []string{"touch"},
	}

	assertEqual(t, CmdViewSnippet(nil, []string{"a.md"}), nil)
	_, err := os.Stat(path.Join(tmpDir, "a.md"))
	assertEqual(t, err, nil)

	assertEqual(t, CmdViewSnippet(nil, []string{"a.yaml"}), nil)
	_, err = os.Stat(path.Join(tmpDir, "a.yaml"))
	assertEqual(t, err, nil)
}

func TestCmdSnippetsDir(t *testing.T) {
	tmpDir := t.TempDir()
	Cfg = &Config{Dir: tmpDir}
	assertEqual(t, CmdSnippetsDir(nil, nil), nil)
	assertEqual(t, CmdSnippetsDir(nil, []string{"abc"}), nil)
}

func TestCmdEditSnippet(t *testing.T) {
	tmpDir := t.TempDir()
	Cfg = &Config{
		Dir:    tmpDir,
		Editor: "touch",
	}

	assertEqual(t, CmdEditSnippet(nil, []string{"a.md"}), nil)
	_, err := os.Stat(path.Join(tmpDir, "a.md"))
	assertEqual(t, err, nil)

	assertEqual(t, CmdEditSnippet(nil, []string{"a/b.md"}), nil)
	_, err = os.Stat(path.Join(tmpDir, "a/b.md"))
	assertEqual(t, err, nil)

	assertEqual(t, CmdEditSnippet(nil, []string{"c"}), nil)
	_, err = os.Stat(path.Join(tmpDir, "c.md"))
	assertEqual(t, err, nil)

	assertEqual(t, CmdEditSnippet(nil, []string{"c.yaml"}), nil)
	_, err = os.Stat(path.Join(tmpDir, "c.yaml"))
	assertEqual(t, err, nil)

	// Open snippets dir
	Cfg.Dir = path.Join(tmpDir, "a.yaml")
	assertEqual(t, CmdEditSnippet(nil, nil), nil)
	_, err = os.Stat(path.Join(tmpDir, "a.yaml"))
	assertEqual(t, err, nil)
}
