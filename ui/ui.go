package ui

import (
	"fmt"
	"github.com/austinlparker/dropsonde/ui/model"
	"log"
	"os"

	tea "github.com/charmbracelet/bubbletea"
)

func RenderUI(tapEndpoint string) {
	if len(os.Getenv("DEBUG")) > 0 {
		f, err := tea.LogToFile("debug.log", "debug")
		if err != nil {
			fmt.Println("fatal:", err)
			os.Exit(1)
		}
		defer f.Close()
	}

	p := tea.NewProgram(model.NewModel(tapEndpoint), tea.WithAltScreen())
	if _, err := p.Run(); err != nil {
		log.Fatalf("error running program: %v", err)
	}
}
