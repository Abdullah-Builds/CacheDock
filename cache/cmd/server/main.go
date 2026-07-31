package main

import (
	"log"
	"net"
	"os"
	"os/signal"
	"path/filepath"
	"syscall"

	"github.com/Abdullah-Builds/cache/config"
	"github.com/Abdullah-Builds/cache/internal/cache"
	"github.com/Abdullah-Builds/cache/internal/handler"
)

func main() {

	// Load configuration
	config, err := config.LoadConfig()
	if err != nil {
		log.Fatal(err)
	}

	// Create data directory
	if err := os.MkdirAll(filepath.Dir(config.DataFile), 0755); err != nil {
		log.Fatal(err)
	}

	// Create cache
	cacheServer := cache.NewWithMaxKeys(config.MaxKeys)

	// Load previous snapshot
	if err := cacheServer.Load(config.DataFile); err != nil {
		log.Println("load error:", err)
	}

	// Start background jobs
	cacheServer.StartCleanup(config.CleanupInterval)
	cacheServer.StartAutoSave(config.DataFile, config.AutoSaveInterval)

	// Graceful shutdown
	signals := make(chan os.Signal, 1)
	signal.Notify(signals, os.Interrupt, syscall.SIGTERM)

	go func() {
		<-signals

		log.Println("Saving cache before shutdown...")

		if err := cacheServer.Save(config.DataFile); err != nil {
			log.Println("save error:", err)
		}

		os.Exit(0)
	}()

	// Start TCP server
	listener, err := net.Listen("tcp", ":"+config.Port)
	if err != nil {
		log.Fatal(err)
	}
	defer listener.Close()

	log.Printf("Cache server listening on :%s\n", config.Port)

	// Accept connections
	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Println("accept error:", err)
			continue
		}

		go handler.HandleConnection(conn, cacheServer)
	}
}
