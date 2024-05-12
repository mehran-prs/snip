package main

import (
	"fmt"
	"os"
	"strings"
	"sync"

	"github.com/knadh/koanf/providers/env"
	"github.com/knadh/koanf/providers/structs"
	"github.com/knadh/koanf/v2"
)

const (
	prefix    = "SNIP_"
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

// SNIP_DEBUG: debug
// SNIP_ABC__DEF: abc.def
func envToPathConverter(prefix string) func(string) string {
	return func(source string) string {
		base := strings.ToLower(strings.TrimPrefix(source, prefix))
		return strings.ReplaceAll(base, separator, delimiter)
	}
}

func loadConfig(envPrefix string) (err error) {
	cfgOnce.Do(func() {
		k := koanf.New(delimiter)

		// load default configuration from default function
		if err = k.Load(structs.Provider(defaultConfig(), structTag), nil); err != nil {
			err = fmt.Errorf("error loading default: %w", err)
			return
		}

		// load default environment variables config
		if err := k.Load(env.Provider(envPrefix, delimiter, envToPathConverter(envPrefix)), nil); err != nil {
			Warn("error loading environment variables", "err", err)
		}

		Cfg = &Config{}
		if err = k.UnmarshalWithConf("", Cfg, koanf.UnmarshalConf{Tag: structTag}); err != nil {
			err = fmt.Errorf("error unmarshalling config: %w", err)
		}

		// Special changes:
		// make exclude list lowercase:
		Cfg.Exclude = LowerAll(Cfg.Exclude...)
	})
	return
}
