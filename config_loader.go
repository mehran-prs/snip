package main

import (
	"fmt"
	"os"
	"path"
	"strings"
	"sync"
)

const prefix = "SNIP_"

var cfgOnce sync.Once // Singleton config instance.
var Cfg *Config

func userHomeDir() string {
	home, err := os.UserHomeDir()
	if err != nil {
		home = "/"
	}

	return home
}

func loadConfig(envPrefix string) (err error) {
	env := func(name string, def ...string) string {
		return defaultStr(os.Getenv(strings.ToUpper(envPrefix+name)), def...)
	}
	cfgOnce.Do(func() {
		Cfg = &Config{
			Dir:            env("dir", path.Join(userHomeDir(), "snippets")),
			Editor:         env("editor", "vim"),
			Verbose:        env("verbose", "") != "",
			LogTmpFileName: env("log_tmp_file"),
		}

		if exclude := env("exclude", ".git,.idea"); exclude != "" {
			Cfg.Exclude = LowerAll(strings.Split(exclude, ",")...)
		}

		Cfg.FileViewerCMD, err = parseCommand(env("file_viewer_cmd", "bat --style plain --paging never"))
		if err != nil {
			return
		}
		Cfg.MarkdownViewerCMD, err = parseCommand(env("markdown_viewer_cmd", "glow"))
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

	Verbose("Config: ", fmt.Sprintf("%#v", Cfg))
	return
}
