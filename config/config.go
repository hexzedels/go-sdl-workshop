package config

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"

	"gopkg.in/yaml.v3"
)

// Config holds the application configuration.
type Config struct {
	Server   ServerConfig   `yaml:"server"`
	Database DatabaseConfig `yaml:"database"`
	JWT      JWTConfig      `yaml:"jwt"`
	APIKey   string         `yaml:"api_key"`
	Internal InternalConfig `yaml:"internal"`
}

type ServerConfig struct {
	Port string `yaml:"port"`
}

type DatabaseConfig struct {
	Driver string `yaml:"driver"`
	DSN    string `yaml:"dsn"`
}

type JWTConfig struct {
	Secret     string `yaml:"secret"`
	ExpiryMins int    `yaml:"expiry_mins"`
}

type InternalConfig struct {
	SigningToken string `yaml:"signing_token"`
}

// Resolve returns the first existing config path, checking in order:
//  1. --config flag value (if non-empty)
//  2. CONFIG_PATH env var
//  3. $(dirname executable)/config.yaml
//  4. ./config.yaml
//
// Returns the empty string only if none of the candidates exist.
func Resolve(flagPath string) string {
	candidates := []string{}
	if flagPath != "" {
		candidates = append(candidates, flagPath)
	}
	if env := os.Getenv("CONFIG_PATH"); env != "" {
		candidates = append(candidates, env)
	}
	if exe, err := os.Executable(); err == nil {
		candidates = append(candidates, filepath.Join(filepath.Dir(exe), "config.yaml"))
	}
	candidates = append(candidates, "config.yaml")

	for _, p := range candidates {
		if _, err := os.Stat(p); err == nil {
			return p
		}
	}
	return ""
}

// Load reads config from the given YAML file path, then applies env overrides
// for secrets so that production deployments do not need secrets in the file.
func Load(path string) (*Config, error) {
	if path == "" {
		return nil, errors.New("config path not resolved (pass --config, CONFIG_PATH env, or place config.yaml next to binary)")
	}
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("read config %s: %w", path, err)
	}
	var cfg Config
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		return nil, fmt.Errorf("parse config %s: %w", path, err)
	}

	if v := os.Getenv("JWT_SECRET"); v != "" {
		cfg.JWT.Secret = v
	}
	if v := os.Getenv("API_KEY"); v != "" {
		cfg.APIKey = v
	}
	if v := os.Getenv("SIGNING_TOKEN"); v != "" {
		cfg.Internal.SigningToken = v
	}
	return &cfg, nil
}
