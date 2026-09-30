package config

import (
	"fmt"
	"os"
	"time"

	"github.com/BurntSushi/toml"
)

type Config struct {
	Site         `toml:"site"`
	Server       `toml:"server"`
	Session      `toml:"session"`
	Registration `toml:"registration"`
}

type Site struct {
	Name         string `toml:"name"`
	Description  string `toml:"description"`
	DefaultTheme string `toml:"default_theme"`
}

type Server struct {
	Address string `toml:"address"`
}

type Session struct {
	Secure   bool          `toml:"secure"`
	Lifetime time.Duration `toml:"lifetime"`
}

type Registration struct {
	Enabled            bool `toml:"enabled"`
	RequireEmailVerify bool `toml:"require_email_verify"`
	InviteOnly         bool `toml:"invite_only"`
}

var instance *Config

// Get provides global access to application config
func Get() *Config {
	if instance == nil {
		panic("Get called before initialization")
	}
	return instance
}

// MustLoad reads the config.toml file and globally initializes a Config.
// Panics on os errors or unparsable config
func MustLoad() {
	// config.toml
	data, err := os.ReadFile("config.toml")
	if err != nil {
		panic(fmt.Sprintf("unable to read config file: %v", err))
	}
	cfg := &Config{}
	if _, err := toml.Decode(string(data), cfg); err != nil {
		panic(fmt.Sprintf("error parsing config: %v", err))
	}
	instance = cfg
}

// MustEnv asserts that all required environment vars are set.
// Panics on missing environment vars.
func MustEnv(required ...string) {
	var missing []string
	for _, key := range required {
		if _, ok := os.LookupEnv(key); !ok {
			missing = append(missing, key)
		}
	}
	if len(missing) > 0 {
		panic(fmt.Sprintf("missing env values: %v", missing))
	}
}
