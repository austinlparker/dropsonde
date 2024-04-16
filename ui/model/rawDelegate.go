package model

import (
	"github.com/austinlparker/dropsonde/internal/otlptext"
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"go.opentelemetry.io/collector/pdata/plog"
	"go.opentelemetry.io/collector/pdata/pmetric"
	"go.opentelemetry.io/collector/pdata/ptrace"
)

func newRawDelegateKeyMap() *rawDelegateKeyMap {
	return &rawDelegateKeyMap{
		details: key.NewBinding(
			key.WithKeys("enter"),
			key.WithHelp("enter", "show details"),
		),
	}
}

func newRawItemDelegate(keys *rawDelegateKeyMap) list.DefaultDelegate {
	d := list.NewDefaultDelegate()

	lm := otlptext.NewTextLogsMarshaler()
	mm := otlptext.NewTextMetricsMarshaler()
	tm := otlptext.NewTextTracesMarshaler()

	d.Styles.SelectedTitle = d.Styles.
		SelectedTitle.
		Foreground(lipgloss.Color(OTelBlue)).
		BorderLeftForeground(lipgloss.Color(OTelBlue))
	d.Styles.SelectedDesc = d.Styles.
		SelectedTitle.
		Copy()

	d.UpdateFunc = func(msg tea.Msg, m *list.Model) tea.Cmd {
		switch msgType := msg.(type) {
		case tea.KeyMsg:
			switch {
			case key.Matches(msgType,
				keys.details,
				list.DefaultKeyMap().CursorUp,
				list.DefaultKeyMap().CursorDown,
				list.DefaultKeyMap().GoToStart,
				list.DefaultKeyMap().GoToEnd,
				list.DefaultKeyMap().NextPage,
				list.DefaultKeyMap().PrevPage,
			):
				if item, ok := m.SelectedItem().(RawDataItem); ok {
					var res []byte
					switch item.itemType {
					case MetricData:
						unmarshaler := &pmetric.JSONUnmarshaler{}
						metrics, err := unmarshaler.UnmarshalMetrics([]byte(item.item))
						if err != nil {
							return nil
						}
						res, err = mm.MarshalMetrics(metrics)
					case LogData:
						unmarshaler := &plog.JSONUnmarshaler{}
						logs, err := unmarshaler.UnmarshalLogs([]byte(item.item))
						if err != nil {
							return nil
						}
						res, err = lm.MarshalLogs(logs)
					case TraceData:
						unmarshaler := &ptrace.JSONUnmarshaler{}
						traces, err := unmarshaler.UnmarshalTraces([]byte(item.item))
						if err != nil {
							return nil
						}
						res, err = tm.MarshalTraces(traces)
					}
					return showRawData(string(res))
				} else {
					return nil
				}
			}
		}
		return nil
	}
	help := []key.Binding{keys.details}

	d.ShortHelpFunc = func() []key.Binding {
		return help
	}

	d.FullHelpFunc = func() [][]key.Binding {
		return [][]key.Binding{help}
	}
	return d

}
