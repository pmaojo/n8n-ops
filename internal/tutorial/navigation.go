package tutorial

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/pmaojo/n8n-ops/internal/ascii"
)

func readKey(r io.Reader) string {
	scanner := bufio.NewScanner(r)
	if !scanner.Scan() {
		return ""
	}
	input := strings.ToLower(scanner.Text())
	switch input {
	case "q", "quit", "exit":
		return "q"
	case "up", "u":
		return "up"
	case "down", "d":
		return "down"
	default:
		return "enter"
	}
}

// WaitForKey reads a navigation key from stdin.
func WaitForKey() string {
	fmt.Print("> ")
	return readKey(os.Stdin)
}

// WaitForEnter blocks until the user presses enter.
func WaitForEnter() {
	fmt.Printf("\n%s%sPress Enter to continue...%s", ascii.Bold, ascii.Cyan, ascii.Reset)
	readKey(os.Stdin)
}
