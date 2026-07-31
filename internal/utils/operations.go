package utils

import (
	"net"
	"strings"

	"github.com/Abdullah-Builds/cache/internal/cache"
	"github.com/Abdullah-Builds/cache/internal/commands"
)

func PerformOperations(line string, cacheServer *cache.Cache, conn net.Conn) {
	line = strings.TrimSpace(line)

	if line == "" {
		return
	}

	parts := strings.Fields(line)
	command := strings.ToUpper(parts[0])

	handler, exists := commands.Registry[command]
	if !exists {
		conn.Write([]byte("ERROR unknown command\n"))
		return
	}

	response, ok := handler(parts, cacheServer, conn)

	if response != "" {
		conn.Write([]byte(response + "\n"))
	}

	if !ok {
		return
	}
}
