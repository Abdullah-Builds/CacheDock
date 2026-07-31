package handler

import (
	"bufio"
	"net"
	"testing"
	"time"

	"github.com/Abdullah-Builds/cache/internal/cache"
)

func TestServerCommands(t *testing.T) {
	client, server := net.Pipe()
	defer client.Close()
	c := cache.New()
	go HandleConnection(server, c)

	writeErr := make(chan error, 1)
	go func() {
		_, err := client.Write([]byte("SET user Alice\nGET user\nDELETE user\nGET user\n"))
		writeErr <- err
	}()

	_ = client.SetReadDeadline(time.Now().Add(time.Second))
	reader := bufio.NewReader(client)
	for _, want := range []string{"OK\n", "Alice\n", "OK\n", "NOT_FOUND\n"} {
		got, err := reader.ReadString('\n')
		if err != nil {
			t.Fatalf("read response: %v", err)
		}
		if got != want {
			t.Errorf("response = %q, want %q", got, want)
		}
	}
	if err := <-writeErr; err != nil {
		t.Fatalf("write commands: %v", err)
	}
}
