package entrypoint

import (
	"log"
	"time"

	"github.com/fsnotify/fsnotify"
)

func (ep *EntryPoint) setupWatcher() {
	var err error
	ep.watcher, err = fsnotify.NewWatcher()
	if err != nil {
		log.Fatalf("Failed to create file watcher: %v", err)
	}

	// Add all input files to the watcher
	for _, gen := range ep.appConfig.Generate {
		for _, input := range gen.Inputs {
			err = ep.watcher.Add(input.Path)
			if err != nil {
				log.Printf("Warning: Could not watch file %s: %v", input.Path, err)
			} else {
				log.Printf("Watching file: %s", input.Path)
			}
		}

		// If using template strategy, also watch the template file
		if gen.Strategy == "template" && gen.Template != "" {
			err = ep.watcher.Add(gen.Template)
			if err != nil {
				log.Printf("Warning: Could not watch template file %s: %v", gen.Template, err)
			} else {
				log.Printf("Watching template file: %s", gen.Template)
			}
		}
	}
}

func (ep *EntryPoint) WatchForChanges() {
	var lastEventTime time.Time
	debounceInterval := 100 * time.Millisecond

	for {
		select {
		case event, ok := <-ep.watcher.Events:
			if !ok {
				return
			}
			if event.Op&fsnotify.Write == fsnotify.Write {
				// Debounce events to prevent multiple rapid regenerations
				now := time.Now()
				if now.Sub(lastEventTime) > debounceInterval {
					lastEventTime = now
					log.Printf("File modified: %s", event.Name)

					// Find which config this file belongs to
					for _, gen := range ep.appConfig.Generate {
						needsRegeneration := false

						// Check if it's one of the input files
						for _, input := range gen.Inputs {
							if input.Path == event.Name {
								needsRegeneration = true
								break
							}
						}

						// Check if it's the template file
						if gen.Strategy == "template" && gen.Template == event.Name {
							needsRegeneration = true
						}

						if needsRegeneration {
							log.Printf("Regenerating output: %s", gen.Name)
							generateFile(gen)

							// If reload is enabled, reload the managed process
							if ep.appConfig.Process.Reload.Enabled {
								ep.reloadManagedProcess()
							}
						}
					}
				}
			}
		case err, ok := <-ep.watcher.Errors:
			if !ok {
				return
			}
			log.Printf("Watcher error: %v", err)
		}
	}
}
