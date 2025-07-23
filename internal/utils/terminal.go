package utils

import (
	"os"
	"os/exec"
	"runtime"
)

// ClearTerminalScreen clears the terminal screen cross-platform
func ClearTerminalScreen() error {
	var cmd *exec.Cmd
	switch runtime.GOOS {
	case "windows":
		cmd = exec.Command("cmd", "/c", "cls")
	default:
		cmd = exec.Command("clear")
	}
	cmd.Stdout = os.Stdout
	return cmd.Run()
}
