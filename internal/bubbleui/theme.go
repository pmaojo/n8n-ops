package bubbleui

import "github.com/charmbracelet/lipgloss"

// Theme defines visual properties for dashboard components.
type Theme struct {
	SuccessColor   lipgloss.Color
	ErrorColor     lipgloss.Color
	WarningColor   lipgloss.Color
	NeutralColor   lipgloss.Color
	GaugeStyle     lipgloss.Style
	GaugeBarStyle  lipgloss.Style
	HighlightStyle lipgloss.Style
}

// DefaultTheme provides colors similar to the previous hard-coded scheme.
var DefaultTheme = Theme{
	SuccessColor:   lipgloss.Color("2"),
	ErrorColor:     lipgloss.Color("1"),
	WarningColor:   lipgloss.Color("3"),
	NeutralColor:   lipgloss.Color("7"),
	GaugeStyle:     lipgloss.NewStyle().Foreground(lipgloss.Color("2")),
	GaugeBarStyle:  lipgloss.NewStyle().Foreground(lipgloss.Color("10")),
	HighlightStyle: lipgloss.NewStyle().Foreground(lipgloss.Color("0")).Background(lipgloss.Color("2")),
}

// CyberpunkTheme is a neon-styled variant.
var CyberpunkTheme = Theme{
	SuccessColor:   lipgloss.Color("#39FF14"),
	ErrorColor:     lipgloss.Color("#FF2079"),
	WarningColor:   lipgloss.Color("#F5F543"),
	NeutralColor:   lipgloss.Color("7"),
	GaugeStyle:     lipgloss.NewStyle().Foreground(lipgloss.Color("#39FF14")),
	GaugeBarStyle:  lipgloss.NewStyle().Foreground(lipgloss.Color("#FF007A")),
	HighlightStyle: lipgloss.NewStyle().Foreground(lipgloss.Color("#000000")).Background(lipgloss.Color("#39FF14")),
}
