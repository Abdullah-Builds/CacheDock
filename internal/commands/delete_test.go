package commands

import (
	"net"
	"testing"

	"github.com/Abdullah-Builds/cache/internal/cache"
)

func TestDelete(t *testing.T) {
	c := cache.New()
	c.Set("user", "Alice", 0)
	client, server := net.Pipe()
	defer client.Close()
	defer server.Close()

	response, ok := Delete([]string{"DELETE", "user"}, c, server)
	if !ok || response != "OK" {
		t.Fatalf("Delete() = (%q, %t), want (OK, true)", response, ok)
	}
	if _, found := c.Get("user"); found {
		t.Fatal("expected deleted key to be absent")
	}
}
