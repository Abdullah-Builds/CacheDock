package commands

import (
	"net"
	"strconv"
	"time"

	"github.com/Abdullah-Builds/cache/internal/cache"
)

func Set(parts []string, cacheServer *cache.Cache, conn net.Conn) (string, bool) {

	var ttl time.Duration
	if len(parts) != 3 && len(parts) != 4 {
		return "ERROR usage: SET <key> <value> [ttl_seconds]", false
	}

	if len(parts) == 4 {
		seconds, err := strconv.Atoi(parts[3])
		if err != nil {
			return "ERROR invalid TTL", false
		}

		ttl = time.Duration(seconds) * time.Second
	}

	key := parts[1]
	value := parts[2]

	msg := cacheServer.Set(key, value, ttl)
	return msg, true
}
