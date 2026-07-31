package commands

import (
	"fmt"
	"net"

	"github.com/Abdullah-Builds/cache/internal/cache"
)

func Info(parts []string, cacheServer *cache.Cache, conn net.Conn) (string, bool) {
	
	stats := cacheServer.Stats()

	msg := fmt.Sprintf(
		"Requests=%d Hits=%d Misses=%d Sets=%d Deletes=%d",
		stats.Requests,
		stats.Hits,
		stats.Misses,
		stats.Sets,
		stats.Deletes,
	)

	return msg, true
}
