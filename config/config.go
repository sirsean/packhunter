package config

import (
	"code.google.com/p/gcfg"
	"log"
)

type Config struct {
	Mongo struct {
		Hosts        string
		AuthDatabase string
		AuthUsername string
		AuthPassword string
		Database     string
	}
	Host struct {
		Port string
		Path string
	}
	CookieStore struct {
		AuthenticationKey string
	}
}

var cfg Config
var loaded bool

func Get() Config {
	if !loaded {
		Load()
	}
	return cfg
}

func Load() {
	err := gcfg.ReadFileInto(&cfg, "/etc/friendly-ph/friendly-ph.gcfg")
	if err != nil {
		log.Printf("Failed to read config: %v", err)
	} else {
		loaded = true
	}
}
