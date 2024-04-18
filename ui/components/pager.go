package components

import (
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type Pager struct {
	viewport    viewport.Model
	ready       bool
	width       int
	isDataReady bool
}

func NewPager() *Pager { return &Pager{} }

func (p *Pager) SetContent(msg tea.Msg) {
	s, ok := msg.(string)
	if !ok {
		return
	}

	w := lipgloss.Width(s)
	if w > p.width {
		p.viewport.SetContent("\n !! window size too small to show all data")
		return
	}
	p.viewport.SetContent(s)
	p.isDataReady = true
}

func (p *Pager) IsReady() bool { return p.isDataReady }

func (p *Pager) Init() tea.Cmd { return nil }

func (p *Pager) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var (
		cmd  tea.Cmd
		cmds []tea.Cmd
	)
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		p.width = msg.Width
		headerHeight := lipgloss.Height(Header())
		footerHeight := lipgloss.Height(NewHelpModel().View())
		verticalMarginHeight := headerHeight + footerHeight
		if !p.ready {
			p.viewport = viewport.New(msg.Width, msg.Height-verticalMarginHeight)
			p.viewport.YPosition = headerHeight
			p.ready = true
		} else {
			p.viewport.Width = msg.Width
			p.viewport.Height = msg.Height - verticalMarginHeight
		}
	}

	// Handle keyboard and mouse events in the viewport
	p.viewport, cmd = p.viewport.Update(msg)
	cmds = append(cmds, cmd)

	return p, tea.Batch(cmds...)
}
func (p *Pager) View() string {
	return p.viewport.View()
}

func (p *Pager) SetUnready() { p.isDataReady = false }
