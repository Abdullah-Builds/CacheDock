package commands

import (
	"net"

	"github.com/Abdullah-Builds/cache/internal/cache"
)

type Handler func(parts []string, cacheServer *cache.Cache, conn net.Conn) (string, bool)

var Registry = map[string]Handler{
	"SET":    Set,
	"GET":    Get,
	"DELETE": Delete,
	"PING":   Ping,
	"INFO":   Info,
}
