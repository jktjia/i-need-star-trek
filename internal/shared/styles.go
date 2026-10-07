package shared

import "github.com/charmbracelet/lipgloss"

var (
	AppStyle = lipgloss.NewStyle().Padding(1, 2)

	TitleStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#FFFDF5")).
			Background(lipgloss.Color("#25A065")).
			Padding(0, 1)

	accent = lipgloss.AdaptiveColor{Light: "#04B575", Dark: "#04B575"}

	statusMessageStyle = lipgloss.NewStyle().Foreground(accent).Render

	selectedStyle = lipgloss.NewStyle().Foreground(accent).BorderLeftForeground(accent)
	filterStyle   = lipgloss.NewStyle().Padding(1, 1)

	DescStyle = lipgloss.NewStyle().Foreground(lipgloss.AdaptiveColor{
		Light: "#B2B2B2",
		Dark:  "#4A4A4A",
	})
)
