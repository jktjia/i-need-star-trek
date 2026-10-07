package interactive

import (
	huh "charm.land/huh/v2"
	"charm.land/lipgloss/v2"
)

var (
	theme = huh.ThemeCatppuccin

	appStyle = lipgloss.NewStyle().Padding(1, 2)

	titleStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#FFFDF5")).
			Background(lipgloss.Color("#25A065")).
			Padding(0, 1)

	accent = lipgloss.Color("#04B575")

	statusMessageStyle = lipgloss.NewStyle().Foreground(accent).Render

	selectedStyle = lipgloss.NewStyle().Foreground(accent).BorderLeftForeground(accent)
	filterStyle   = lipgloss.NewStyle().Padding(1, 1)

	descStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("#4A4A4A"))

	highlight = lipgloss.Color("#7D56F4")
)
