package commands

import (
	"net"
	"testing"

	"github.com/Abdullah-Builds/cache/internal/cache"
)

func TestGet(t *testing.T) {
	c := cache.New()
	c.Set("user", "Alice", 0)
	client, server := net.Pipe()
	defer client.Close()
	defer server.Close()

	response, ok := Get([]string{"GET", "user"}, c, server)
	if !ok || response != "Alice" {
		t.Fatalf("Get() = (%q, %t), want (Alice, true)", response, ok)
	}
}
