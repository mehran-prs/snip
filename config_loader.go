package main

import (
	"fmt"
	"os"
	"sync"

	"github.com/knadh/koanf/parsers/yaml"
	"github.com/knadh/koanf/providers/file"
	"github.com/knadh/koanf/providers/structs"
	"github.com/knadh/koanf/v2"
)

const (
	delimiter = "."
	separator = "__"
	structTag = "cfg" // the tag on the config struct that we use to map config path to the struct fields.
)

var cfgOnce sync.Once // Singleton config instance.
var Cfg *Config

func userHomeDir() string {
	home, err := os.UserHomeDir()
	if err != nil {
		home = "/"
	}
	return home
}

func loadConfig(configFile string, logLevel string) (err error) {
	cfgOnce.Do(func() {
		k := koanf.New(delimiter)

		// load default configuration from default function
		if err = k.Load(structs.Provider(defaultConfig(), structTag), nil); err != nil {
			err = fmt.Errorf("error loading default: %w", err)
			return
		}

		// load from config file
		if err = k.Load(file.Provider(configFile), yaml.Parser()); err != nil {
			if !os.IsNotExist(err) { // Ignore and just log it
				err = fmt.Errorf("can not load config file %s: %w", configFile, err)
				return
			}
			Debug("config file not found, ignoring it.", "file:", configFile)
		}

		Cfg = &Config{}
		if err = k.UnmarshalWithConf("", Cfg, koanf.UnmarshalConf{Tag: structTag}); err != nil {
			err = fmt.Errorf("error unmarshalling config: %w", err)
		}

		// Override log level
		if logLevel != "" {
			Cfg.LogLevel = logLevel
		}
	})
	return
}
