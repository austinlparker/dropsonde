package components

import (
	"github.com/austinlparker/dropsonde/ui"
	"github.com/austinlparker/dropsonde/ui/theme"
	"github.com/charmbracelet/bubbles/help"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type HelpModel struct {
	help       help.Model
	keys       ui.KeyMap
	inputStyle lipgloss.Style
	quitting   bool
}

func NewHelpModel() HelpModel {
	return HelpModel{
		keys:       ui.KeyBindings,
		help:       help.New(),
		inputStyle: lipgloss.NewStyle().Foreground(theme.DefaultTheme().Primary),
	}
}

func (m HelpModel) Init() tea.Cmd {
	return nil
}

func (m HelpModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.help.Width = msg.Width
	}

	return m, nil
}

func (m HelpModel) View() string {
	if m.quitting {
		return "Quitting...\n"
	}

	style := lipgloss.NewStyle().MarginTop(1)
	return style.Render(m.help.View(m.keys))
}

func (m HelpModel) SetWidth(w int) {
	m.help.Width = w
}
