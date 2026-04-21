package tui

import (
	"floco/internal/presets"
	"fmt"
	"io"

	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
)

type presetDelegate struct{}

func (d presetDelegate) Height() int                               { return 1 }
func (d presetDelegate) Spacing() int                              { return 0 }
func (d presetDelegate) Update(msg tea.Msg, m *list.Model) tea.Cmd { return nil }

func (d presetDelegate) Render(w io.Writer, m list.Model, index int, listItem list.Item) {
	// Aqui está o pulo do gato:
	p, ok := listItem.(presets.Preset)
	if !ok {
		return
	}

	str := p.Name // Pegamos apenas o nome da Stack para o menu

	// Estilização simples para o item selecionado
	if index == m.Index() {
		fmt.Fprintf(w, "➤ %s", styleSelected.Render(str))
	} else {
		fmt.Fprintf(w, "  %s", str)
	}
}
