package config

import (
	"os"
	"path/filepath"
	"testing"
	"time"
)

func TestLoadConfigReadsYAMLAndEnvironmentOverrides(t *testing.T) {
	directory := t.TempDir()
	filename := filepath.Join(directory, "config.yaml")
	contents := "port: 9000\nmaxKeys: 500\ncleanupInterval: 2s\nautosaveInterval: 1m\ndataFile: snapshots/cache.json\n"
	if err := os.WriteFile(filename, []byte(contents), 0600); err != nil {
		t.Fatal(err)
	}
	t.Setenv("CONFIG_FILE", filename)
	t.Setenv("PORT", "9010")

	config, err := LoadConfig()
	if err != nil {
		t.Fatal(err)
	}
	if config.Port != "9010" || config.MaxKeys != 500 || config.DataFile != "snapshots/cache.json" {
		t.Fatalf("unexpected config: %+v", config)
	}
	if config.CleanupInterval != 2*time.Second || config.AutoSaveInterval != time.Minute {
		t.Fatalf("unexpected intervals: cleanup=%s autosave=%s", config.CleanupInterval, config.AutoSaveInterval)
	}
}

func TestLoadConfigRejectsInvalidValues(t *testing.T) {
	filename := filepath.Join(t.TempDir(), "config.yaml")
	if err := os.WriteFile(filename, []byte("maxKeys: 0\n"), 0600); err != nil {
		t.Fatal(err)
	}
	t.Setenv("CONFIG_FILE", filename)

	if _, err := LoadConfig(); err == nil {
		t.Fatal("expected invalid maxKeys error")
	}
}
