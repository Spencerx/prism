package internal

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
)

const (
	configDirName  = "prism"
	configFileName = "config.json"
)

var configPathOverride = os.Getenv("PRISM_CONFIG_FILE")

type Config struct {
	Verbose   bool `json:"verbose"`
	OnlyFails bool `json:"only_fails"`
	NoLogo    bool `json:"no_logo"`
}

// GlobalConfig holds the active configuration for the current process.
var GlobalConfig Config

// LoadConfig reads the persisted configuration file, if it exists.
func LoadConfig() (Config, error) {
	cfg := Config{}

	path, err := configFilePath()
	if err != nil {
		return cfg, err
	}

	data, err := os.ReadFile(path)
	if err != nil {
		if errors.Is(err, fs.ErrNotExist) {
			return cfg, nil
		}
		return cfg, fmt.Errorf("read config file: %w", err)
	}

	if len(data) == 0 {
		return cfg, nil
	}

	if err := json.Unmarshal(data, &cfg); err != nil {
		return cfg, fmt.Errorf("parse config file: %w", err)
	}

	return cfg, nil
}

// SaveConfig writes the provided configuration to disk, persisting supported settings.
func SaveConfig(cfg Config) error {
	path, err := configFilePath()
	if err != nil {
		return err
	}

	dir := filepath.Dir(path)
	if err := os.MkdirAll(dir, 0o755); err != nil {
		return fmt.Errorf("create config directory: %w", err)
	}

	payload, err := json.MarshalIndent(cfg, "", "  ")
	if err != nil {
		return fmt.Errorf("encode config: %w", err)
	}
	payload = append(payload, '\n')

	tmpPath := path + ".tmp"
	if err := os.WriteFile(tmpPath, payload, 0o644); err != nil {
		return fmt.Errorf("write config temp file: %w", err)
	}

	if err := os.Rename(tmpPath, path); err != nil {
		if removeErr := os.Remove(path); removeErr != nil && !errors.Is(removeErr, fs.ErrNotExist) {
			_ = os.Remove(tmpPath)
			return fmt.Errorf("replace config file: %w", err)
		}
		if err := os.Rename(tmpPath, path); err != nil {
			_ = os.Remove(tmpPath)
			return fmt.Errorf("replace config file: %w", err)
		}
	}

	return nil
}

func configFilePath() (string, error) {
	if configPathOverride != "" {
		return configPathOverride, nil
	}

	dir, err := os.UserConfigDir()
	if err != nil {
		return "", fmt.Errorf("resolve user config dir: %w", err)
	}

	return filepath.Join(dir, configDirName, configFileName), nil
}
