package components

import (
	"github.com/austinlparker/dropsonde/ui/theme"
	"github.com/charmbracelet/bubbles/spinner"
	"github.com/charmbracelet/lipgloss"
)

var t = theme.DefaultTheme()

func NewSpinner() spinner.Model {
	s := spinner.New()
	s.Style = lipgloss.NewStyle().Foreground(t.SecondaryFG)
	s.Spinner = spinner.Points
	return s
}
