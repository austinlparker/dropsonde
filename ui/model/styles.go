package model

import "github.com/charmbracelet/lipgloss"

const (
	OTelYellow     = "#F5A800"
	OTelBlue       = "#425CC7"
	DropsondeWhite = "#FFFFFF"
)

var (
	baseStyle = lipgloss.NewStyle().
			PaddingLeft(1).
			PaddingRight(1).
			Foreground(lipgloss.Color(DropsondeWhite))

	baseListStyle = lipgloss.NewStyle().
			PaddingTop(1).
			PaddingRight(2).
			PaddingLeft(1).
			PaddingBottom(1)

	stackListStyle = baseListStyle.Copy().
			Width(listWidth+5).
			Border(lipgloss.RoundedBorder(), true).
			BorderForeground(lipgloss.Color(OTelYellow))

	viewPortStyle = baseListStyle.Copy().
			Width(150).
			PaddingLeft(4).
			Border(lipgloss.RoundedBorder(), true).
			BorderForeground(lipgloss.Color(OTelYellow))

	modeStyle = baseStyle.Copy().
			Align(lipgloss.Center).
			Bold(true).
			Background(lipgloss.Color(OTelBlue))

	valueTitleStyle = baseStyle.Copy().
			Bold(true).
			Background(lipgloss.Color(OTelBlue)).
			Align(lipgloss.Left)

	helpMsgStyle = baseStyle.Copy().
			Bold(true).
			Background(lipgloss.Color(OTelBlue))

	// status bar styles

	statusBarStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color(DropsondeWhite)).
			Background(lipgloss.Color(OTelBlue))

	statusStyle = lipgloss.NewStyle().
			Inherit(statusBarStyle).
			Foreground(lipgloss.Color(DropsondeWhite)).
			Background(lipgloss.Color(OTelBlue)).
			Padding(0, 1).
			MarginRight(1)

	statusBarTextStyle = lipgloss.NewStyle().Inherit(statusBarStyle)
)
