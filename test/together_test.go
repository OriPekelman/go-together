package together_test

import (
	"testing"
	"time"

	"github.com/oripekelman/go-together/pkg/together"
)

func TestTogether(t *testing.T) {
	tg := together.NewTogether()

	// start a process that sleeps for a short time
	go tg.RunCmd("sleep 1")

	// wait a moment to ensure the process starts
	time.Sleep(time.Millisecond * 100)

	// ensure that the process is running
	if len(tg.Processes()) != 1 {
		t.Error("Expected 1 process to be running")
	}

	// wait for the process to finish
	time.Sleep(time.Second * 2)

	// ensure that the process has finished
	if len(tg.Processes()) != 0 {
		t.Error("Expected all processes to have finished")
	}
}
