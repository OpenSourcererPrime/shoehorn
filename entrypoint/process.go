package entrypoint

import (
	"log"
	"os"
	"os/exec"
	"strings"
	"syscall"
)

func (ep *EntryPoint) StartManagedProcess() {
	if ep.appConfig.Process.Path == "" {
		log.Printf("No process specified to manage, entrypoint will only manage configurations")
		return
	}

	log.Printf("Starting managed process: %s %s", ep.appConfig.Process.Path, strings.Join(ep.appConfig.Process.Args, " "))

	c := exec.Command(ep.appConfig.Process.Path, ep.appConfig.Process.Args...)

	ep.managedCmd = c

	// Connect process stdin/stdout/stderr to the entrypoint's
	ep.managedCmd.Stdin = os.Stdin
	ep.managedCmd.Stdout = os.Stdout
	ep.managedCmd.Stderr = os.Stderr

	err := ep.managedCmd.Start()
	if err != nil {
		log.Fatalf("Failed to start managed process: %v", err)
	}

	// Handle process completion in a goroutine
	go func() {
		err := ep.managedCmd.Wait()
		if err != nil {
			log.Printf("Inside goroutine: managed process exited with error: %v", err)
			log.Printf("Managed process exited with error: %v", err)
			os.Exit(1)
		} else {
			log.Printf("Managed process completed successfully")
			os.Exit(0)
		}
	}()
}

func (ep *EntryPoint) reloadManagedProcess() {
	if ep.managedCmd == nil || ep.managedCmd.Process == nil {
		log.Printf("No managed process to reload")
		return
	}

	switch ep.appConfig.Process.Reload.Method {
	case "restart":
		log.Printf("Restarting managed process")
		// Kill the existing process
		if err := ep.managedCmd.Process.Signal(syscall.SIGTERM); err != nil {
			log.Printf("Failed to kill managed process: %v", err)
		}

		// Wait for it to exit (should be near-immediate since we killed it)
		ep.managedCmd.Wait()
		ep.managedCmd.Process.Wait()

		// Start it again with the same arguments
		ep.StartManagedProcess()

	case "signal":
		log.Printf("Sending %s to managed process", ep.appConfig.Process.Reload.Signal)
		var sig syscall.Signal

		// Convert string signal name to actual signal
		switch ep.appConfig.Process.Reload.Signal {
		case "SIGHUP":
			sig = syscall.SIGHUP
		case "SIGUSR1":
			sig = syscall.SIGUSR1
		case "SIGUSR2":
			sig = syscall.SIGUSR2
		case "SIGTERM":
			sig = syscall.SIGTERM
		case "SIGINT":
			sig = syscall.SIGINT
		default:
			log.Printf("Unsupported signal: %s, using SIGHUP instead", ep.appConfig.Process.Reload.Signal)
			sig = syscall.SIGHUP
		}

		err := ep.managedCmd.Process.Signal(sig)
		if err != nil {
			log.Printf("Failed to send signal to managed process: %v", err)
		}
	}
}
