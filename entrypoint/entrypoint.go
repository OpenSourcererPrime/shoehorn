package entrypoint

import (
	"log"
	"os"
	"os/exec"
	"os/signal"
	"syscall"
	"time"

	"github.com/OpenSourcererPrime/shoehorn/config"
	"github.com/fsnotify/fsnotify"
)

type EntryPoint struct {
	managedCmd *exec.Cmd
	appConfig  config.Config
	watcher    *fsnotify.Watcher
}

func NewEntryPoint(appConfig *config.Config) (*EntryPoint, error) {
	ep := &EntryPoint{
		appConfig: *appConfig,
	}

	// Setup file watcher
	ep.setupWatcher()

	// Generate initial files
	ep.generateAllFiles()

	return ep, nil
}

func (ep *EntryPoint) Close() {
	if ep.watcher != nil {
		ep.watcher.Close()
	}
	if ep.managedCmd != nil && ep.managedCmd.Process != nil {
		ep.managedCmd.Process.Kill()
	}
}

func (ep *EntryPoint) HandleSignals() {
	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, syscall.SIGINT, syscall.SIGTERM)

	sig := <-signalChan
	log.Printf("Received signal: %v, shutting down...", sig)

	// Forward the signal to the managed process if it exists
	if ep.managedCmd != nil && ep.managedCmd.Process != nil {
		log.Printf("Forwarding signal to managed process")
		ep.managedCmd.Process.Signal(sig)

		// Give the process a short time to exit gracefully
		done := make(chan error, 1)
		go func() {
			done <- ep.managedCmd.Wait()
		}()

		select {
		case <-done:
			log.Printf("Managed process exited gracefully")
		case <-time.After(5 * time.Second):
			log.Printf("Timeout waiting for managed process to exit, forcing termination")
			ep.managedCmd.Process.Kill()
		}
	}

	os.Exit(0)
}
