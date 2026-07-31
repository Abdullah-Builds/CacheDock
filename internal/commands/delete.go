package commands

import (
	"net"

	"github.com/Abdullah-Builds/cache/internal/cache"
)

func Delete(parts []string, cacheServer *cache.Cache, conn net.Conn) (string, bool) {
	if len(parts) != 2 {
		return "ERROR usage: DELETE <key>", false
	}

	key := parts[1]

	cacheServer.Delete(key)
	return "OK", true
}
