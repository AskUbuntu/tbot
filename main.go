package main

import (
	"github.com/AskUbuntu/tbot/config"
	"github.com/AskUbuntu/tbot/server"

	"log"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	var (
		cfg *config.Config
		err error
	)

	// Ensure that exactly one or two arguments were supplied
	// - the application itself
	// - the configuration file [OPTIONAL]
	// With no config file supplied, os.Environ is used instead
	switch len(os.Args) {
	case 1:
		log.Print("loading configuration from environment")
		cfg = config.LoadFromEnvironment()
	case 2:
		log.Printf("loading configuration from '%s'...", os.Args[1])
		cfg, err = config.LoadFromDisk(os.Args[1])
		if err != nil {
			log.Fatal(err)
		}
	default:
		log.Fatalf("Usage: %s [CONFIG]", os.Args[0])
	}

	log.Printf("data path: '%s'", cfg.DataPath)
	log.Printf("root path: '%s'", cfg.RootPath)

	// Finally, create the server that will listen for requests
	log.Print("initializing server...")
	server, err := server.New(cfg)
	if err != nil {
		log.Fatal(err)
	}
	defer server.Close()

	log.Println("application is running")

	// Wait for SIGINT to be sent
	ch := make(chan os.Signal)
	signal.Notify(ch, syscall.SIGINT, syscall.SIGTERM)
	<-ch

	log.Print("shutting down...")
}
