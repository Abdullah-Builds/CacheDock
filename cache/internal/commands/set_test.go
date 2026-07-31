package commands

import (
	"net"
	"testing"

	"github.com/Abdullah-Builds/cache/internal/cache"
)

func TestSet(t *testing.T) {
	c := cache.New()
	client, server := net.Pipe()
	defer client.Close()
	defer server.Close()

	response, ok := Set([]string{"SET", "user", "Alice"}, c, server)
	if !ok || response != "OK" {
		t.Fatalf("Set() = (%q, %t), want (OK, true)", response, ok)
	}
	if value, found := c.Get("user"); !found || value != "Alice" {
		t.Fatalf("cache value = %q (found=%t)", value, found)
	}
}
