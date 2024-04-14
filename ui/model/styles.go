package model

import "github.com/charmbracelet/lipgloss"

var (
	baseStyle = lipgloss.NewStyle().
			PaddingLeft(1).
			PaddingRight(1).
			Foreground(lipgloss.Color(OTEL_BLUE))

	baseListStyle = lipgloss.NewStyle().
			PaddingTop(1).
			PaddingRight(2).
			PaddingLeft(1).
			PaddingBottom(1)

	stackListStyle = baseListStyle.Copy().
			Width(listWidth+5).
			Border(lipgloss.RoundedBorder(), true).
			BorderForeground(lipgloss.Color("#3c3836"))

	viewPortStyle = baseListStyle.Copy().
			Width(150).
			PaddingLeft(4)

	modeStyle = baseStyle.Copy().
			Align(lipgloss.Center).
			Bold(true).
			Background(lipgloss.Color("#b8bb26"))

	kMsgValueTitleStyle = baseStyle.Copy().
				Bold(true).
				Background(lipgloss.Color("#8ec07c")).
				Align(lipgloss.Left)

	helpMsgStyle = baseStyle.Copy().
			Bold(true).
			Foreground(lipgloss.Color("#83a598"))
)
