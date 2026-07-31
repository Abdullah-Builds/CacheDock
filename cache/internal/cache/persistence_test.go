package cache

import (
	"os"
	"path/filepath"
	"testing"
)

func TestSave(t *testing.T) {
	c := New()
	c.Set("user", "Alice", 0)
	filename := filepath.Join(t.TempDir(), "cache.json")

	if err := c.Save(filename); err != nil {
		t.Fatalf("Save() error = %v", err)
	}
	if info, err := os.Stat(filename); err != nil || info.Size() == 0 {
		t.Fatalf("snapshot was not written: %v", err)
	}
}

func TestLoad(t *testing.T) {
	filename := filepath.Join(t.TempDir(), "cache.json")
	source := New()
	source.Set("user", "Alice", 0)
	if err := source.Save(filename); err != nil {
		t.Fatalf("Save() error = %v", err)
	}

	loaded := New()
	if err := loaded.Load(filename); err != nil {
		t.Fatalf("Load() error = %v", err)
	}
	value, ok := loaded.Get("user")
	if !ok || value != "Alice" {
		t.Fatalf("expected restored value Alice, got %q (found=%t)", value, ok)
	}
}
