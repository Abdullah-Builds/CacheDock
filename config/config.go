package config

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"
)

const defaultConfigFile = "config/config.yaml"

// LoadConfig reads a small YAML configuration file and then applies environment
// variable overrides. CONFIG_FILE selects a different YAML file when needed.
func LoadConfig() (*Config, error) {
	config := DefaultConfig()
	configFile := os.Getenv("CONFIG_FILE")
	if configFile == "" {
		configFile = defaultConfigFile
	}

	if err := loadYAML(configFile, &config); err != nil {
		return nil, err
	}

	if port := os.Getenv("PORT"); port != "" {
		config.Port = port
	}
	if file := os.Getenv("DATA_FILE"); file != "" {
		config.DataFile = file
	}
	if cleanup := os.Getenv("CLEANUP_INTERVAL"); cleanup != "" {
		if d, err := time.ParseDuration(cleanup); err == nil {
			config.CleanupInterval = d
		}
	}
	if autosave := os.Getenv("AUTOSAVE_INTERVAL"); autosave != "" {
		if d, err := time.ParseDuration(autosave); err == nil {
			config.AutoSaveInterval = d
		}
	}
	return &config, nil
}

type Config struct {
	Port             string
	MaxKeys          int
	DataFile         string
	CleanupInterval  time.Duration
	AutoSaveInterval time.Duration
}

func DefaultConfig() Config {
	return Config{
		Port:             "8080",
		MaxKeys:          100,
		DataFile:         "data/cache.json",
		CleanupInterval:  time.Hour,
		AutoSaveInterval: 5 * time.Minute,
	}
}

func loadYAML(filename string, config *Config) error {
	file, err := os.Open(filename)
	if err != nil {
		return fmt.Errorf("open config %q: %w", filename, err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for lineNumber := 1; scanner.Scan(); lineNumber++ {
		line := strings.TrimSpace(scanner.Text())
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}

		key, value, found := strings.Cut(line, ":")
		if !found {
			return fmt.Errorf("config %q line %d: expected key: value", filename, lineNumber)
		}
		key = strings.TrimSpace(key)
		value = strings.Trim(strings.TrimSpace(value), "\"'")

		var parseErr error
		switch key {
		case "port":
			config.Port = value
		case "maxKeys":
			config.MaxKeys, parseErr = strconv.Atoi(value)
			if parseErr == nil && config.MaxKeys <= 0 {
				parseErr = fmt.Errorf("must be greater than zero")
			}
		case "cleanupInterval":
			config.CleanupInterval, parseErr = time.ParseDuration(value)
		case "autosaveInterval":
			config.AutoSaveInterval, parseErr = time.ParseDuration(value)
		case "dataFile":
			config.DataFile = value
		default:
			return fmt.Errorf("config %q line %d: unknown key %q", filename, lineNumber, key)
		}
		if parseErr != nil {
			return fmt.Errorf("config %q line %d: invalid %s: %w", filename, lineNumber, key, parseErr)
		}
	}
	if err := scanner.Err(); err != nil {
		return fmt.Errorf("read config %q: %w", filename, err)
	}
	return nil
}
