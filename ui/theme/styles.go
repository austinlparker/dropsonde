package theme

import (
	"github.com/charmbracelet/bubbles/table"
	"github.com/charmbracelet/lipgloss"
)

var theme = DefaultTheme()

var ActiveItemStyle = lipgloss.NewStyle().
	BorderStyle(lipgloss.NormalBorder()).
	BorderForeground(lipgloss.Color("240")).
	BorderLeft(true).
	PaddingLeft(1)

var ActiveTableStyle = lipgloss.NewStyle().
	BorderStyle(lipgloss.HiddenBorder()).
	BorderForeground(lipgloss.Color("240"))

var InactiveTableStyle = lipgloss.NewStyle().
	BorderStyle(lipgloss.NormalBorder()).
	BorderForeground(lipgloss.Color("240"))

var TitleTextStyle = lipgloss.NewStyle().
	Bold(true).
	Background(lipgloss.Color("55")).
	PaddingLeft(1).
	PaddingRight(1).
	MarginTop(1)

var ContainerStyle = lipgloss.NewStyle().PaddingLeft(1)

var TabSeparatorStyle = lipgloss.NewStyle().
	Foreground(theme.SecondaryFG)

var TabStyle = lipgloss.NewStyle().
	PaddingLeft(1).
	PaddingRight(1)

var ActiveTabStyle = TabStyle.Copy().
	Bold(true).
	Background(theme.Secondary).
	Foreground(theme.Primary)

var InactiveTabStyle = TabStyle.Copy().
	Bold(false).
	Foreground(theme.Secondary)

var TabGroupStyle = lipgloss.NewStyle().
	MarginRight(1).
	MarginLeft(1).
	PaddingBottom(1).
	BorderStyle(lipgloss.NormalBorder()).
	BorderForeground(theme.PrimaryFG).
	BorderBottom(true)

func GetTableStyle() table.Styles {
	s := table.DefaultStyles()
	s.Header = s.Header.
		BorderStyle(lipgloss.NormalBorder()).
		BorderForeground(theme.SecondaryFG).
		BorderTop(true).
		BorderBottom(true).
		Bold(true)
	s.Selected = s.Selected.
		Foreground(theme.PrimaryFG).
		Background(theme.Accent).
		Bold(false)

	return s
}
