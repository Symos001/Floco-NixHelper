package tui

import (
	"floco/internal/deps"
	"floco/internal/flake"
	"floco/internal/presets"

	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
)

type state int

const (
	StateSplash state = iota
	StatePreset
	StateDeps
	StateManual
	StatePreview
)

type Model struct {
	State          state           // Exportado para o Update acessar
	PresetList     list.Model      // Exportado
	DepsList       list.Model      // Exportado
	ManualInput    textinput.Model // Exportado
	SelectedDeps   []string
	SelectedPreset presets.Preset  // Agora guardamos a struct completa do TOML
	FlakeGen       flake.Generator // Interface injetada
	Preview        string
	Checked        map[string]bool
}

// InitialModel recebe os presets (do TOML) e o gerador (DI)
func InitialModel(ps []presets.Preset, gen flake.Generator) Model {
	// 1. Configuração da Lista de Presets (Stacks)
	// Transformamos os presets do TOML em itens da lista
	presetItems := make([]list.Item, len(ps))
	for i, p := range ps {
		presetItems[i] = p
	}

	// Criamos a lista de Presets usando seu delegate customizado
	pl := list.New(presetItems, presetDelegate{}, 40, 15)
	pl.Title = "Stacks Disponíveis"

	// 2. Configuração da Lista de Dependências Sugeridas
	// Usamos o tipo 'item' (string) para as sugestões extras
	depItems := make([]list.Item, len(deps.Suggested))
	for i, d := range deps.Suggested {
		depItems[i] = item(d)
	}

	dl := list.New(depItems, list.NewDefaultDelegate(), 40, 15)
	dl.Title = "Dependências Sugeridas"

	// 3. Configuração do Input Manual
	ti := textinput.New()
	ti.Placeholder = "Ex: pkgs.hello"
	ti.Width = 40

	return Model{
		State:       StateSplash,
		PresetList:  pl,
		DepsList:    dl,
		ManualInput: ti,
		FlakeGen:    gen,
		Checked:     make(map[string]bool),
	}
}

// item é usado apenas para a lista de dependências sugeridas (strings simples)
type item string

func (i item) FilterValue() string { return string(i) }
func (i item) Title() string       { return string(i) }
func (i item) Description() string { return "" }

func (m Model) Init() tea.Cmd {
	return nil
}
