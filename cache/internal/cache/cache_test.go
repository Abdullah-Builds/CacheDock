package cache

import (
	"testing"
	"time"
)

func TestSetAndGet(t *testing.T) {

	cache := New()

	cache.Set("user", "Alice", 0)

	value, ok := cache.Get("user")

	if !ok {
		t.Fatal("expected key to exist")
	}

	if value != "Alice" {
		t.Fatalf("expected Alice got %s", value)
	}
}

func TestDelete(t *testing.T) {
	c := New()
	c.Set("user", "Alice", 0)
	c.Delete("user")

	if _, ok := c.Get("user"); ok {
		t.Fatal("expected deleted key to be absent")
	}
}

func TestMissingKey(t *testing.T) {
	c := New()
	if value, ok := c.Get("missing"); ok || value != "" {
		t.Fatalf("expected missing key, got %q (found=%t)", value, ok)
	}
}

func TestTTLExpiration(t *testing.T) {
	c := New()
	c.Set("session", "active", 20*time.Millisecond)
	time.Sleep(40 * time.Millisecond)

	if _, ok := c.Get("session"); ok {
		t.Fatal("expected expired key to be absent")
	}
}

func TestOverwrite(t *testing.T) {
	c := New()
	c.Set("user", "Alice", 0)
	c.Set("user", "Bob", 0)

	value, ok := c.Get("user")
	if !ok || value != "Bob" {
		t.Fatalf("expected overwritten value Bob, got %q (found=%t)", value, ok)
	}
}
