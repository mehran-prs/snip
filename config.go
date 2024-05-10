package main

import (
	"os"
	"os/exec"
	"path"
	"path/filepath"
	"strings"
)

type Config struct {
	SnippetsDir       string        `cfg:"snippets_dir"`
	FileViewerCMD     CommandConfig `cfg:"file_viewer"`
	MarkdownViewerCMD CommandConfig `cfg:"markdown_viewer"`
	Editor            string        `cfg:"editor"`
	Git               string        `cfg:"git"`
	Exclude           []string      `cfg:"exclude"` // exclude dirs/files. e.g., .git, .idea,...
	LogLevel          string        `cfg:"log_level"`
}

type CommandConfig struct {
	Name string   `cfg:"name"`
	Args []string `cfg:"args"`
}

func defaultConfig() *Config {
	return &Config{
		SnippetsDir: path.Join(userHomeDir(), "snippets"),
		FileViewerCMD: CommandConfig{
			Name: "bat",
			Args: []string{"--style", "plain", "--paging", "never"},
		},
		MarkdownViewerCMD: CommandConfig{Name: "glow"},
		Editor:            os.Getenv("EDITOR"),
		Git:               "git",
	}
}

func (c *Config) ViewerCmd(fpath string) *exec.Cmd {
	// if it's markdown, use markdown viewer, otherwise use file viewer
	params := c.FileViewerCMD
	if strings.HasSuffix(fpath, ".md") {
		params = c.MarkdownViewerCMD
	}

	cmd := Command(params.Name, params.Args...)
	cmd.Args = append(cmd.Args, fpath)
	return cmd
}

func (c *Config) SnippetPath(name string) string {
	fname := name
	if !strings.HasPrefix(name, c.SnippetsDir) {
		fname = path.Join(c.SnippetsDir, name)
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
