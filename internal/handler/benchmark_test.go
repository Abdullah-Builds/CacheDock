package handler

import (
	"bufio"
	"io"
	"log"
	"net"
	"testing"

	"github.com/Abdullah-Builds/cache/internal/cache"
)

// BenchmarkProtocolRoundTripParallel exercises parsing, command dispatch, and
// responses with concurrent persistent client connections. net.Pipe keeps this
// focused on the application protocol rather than local TCP stack overhead.
func BenchmarkProtocolRoundTripParallel(b *testing.B) {
	previousOutput := log.Writer()
	log.SetOutput(io.Discard)
	b.Cleanup(func() { log.SetOutput(previousOutput) })

	c := cache.New()
	b.ReportAllocs()
	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		client, server := net.Pipe()
		go HandleConnection(server, c)
		defer client.Close()

		reader := bufio.NewReader(client)
		for pb.Next() {
			if _, err := client.Write([]byte("SET shared value\nGET shared\n")); err != nil {
				b.Fatal(err)
			}
			for _, want := range []string{"OK\n", "value\n"} {
				got, err := reader.ReadString('\n')
				if err != nil {
					b.Fatal(err)
				}
				if got != want {
					b.Fatalf("response = %q, want %q", got, want)
				}
			}
		}
	})
}
