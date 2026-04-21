package main

import (
	"fmt"
	"os"

	"floco/internal/flake"
	"floco/internal/presets"
	"floco/internal/tui"

	tea "github.com/charmbracelet/bubbletea"
)

func main() {
	// 1. Carregamos os dados (Presets)
	allPresets, err := presets.LoadPresets()
	if err != nil {
		fmt.Printf("Erro ao carregar presets: %v\n", err)
		os.Exit(1)
	}

	// 2. Instanciamos o serviço de geração (A "Dependência")
	generator := flake.NewNixGenerator()

	// 3. Injetamos tudo no modelo inicial da TUI
	m := tui.InitialModel(allPresets, generator)

	p := tea.NewProgram(m, tea.WithAltScreen())
	if _, err := p.Run(); err != nil {
		fmt.Printf("Erro ao iniciar aplicação: %v\n", err)
		os.Exit(1)
	}
}
