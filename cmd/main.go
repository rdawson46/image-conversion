package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/rdawson46/pic-conversion/internal/server"
)

func run() error {
    config := server.NewConfig(8000, 10, 10, 30, 10, server.Test)

    s := server.NewServer(config)


    if err := s.Start(); err != nil {
        log.Fatal(err)
    }

    quit := make(chan os.Signal, 1)
    signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
    <- quit

    return s.Shutdown()
}

func main() {
    if err := run(); err != nil {
        log.Fatal(err)
    }
}
