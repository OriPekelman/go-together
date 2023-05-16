// cmd/main.go

package main

import (
	"fmt"
	"os"
	"os/signal"
	"sync"
	"syscall"

	"github.com/oripekelman/go-together/pkg/together"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println(`together can run multiple processes, together, such that they live and die together.
  
Usage: 
together "sleep 5" "sleep 10"
If any of the processes die, together will kill the others
If it receives a SIGTERM or a SIGINT it will kill the spawned processes`)
		os.Exit(0)
	}

	// Initialize Together
	tg := together.NewTogether()

	// Set up signal handling
	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT, syscall.SIGHUP)
	go func() {
		sig := <-sigCh
		fmt.Printf("Received %s, killing all processes\n", sig)
		tg.KillAll()
	}()

	// Spawn child processes
	var wg sync.WaitGroup
	for _, cmdStr := range os.Args[1:] {
		wg.Add(1)
		go func(cmdStr string) {
			defer wg.Done()
			tg.RunCmd(cmdStr)
		}(cmdStr)
	}

	// Wait for all child processes to finish
	wg.Wait()
}
