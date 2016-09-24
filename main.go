package main

import (
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
	config, err := LoadConfig(os.Args[1])
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("Data path: '%s'", config.DataPath)
	log.Printf("Root path: '%s'", config.RootPath)

	// Initialize the user-customizable settings
	log.Print("Loading settings...\n")
	settings, err := NewSettings(config)
	if err != nil {
		log.Fatal(err)
	}
	defer settings.Close()

	// Initialize the list of users
	log.Print("Loading users...\n")
	auth, err := NewAuth(config)
	if err != nil {
		log.Fatal(err)
	}

	// Create and initialize the queue
	log.Print("Loading queue...\n")
	queue, err := NewQueue(config, settings)
	if err != nil {
		log.Fatal(err)
	}
	defer queue.Close()

	// Create the Twitter client
	log.Print("Initializing Twitter client...\n")
	client, err := NewClient(config, queue.Tweet)
	if err != nil {
		log.Fatal(err)
	}
	defer client.Close()

	// Finally, create the server that will listen for requests
	log.Print("Initializing server...\n")
	server, err := NewServer(config, settings, auth, queue)
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
