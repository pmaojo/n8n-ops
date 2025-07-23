package ascii

import (
	"fmt"
	"math/rand"
	"strings"
	"time"
)

// Colors for terminal output
const (
	Reset     = "\033[0m"
	Cyan      = "\033[96m"
	Blue      = "\033[94m"
	Purple    = "\033[95m"
	Green     = "\033[92m"
	Yellow    = "\033[93m"
	Red       = "\033[91m"
	Bold      = "\033[1m"
	Dim       = "\033[2m"
	Underline = "\033[4m"
)

// N8nLogo returns the main n8n CLI futuristic logo
func N8nLogo() string {
	return fmt.Sprintf(`%s%s
                   ad88888ba                                                        
                  d8"     "8b                                                       
                  88       88                                                       
                  Y8a     a8P                                                       
                   "Y8aaa8P"                                                        
  ,ggg,,ggg,       ,d8"""8b,        ,ggg,,ggg,                                      
 ,8" "8P" "8,     d8"     "8b      ,8" "8P" "8,                                     
      8I   8I     88       88           8I   8I                                     
      8I   Yb,    Y8a     a8P           8I   Yb,                                    
      8I   Y8     "Y88888P"             8I   Y8       
                              
                                                                                    
         8I                                   ,dPYb,                                
         8I                                   IP' Yb                                
         8I                                   I8  8I                                
         8I                                   I8  8'                                
   ,gggg,8I       ,ggg,       gg,gggg,        I8 dP        ,ggggg,        gg     gg 
  dP"  "Y8I      i8" "8i      I8P"  "Yb       I8dP        dP"  "Y8ggg     I8     8I 
 i8'    ,8I      I8, ,8I      I8'    ,8i      I8P        i8'    ,8I       I8,   ,8I 
,d8,   ,d8b,     YbadP'     ,I8 _  ,d8'     ,d8b,_     ,d8,   ,d8'      ,d8b, ,d8I 
P"Y8888P" Y8    888P"Y888    PI8 YY88888P    8P'"Y88    P"Y8888P"        P""Y88P"888
                              I8                                               ,d8I'
                              I8                                             ,dP'8I 
                              I8                                            ,8"  8I 
                              I8                                            I8   8I 
                              I8                                            8, ,8I  
                              I8                                             Y8P"   
%s
             %s⚡ n8n Operations Tool ⚡%s`, Cyan, Bold, Reset, Yellow, Reset)
}

// SmallLogo returns a compact version for headers
func SmallLogo() string {
	return fmt.Sprintf(`%s%s
              __                                                
      /  |                / /                            
 ___ (___| ___       ___   (___       ___       ___  ___ 
|   )|   )|   )     |   )| |         |___ \   )|   )|    
|  / |__/ |  /      |__/ | |__        __/  \_/ |  / |__  
                    __/                     /                               
%s`, Blue, Bold, Reset)
}

// Banner creates a futuristic banner with environment info
func Banner(env string) string {
	envColor := getEnvironmentColor(env)
	return fmt.Sprintf(`%s%s
╔══════════════════════════════════════════════════════════════════════╗
║  %s⚡ N8N WORKFLOW OPERATIONS TOOL ⚡%s                                 ║
║                                                                      ║
║  %sEnvironment: %s%s%-12s%s                                            ║
║  %sMode: %sFuturistic Workflow Management%s                           ║
╚══════════════════════════════════════════════════════════════════════╝%s
`, Purple, Bold, Yellow, Purple, Dim, envColor, Bold, strings.ToUpper(env), Purple, Dim, Cyan, Purple, Reset)
}

// SuccessMessage creates an animated success display
func SuccessMessage(action string) string {
	return fmt.Sprintf(`%s%s
    ╔═══════════════════════════════════════════════════════════════════╗
    ║  %s✨ SUCCESS! ✨%s                                                  ║
    ║                                                                   ║
    ║  %s%s%s%s                                                           ║
    ╚═══════════════════════════════════════════════════════════════════╝%s
`, Green, Bold, Yellow, Green, Cyan, action, Green, Bold, Reset)
}

// ErrorMessage creates a futuristic error display
func ErrorMessage(err string) string {
	return fmt.Sprintf(`%s%s
    ╔═══════════════════════════════════════════════════════════════════╗
    ║  %s⚠️  ERROR DETECTED ⚠️%s                                           ║
    ║                                                                   ║
    ║  %s%s%s                                                             ║
    ╚═══════════════════════════════════════════════════════════════════╝%s
`, Red, Bold, Yellow, Red, Dim, err, Red, Reset)
}

// LoadingSpinner creates animated loading text
func LoadingSpinner(message string) string {
	spinners := []string{"⠋", "⠙", "⠹", "⠸", "⠼", "⠴", "⠦", "⠧", "⠇", "⠏"}
	rand.Seed(time.Now().UnixNano())
	spinner := spinners[rand.Intn(len(spinners))]

	return fmt.Sprintf(`%s%s%s %s%s%s`, Cyan, Bold, spinner, message, Dim, Reset)
}

// ProgressBar creates a futuristic progress bar
func ProgressBar(current, total int, message string) string {
	percentage := float64(current) / float64(total) * 100
	filled := int(percentage / 5) // 20 characters total
	bar := strings.Repeat("█", filled) + strings.Repeat("░", 20-filled)

	return fmt.Sprintf(`%s%s
    ╔═══════════════════════════════════════════════════════════════════╗
    ║  %s%s%s                                                             ║
    ║                                                                   ║
    ║  [%s%s%s] %s%.1f%%%s                                                ║
    ║  %s%d/%d workflows processed%s                                       ║
    ╚═══════════════════════════════════════════════════════════════════╝%s
`, Purple, Bold, Cyan, message, Purple, Green, bar, Purple, Yellow, percentage, Purple, Dim, current, total, Purple, Reset)
}

// WorkflowInfo displays workflow information in a futuristic format
func WorkflowInfo(name, status, env string) string {
	statusColor := getStatusColor(status)
	envColor := getEnvironmentColor(env)

	return fmt.Sprintf(`%s
    ┌─────────────────────────────────────────────────────────────────┐
    │  %sWorkflow: %s%-20s%s                                         │
    │  %sStatus:   %s%-20s%s                                         │
    │  %sEnv:      %s%-20s%s                                         │
    └─────────────────────────────────────────────────────────────────┘%s
`,
		Dim,
		Cyan, Bold, name, Dim,
		Cyan, statusColor, status, Dim,
		Cyan, envColor, env, Dim,
		Reset)
}

// CommandHelp creates futuristic help display
func CommandHelp(cmd string) string {
	return fmt.Sprintf(`%s%s
╔══════════════════════════════════════════════════════════════════════╗
║  %s🚀 COMMAND: %s%s%s                                                   ║
╚══════════════════════════════════════════════════════════════════════╝%s
`, Blue, Bold, Yellow, Cyan, strings.ToUpper(cmd), Blue, Reset)
}

// Matrix effect for fun
func MatrixEffect() string {
	chars := []string{"0", "1", "⚡", "▓", "░", "█", "▒"}
	var matrix strings.Builder

	for i := 0; i < 5; i++ {
		for j := 0; j < 50; j++ {
			char := chars[rand.Intn(len(chars))]
			matrix.WriteString(fmt.Sprintf("%s%s", Green, char))
		}
		matrix.WriteString(fmt.Sprintf("%s\n", Reset))
	}

	return matrix.String()
}

// Helper functions for colors
func getEnvironmentColor(env string) string {
	switch strings.ToLower(env) {
	case "production":
		return Red + Bold
	case "staging":
		return Yellow + Bold
	case "development":
		return Green + Bold
	default:
		return Cyan + Bold
	}
}

func getStatusColor(status string) string {
	switch strings.ToLower(status) {
	case "success", "completed", "active":
		return Green + Bold
	case "error", "failed":
		return Red + Bold
	case "warning", "pending":
		return Yellow + Bold
	default:
		return Cyan + Bold
	}
}

// WelcomeScreen creates the main futuristic welcome display
func WelcomeScreen() string {
	matrix := MatrixEffect()
	logo := N8nLogo()

	return fmt.Sprintf(`%s
%s
%s%s╔══════════════════════════════════════════════════════════════════════╗
║                                                                      ║
║  %sWelcome to the Future of Workflow Automation%s                      ║
║                                                                      ║
║  %s🔮 Powered by n8n Technology%s                                      ║
║  %s⚡ Lightning-fast Git Integration%s                                 ║
║  %s🚀 Multi-environment Support%s                                      ║
║  %s🛡️  Enterprise-grade Security%s                                     ║
║                                                                      ║
╚══════════════════════════════════════════════════════════════════════╝%s

%sType '%sn8n-ops --help%s' to begin your journey...%s
`, matrix, logo,
		Purple, Bold,
		Cyan, Purple,
		Green, Purple,
		Blue, Purple,
		Yellow, Purple,
		Red, Purple,
		Reset,
		Dim, Cyan+Bold, Dim, Reset)
}
