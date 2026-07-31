package handler

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"net"
	"os"
	"sync"
	"testing"
	"time"

	"github.com/Abdullah-Builds/cache/internal/cache"
)

const (
	stressClients = 128
	stressRounds  = 10000
)

// TestTCPStress runs 2.56 million application commands through an actual local
// TCP listener. It is intentionally opt-in: run it with `go test -run TestTCPStress`.
func TestTCPStress(t *testing.T) {
	if os.Getenv("RUN_STRESS") != "1" {
		t.Skip("set RUN_STRESS=1 to run the TCP stress test")
	}

	previousOutput := log.Writer()
	log.SetOutput(io.Discard)
	defer log.SetOutput(previousOutput)

	listener, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatal(err)
	}
	defer listener.Close()

	c := cache.New()
	var handlers sync.WaitGroup
	acceptDone := make(chan struct{})
	go func() {
		defer close(acceptDone)
		for {
			conn, err := listener.Accept()
			if err != nil {
				return
			}
			handlers.Add(1)
			go func() {
				defer handlers.Done()
				HandleConnection(conn, c)
			}()
		}
	}()

	start := time.Now()
	errCh := make(chan error, stressClients)
	var clients sync.WaitGroup
	for clientID := 0; clientID < stressClients; clientID++ {
		clients.Add(1)
		go func(id int) {
			defer clients.Done()
			conn, err := net.Dial("tcp", listener.Addr().String())
			if err != nil {
				errCh <- err
				return
			}
			defer conn.Close()

			key := fmt.Sprintf("key-%d", id%100)
			reader := bufio.NewReader(conn)
			for round := 0; round < stressRounds; round++ {
				if _, err := fmt.Fprintf(conn, "SET %s value\nGET %s\n", key, key); err != nil {
					errCh <- err
					return
				}
				for _, want := range []string{"OK\n", "value\n"} {
					got, err := reader.ReadString('\n')
					if err != nil {
						errCh <- err
						return
					}
					if got != want {
						errCh <- fmt.Errorf("client %d: response = %q, want %q", id, got, want)
						return
					}
				}
			}
		}(clientID)
	}
	clients.Wait()
	duration := time.Since(start)
	listener.Close()
	<-acceptDone
	handlers.Wait()
	close(errCh)

	for err := range errCh {
		if err != nil {
			t.Fatal(err)
		}
	}

	commands := stressClients * stressRounds * 2
	t.Logf("completed %d commands in %s (%.0f commands/s)", commands, duration.Round(time.Millisecond), float64(commands)/duration.Seconds())
}
