package main

import (
    "flag"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/rdawson46/pic-conversion/internal/server"
)

func run() error {
    modePtr := flag.String("mode", "test", "mode")

    var config server.ServerConfig
    if *modePtr == "prod" {
        config = server.NewConfig(8000, 2000, 10, 30, 10, server.Prod)
    } else {
        config = server.NewConfig(8000, 2000, 10, 30, 10, server.Test)
    }

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
