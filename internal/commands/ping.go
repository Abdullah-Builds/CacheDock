package commands

import (
	"net"

	"github.com/Abdullah-Builds/cache/internal/cache"
)

func Ping(parts []string, cacheServer *cache.Cache, conn net.Conn) (string, bool) {
	
	return "PONG", true
}
