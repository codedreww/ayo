package cmd

// root.go starts the Bubble Tea program with the app root model.
import (
	"ayo/internal/app"
	"fmt"
	"os"

	tea "github.com/charmbracelet/bubbletea"
)

func Execute() {
	p := tea.NewProgram(app.InitialModel(), tea.WithAltScreen())

	if _, err := p.Run(); err != nil {
		fmt.Println("Error:", err)
		os.Exit(1)
	}
}
