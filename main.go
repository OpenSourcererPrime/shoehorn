package main

import (
	"log"
	"os"

	"github.com/OpenSourcererPrime/shoehorn/config"
	"github.com/OpenSourcererPrime/shoehorn/entrypoint"
)

func main() {
	if len(os.Args) < 2 {
		log.Fatal("Usage: shoehorn <shoehorn.yaml> [args for managed process...]")
	}

	configPath := os.Args[1]
	extraArgs := os.Args[2:]

	log.Printf("Starting shoehorn with config: %s", configPath)

	r, err := os.Open(configPath)
	if err != nil {
		log.Fatalf("Failed to open config file: %v", err)
	}

	defer r.Close()
	appConfig, err := config.LoadConfig(r)
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}
	appConfig.Process.Args = append(appConfig.Process.Args, extraArgs...)

	ep, err := entrypoint.NewEntryPoint(appConfig)
	if err != nil {
		log.Fatalf("Failed to create entrypoint: %v", err)
	}
	defer ep.Close()

	// Start the managed process
	ep.StartManagedProcess()

	// Watch for changes in a separate goroutine
	go ep.WatchForChanges()

	// Handle signals for graceful termination
	ep.HandleSignals()
}
