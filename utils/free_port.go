package utils

import (
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"time"
)

func FreePort(port string) {
	portNum := strings.TrimPrefix(port, ":")
	cmd := exec.Command("lsof", "-t", "-i:"+portNum)
	output, err := cmd.Output()
	if err != nil {
		return
	}

	currentPID := os.Getpid()

	pids := strings.Fields(string(output))

	for _, pidStr := range pids {
		pid, err := strconv.Atoi(pidStr)
		if err == nil && pid != currentPID {
			if proc, err := os.FindProcess(pid); err == nil {
				fmt.Printf("Freeing port %s (killing old process PID %d)...\n", portNum, pid)

				_ = proc.Kill()
				time.Sleep(200 * time.Millisecond)
			}
		}
	}
}
