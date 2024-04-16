package model

import (
	"fmt"

	"github.com/charmbracelet/lipgloss"
)

var (
	listWidth = 60
)

func (m model) View() string {
	var content string
	var msgsViewPtr string
	var valueViewPtr string
	var statusBar string
	var debugMsg string

	m.rawDataList.Title += msgsViewPtr

	var msgValueVP string
	if !m.valueVPReady {
		msgValueVP = "\n  Initializing..."
	} else {
		msgValueVP = viewPortStyle.Render(fmt.Sprintf("%s%s\n\n%s\n", valueTitleStyle.Render("Message Value"), valueViewPtr, m.valueVP.View()))
	}
	var helpVP string
	if !m.helpVPReady {
		helpVP = "\n  Initializing..."
	} else {
		helpVP = viewPortStyle.Render(fmt.Sprintf("  %s\n\n%s\n", valueTitleStyle.Width(m.terminalWidth).Render("Help"), m.helpVP.View()))
	}
	var opampVP string
	if !m.opampVPReady {
		opampVP = "\n  Initializing..."
	} else {
		opampVP = viewPortStyle.Render(fmt.Sprintf(" %s\n\n%s\n", valueTitleStyle.Width(m.terminalWidth).Render("OpAMP"), m.opampVP.View()))
	}

	switch m.vpFullScreen {
	case false:
		content = lipgloss.JoinHorizontal(
			lipgloss.Top,
			stackListStyle.Render(m.rawDataList.View()),
			msgValueVP,
		)
	case true:
		switch m.activeView {
		case ValueView:
			content = msgValueVP
		case HelpView:
			content = helpVP
		case OpAmpView:
			content = opampVP
		}
	}

	statusNameKey := statusStyle.Render("🔭 dropsonde")
	statusModeKey := statusBarTextStyle.Render(fmt.Sprintf("🖥️ mode: %v", m.activeView.Name()))
	statusBarHelpKey := statusBarTextStyle.Render(" ? for help")

	bar := lipgloss.JoinHorizontal(lipgloss.Top, statusNameKey, statusModeKey, statusBarHelpKey)
	statusBar = statusBarStyle.Width(m.terminalWidth).Render(bar)

	if m.debugMode {
		debugMsg += fmt.Sprintf(" %v", m.activeView)
	}

	return lipgloss.JoinVertical(lipgloss.Left,
		content,
		statusBar,
	)
}
