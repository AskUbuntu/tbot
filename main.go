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

	// Ensure that exactly two arguments were supplied
	// - the application itself
	// - the configuration file
	if len(os.Args) != 2 {
		log.Fatalf("Usage: %s CONFIG", os.Args[0])
	}

	// Load the configuration from disk
	log.Printf("Loading '%s'...\n", os.Args[1])
	config, err := config.Load(os.Args[1])
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("Data path: '%s'", config.DataPath)
	log.Printf("Root path: '%s'", config.RootPath)

	// Finally, create the server that will listen for requests
	log.Print("Initializing server...\n")
	server, err := server.New(config)
	if err != nil {
		log.Fatal(err)
	}
	defer server.Close()

	log.Println("Application is running")

	// Wait for SIGINT to be sent
	ch := make(chan os.Signal)
	signal.Notify(ch, syscall.SIGINT)
	<-ch

	log.Print("Shutting down...")
}
