package main

import (
	"fmt"
	"os"
	"os/exec"
	"path"
	"path/filepath"
	"strings"
	"sync"
)

const prefix = "SNIP" // Env prefix.

var cfgOnce sync.Once // Singleton config instance.
var Cfg *Config

type Config struct {
	Dir               string // Snippets dir.
	FileViewerCMD     []string
	MarkdownViewerCMD []string
	Editor            string
	Git               string
	Exclude           []string // exclude dirs/files. e.g., .git, .idea,...
	Verbose           bool
	LogTmpFileName    string
}

func loadConfig(globalPrefix string, appPrefix string) (err error) {
	env := func(name string, def ...string) string {
		return DefaultStr(
			os.Getenv(strings.ToUpper(appPrefix+"_"+name)),
			append([]string{os.Getenv(strings.ToUpper(globalPrefix + "_" + name))}, def...)...,
		)
	}
	cfgOnce.Do(func() {
		var homeDir string
		homeDir, err = os.UserHomeDir()
		if err != nil {
			return
		}
		Cfg = &Config{
			Dir:            strings.TrimSuffix(env("dir", path.Join(homeDir, "snippets")), string(os.PathSeparator)),
			Editor:         env("editor", os.Getenv("EDITOR"), "vim"),
			Git:            env("git", "git"),
			Verbose:        env("verbose", "") != "",
			LogTmpFileName: env("log_tmp_filename"),
		}

		if exclude := env("exclude", ".git,.idea"); exclude != "" {
			Cfg.Exclude = strings.Split(exclude, ",")
		}

		Cfg.FileViewerCMD, err = parseCommand(env("file_viewer_cmd", "cat"))
		if err != nil {
			return
		}
		Cfg.MarkdownViewerCMD, err = parseCommand(env("markdown_viewer_cmd", "cat"))
		if err != nil {
			return
		}

		// Validation
		if len(Cfg.FileViewerCMD) == 0 || len(Cfg.MarkdownViewerCMD) == 0 {
			err = fmt.Errorf(
				`invalid viewer commands. file viwer cmd: %s, markdown view cmd: %s`,
				Cfg.FileViewerCMD,
				Cfg.MarkdownViewerCMD,
			)
			return
		}
	})

	Verbose("Prefix: ", prefix, " app_prefix: ", appPrefix, " Config: ", fmt.Sprintf("%#v", Cfg))
	return
}

func (c *Config) ViewerCmd(fname string) *exec.Cmd {
	// if it's markdown, use markdown viewer, otherwise use file viewer
	params := c.FileViewerCMD
	if strings.HasSuffix(fname, ".md") {
		params = c.MarkdownViewerCMD
	}

	cmd := Command(params[0], params[1:]...)
	cmd.Args = append(cmd.Args, fname)
	return cmd
}

func (c *Config) SnippetPath(name string) string {
	fname := name
	if !strings.HasPrefix(name, c.Dir) { // join the snippets dir with the file.
		fname = JoinPaths(c.Dir, name)
	}

	if EndsWithDirectoryPath(fname) { // If it's a path to a directory, return it.
		return fname
	}

	// if file doesn't have an extension, append markdown extension to it.
	if filepath.Ext(fname) == "" {
		return fname + ".md"
	}

	return fname // return the file with its own extension.
}
