package config

import (
	"os"
	"path/filepath"

	"github.com/spf13/viper"
)

// avairable config and default
var configs = map[string]interface{}{
	"global.snippets_path": filepath.Join(dataPath(), "snippets"),
}

// Load config.toml
func Load() {
	for k, v := range configs {
		viper.SetDefault(k, v)
	}
	viper.SetConfigType("toml")
	viper.SetConfigName("config")
	viper.AddConfigPath(configPath())
	_ = viper.ReadInConfig()
}

// Set value to config
func Set(k, v string) error {
	viper.Set(k, v)
	return nil
}

// Get value from config
func Get(k string) string {
	switch v := viper.Get(k).(type) {
	case string:
		return v
	}
	return ""
}

// Save config.toml
func Save() error {
	return viper.WriteConfigAs(filepath.Join(configPath(), "config.toml"))
}

func configPath() string {
	p := os.Getenv("XDG_CONFIG_HOME")
	if p == "" {
		home := os.Getenv("HOME")
		p = filepath.Join(home, ".config")
	}
	return filepath.Join(p, "frgm")
}

func dataPath() string {
	p := os.Getenv("XDG_DATA_HOME")
	if p == "" {
		home := os.Getenv("HOME")
		p = filepath.Join(home, ".local", "share")
	}
	return filepath.Join(p, "frgm")
}
