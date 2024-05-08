package main

import (
	"log"
	"strings"
	"sync"

	"github.com/knadh/koanf/providers/env"
	"github.com/knadh/koanf/providers/structs"
	"github.com/knadh/koanf/v2"
)

var cfgOnce sync.Once // Singleton config instance.
var Cfg *Config

// BIFROST_DEBUG -> DEBUG -> debug
// BIFROST_DATABASE__HOST -> DATABASE__HOST -> database__host -> database.host

func envToPathConverter(source string) string {
	base := strings.ToLower(strings.TrimPrefix(source, prefix))
	return strings.ReplaceAll(base, separator, delimiter)
}

func GetConfig() *Config {
	cfgOnce.Do(func() {
		k := koanf.New(".")

		// load default configuration from default function
		if err := k.Load(structs.Provider(defaultConfig(), structTag), nil); err != nil {
			log.Fatalf("error loading default: %s", err)
		}

		// load environment variables
		if err := k.Load(env.Provider(prefix, delimiter, envToPathConverter), nil); err != nil {
			log.Printf("error loading environment variables: %s", err)
		}

		Cfg = &Config{}
		if err := k.UnmarshalWithConf("", Cfg, koanf.UnmarshalConf{Tag: structTag}); err != nil {
			log.Fatalf("error unmarshalling config: %s", err)
		}
	})
	return Cfg
}
