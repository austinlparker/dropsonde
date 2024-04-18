package components

import (
	"github.com/austinlparker/dropsonde/ui/theme"
	"github.com/charmbracelet/lipgloss"
)

func Header() string {
	return lipgloss.NewStyle().
		Bold(true).
		Background(theme.DefaultTheme().PrimaryBG).
		Foreground(theme.DefaultTheme().PrimaryFG).
		MarginTop(1).
		MarginBottom(1).
		PaddingLeft(7).
		PaddingRight(7).
		Render("dropsonde")
}
