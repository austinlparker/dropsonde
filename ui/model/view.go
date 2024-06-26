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
	var footer string
	var msgsViewPtr string
	var valueViewPtr string
	var mode string
	var statusBar string
	var debugMsg string

	if m.msg != "" {
		statusBar = Trim(m.msg, 120)
	}

	m.rawDataList.Title += msgsViewPtr

	var errorMsg string
	if m.errorMsg != "" {
		errorMsg = " error: " + Trim(m.errorMsg, 120)
	}

	var msgValueVP string
	if !m.valueVPReady {
		msgValueVP = "\n  Initializing..."
	} else {
		msgValueVP = viewPortStyle.Render(fmt.Sprintf("%s%s\n\n%s\n", kMsgValueTitleStyle.Render("Message Value"), valueViewPtr, m.valueVP.View()))
	}
	var helpVP string
	if !m.helpVPReady {
		helpVP = "\n  Initializing..."
	} else {
		helpVP = viewPortStyle.Render(fmt.Sprintf("  %s\n\n%s\n", kMsgValueTitleStyle.Render("Help"), m.helpVP.View()))
	}
	var opampVP string
	if !m.opampVPReady {
		opampVP = "\n  Initializing..."
	} else {
		opampVP = viewPortStyle.Render(fmt.Sprintf(" %s\n\n%s\n", kMsgValueTitleStyle.Render("OpAMP"), m.opampVP.View()))
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

	footerStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color(OTEL_BLUE)).
		Background(lipgloss.Color(OTEL_YELLOW))

	var helpMsg string
	if m.showHelpIndicator {
		helpMsg = " " + helpMsgStyle.Render("Press ? for help")
	}

	if m.debugMode {
		debugMsg += fmt.Sprintf(" %v", m.activeView)
	}

	footerStr := fmt.Sprintf("%s%s%s%s%s",
		modeStyle.Render("dropsonde"),
		debugMsg,
		helpMsg,
		mode,
		errorMsg,
	)
	footer = footerStyle.Render(footerStr)

	return lipgloss.JoinVertical(lipgloss.Left,
		content,
		statusBar,
		footer,
	)
}
