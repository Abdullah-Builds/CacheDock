package handler

import (
	"bufio"
	"io"
	"log"
	"net"
	"strings"

	"github.com/Abdullah-Builds/cache/internal/cache"
	"github.com/Abdullah-Builds/cache/internal/utils"
)

type Handler struct{}

func HandleConnection(conn net.Conn, cacheServer *cache.Cache) {
	defer conn.Close()

	log.Printf("Client connected: %s\n", conn.RemoteAddr())

	reader := bufio.NewReader(conn)

	for {
		line, err := reader.ReadString('\n')
		if err != nil {
			if err != io.EOF {
				log.Println("read error:", err)
			}
			break
		}

		log.Println("Received:", strings.TrimSpace(line))

		utils.PerformOperations(line, cacheServer, conn)
	}

	log.Printf("Client disconnected: %s\n", conn.RemoteAddr())
}
