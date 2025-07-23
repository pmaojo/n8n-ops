package cmd

import (
	"fmt"
	"os/exec"
	"runtime"
	"time"

	"github.com/pmaojo/n8n-ops/internal/ascii"
	"github.com/spf13/cobra"
)

var welcomeCmd = &cobra.Command{
	Use:   "welcome",
	Short: "Display the futuristic welcome screen",
	Long:  "Shows the complete n8n Operations Tool welcome screen with Matrix effects and ASCII art",
	Run:   runWelcome,
}

func init() {
	rootCmd.AddCommand(welcomeCmd)
}

func runWelcome(cmd *cobra.Command, args []string) {
	// Display the spectacular welcome screen first
	fmt.Print(ascii.WelcomeScreen())

	// Play robot voice announcement (blocking to ensure it plays)
	playRobotVoice("n8n workflow git based version system")
}

// playRobotVoice plays a robot voice saying the given text
func playRobotVoice(text string) {
	// Try different TTS systems based on the platform (blocking execution)
	var cmd *exec.Cmd

	switch runtime.GOOS {
	case "linux":
		// Use espeak with louder, clearer robot voice
		cmd = exec.Command("espeak", "-s", "150", "-p", "30", "-a", "200", "-v", "en+m3", text)
	case "darwin":
		// Use macOS say command with robot voice
		cmd = exec.Command("say", "-v", "Zarvox", text)
	case "windows":
		// Use PowerShell SAPI for Windows
		psCmd := fmt.Sprintf(`Add-Type -AssemblyName System.speech; $speak = New-Object System.Speech.Synthesis.SpeechSynthesizer; $speak.Rate = -2; $speak.Speak('%s')`, text)
		cmd = exec.Command("powershell", "-Command", psCmd)
	default:
		// Fallback - try espeak with audible settings
		cmd = exec.Command("espeak", "-s", "150", "-p", "30", "-a", "200", "-v", "en+m3", text)
	}

	// Execute the voice command (blocking to ensure it plays)
	if cmd != nil {
		cmd.Run() // Try to play audio
	}

	// Always show the robot text animation as visual feedback
	showRobotTextAnimation(text)
}

// showRobotTextAnimation displays animated robot text when voice is not available
func showRobotTextAnimation(text string) {
	fmt.Println("\n🤖 ROBOT VOICE SIMULATION:")
	fmt.Print("🔊 ")

	// Animate the text character by character like a robot typing
	for _, char := range text {
		fmt.Printf("%c", char)
		time.Sleep(80 * time.Millisecond) // Robot-like typing speed
	}

	fmt.Println("\n🎵 *BEEP BOOP BEEP*")
	time.Sleep(500 * time.Millisecond)
}
