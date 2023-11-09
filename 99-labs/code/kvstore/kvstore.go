package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"path/filepath"
	"syscall"

	"kvstore/pkg/server"
)

func main() {
	// Include date, time and filename in the log messages
	log.SetFlags(log.LstdFlags | log.Lshortfile)

	logFile := filepath.Join(os.TempDir(), "translog.log")
	server, err := server.NewServer(logFile)
	if err != nil {
		log.Fatalf("could not init server: %s", err.Error())
	}

	ctx, cancel := context.WithCancel(context.Background())

	if err = server.Run(ctx, ":8081"); err != nil {
		log.Fatalf("could not start server: %s", err.Error())
	}

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	<-sigChan

	cancel()
}
