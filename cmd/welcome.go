package cmd

import (
        "fmt"
        "os/exec"
        "runtime"

        "github.com/spf13/cobra"
        "github.com/n8n-workflows/n8n-ops/internal/ascii"
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
        // Play robot voice announcement
        playRobotVoice("n8n workflow git based version system")
        
        // Display the spectacular welcome screen
        fmt.Print(ascii.WelcomeScreen())
}

// playRobotVoice plays a robot voice saying the given text
func playRobotVoice(text string) {
        go func() {
                // Try different TTS systems based on the platform
                var cmd *exec.Cmd
                
                switch runtime.GOOS {
                case "linux":
                        // Use espeak with ultra robot-like parameters for maximum sci-fi effect
                        cmd = exec.Command("espeak", "-s", "120", "-p", "10", "-a", "90", "-v", "en+m3", "-k", "5", text)
                case "darwin":
                        // Use macOS say command with robot voice
                        cmd = exec.Command("say", "-v", "Zarvox", text)
                case "windows":
                        // Use PowerShell SAPI for Windows
                        psCmd := fmt.Sprintf(`Add-Type -AssemblyName System.speech; $speak = New-Object System.Speech.Synthesis.SpeechSynthesizer; $speak.Rate = -2; $speak.Speak('%s')`, text)
                        cmd = exec.Command("powershell", "-Command", psCmd)
                default:
                        // Fallback - try espeak with robot settings
                        cmd = exec.Command("espeak", "-s", "120", "-p", "10", "-a", "90", "-v", "en+m3", "-k", "5", text)
                }
                
                // Execute the voice command (ignore errors for better UX)
                if cmd != nil {
                        cmd.Run()
                }
        }()
}