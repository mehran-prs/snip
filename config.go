package main

import (
	"os"
	"os/exec"
	"path"
	"path/filepath"
	"strings"
)

type Config struct {
	Dir               string   `cfg:"dir"` // Snippets dir.
	FileViewerCMD     []string `cfg:"file_viewer"`
	MarkdownViewerCMD []string `cfg:"markdown_viewer"`
	Editor            string   `cfg:"editor"`
	Git               string   `cfg:"git"`
	Exclude           []string `cfg:"exclude"` // exclude dirs/files. e.g., .git, .idea,...
	Verbose           bool     `json:"verbose"`
	LogTmpFileName    string   `cfg:"log_tmp_file"`
}

func (c *Config) ViewerCmd(fpath string) *exec.Cmd {
	// if it's markdown, use markdown viewer, otherwise use file viewer
	params := c.FileViewerCMD
	if strings.HasSuffix(fpath, ".md") {
		params = c.MarkdownViewerCMD
	}

	cmd := Command(params[0], params[1:]...)
	cmd.Args = append(cmd.Args, fpath)
	return cmd
}

func (c *Config) SnippetPath(name string) string {
	fname := name
	if !strings.HasPrefix(name, c.Dir) {
		fname = path.Join(c.Dir, name)
	}

	stat, err := os.Stat(fname)
	if err == nil && !stat.IsDir() { // If file exists and is not a directory, return its name
		return fname
	}

	// if it's in view mode or file doesn't have an extension, append markdown extension to it.
	if filepath.Ext(fname) == "" {
		return fname + ".md"
	}

	return fname // return the file with its own extension.
}
