// pkg/together/together.go

package together

import (
	"fmt"
	"os"
	"os/exec"
	"sync"
	"syscall"
)

type Together struct {
	processes map[int]*exec.Cmd
	mx        sync.Mutex
}

func NewTogether() *Together {
	return &Together{
		processes: make(map[int]*exec.Cmd),
	}
}

func (tg *Together) RunCmd(cmdStr string) {
	cmd := exec.Command("/bin/sh", "-c", cmdStr)
	cmd.SysProcAttr = &syscall.SysProcAttr{Setpgid: true}

	// Forward the output and error streams
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	err := cmd.Start()
	if err != nil {
		fmt.Printf("Failed to start %s: %s\n", cmdStr, err)
		return
	}

	tg.mx.Lock()
	tg.processes[cmd.Process.Pid] = cmd
	tg.mx.Unlock()

	fmt.Printf("Running %s as [%d]\n", cmdStr, cmd.Process.Pid)

	err = cmd.Wait()
	if err != nil {
		fmt.Printf("%s exited with error: %s\n", cmdStr, err)
	}

	tg.mx.Lock()
	delete(tg.processes, cmd.Process.Pid)
	tg.mx.Unlock()

	fmt.Printf("[%d] died\n", cmd.Process.Pid)

	if len(tg.processes) > 0 {
		fmt.Println("attempting to kill siblings")
		tg.KillAll()
	}
}

func (tg *Together) KillAll() {
	tg.mx.Lock()
	defer tg.mx.Unlock()

	for pid, cmd := range tg.processes {
		syscall.Kill(-cmd.Process.Pid, syscall.SIGKILL)
		delete(tg.processes, pid)
		fmt.Printf("killed [%d]\n", pid)
	}
	os.Exit(0)
}

func (tg *Together) Processes() map[int]*exec.Cmd {
	tg.mx.Lock()
	defer tg.mx.Unlock()

	// Create a new map to avoid data races
	processes := make(map[int]*exec.Cmd)
	for pid, cmd := range tg.processes {
		processes[pid] = cmd
	}
	return processes
}
