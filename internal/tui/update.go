package tui

import (
	"strings"

	"floco/internal/envrc"
	"floco/internal/presets"

	tea "github.com/charmbracelet/bubbletea"
)

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	// Tecla global para sair
	if msg, ok := msg.(tea.KeyMsg); ok && msg.String() == "ctrl+c" {
		return m, tea.Quit
	}

	// Switch usando campos maiúsculos do Model
	switch m.State {
	case StateSplash:
		return updateSplash(m, msg)
	case StatePreset:
		return updatePreset(m, msg)
	case StateDeps:
		return updateDeps(m, msg)
	case StateManual:
		return updateManual(m, msg)
	case StatePreview:
		return updatePreview(m, msg)
	}
	return m, nil
}

func updateSplash(m Model, msg tea.Msg) (tea.Model, tea.Cmd) {
	if key, ok := msg.(tea.KeyMsg); ok && key.String() == "enter" {
		m.State = StatePreset
	}
	return m, nil
}

func updatePreset(m Model, msg tea.Msg) (tea.Model, tea.Cmd) {
	if key, ok := msg.(tea.KeyMsg); ok && key.String() == "enter" {
		// Captura o objeto Preset inteiro da lista
		if i, ok := m.PresetList.SelectedItem().(presets.Preset); ok {
			m.SelectedPreset = i

			// Reinicia as dependências com base no preset escolhido
			m.SelectedDeps = []string{}
			m.Checked = make(map[string]bool)

			for _, pkg := range i.Packages {
				m.SelectedDeps = append(m.SelectedDeps, pkg)
				m.Checked[pkg] = true
			}

			// Lógica de fluxo: Custom vai para seleção, Stacks prontas vão para Preview
			if i.Name == "Custom" {
				m.State = StateDeps
			} else {
				// Usa o FlakeGen (Maiúsculo) injetado
				m.Preview = m.FlakeGen.GenerateString(unique(m.SelectedDeps))
				m.State = StatePreview
			}
		}
	}
	var cmd tea.Cmd
	m.PresetList, cmd = m.PresetList.Update(msg)
	return m, cmd
}

func updateDeps(m Model, msg tea.Msg) (tea.Model, tea.Cmd) {
	if key, ok := msg.(tea.KeyMsg); ok {
		switch key.String() {
		case " ":
			if i, ok := m.DepsList.SelectedItem().(item); ok {
				name := string(i)
				m.Checked[name] = !m.Checked[name]

				if m.Checked[name] {
					m.SelectedDeps = append(m.SelectedDeps, name)
				} else {
					m.SelectedDeps = removeFromSlice(m.SelectedDeps, name)
				}
			}
		case "enter":
			m.State = StateManual
		}
	}
	var cmd tea.Cmd
	m.DepsList, cmd = m.DepsList.Update(msg)
	return m, cmd
}

func updateManual(m Model, msg tea.Msg) (tea.Model, tea.Cmd) {
	if key, ok := msg.(tea.KeyMsg); ok && key.String() == "enter" {
		val := strings.TrimSpace(m.ManualInput.Value())
		if val != "" {
			fields := strings.Fields(val)
			m.SelectedDeps = append(m.SelectedDeps, fields...)
		}

		m.Preview = m.FlakeGen.GenerateString(unique(m.SelectedDeps))
		m.State = StatePreview
	}
	var cmd tea.Cmd
	m.ManualInput, cmd = m.ManualInput.Update(msg)
	return m, cmd
}

func updatePreview(m Model, msg tea.Msg) (tea.Model, tea.Cmd) {
	if key, ok := msg.(tea.KeyMsg); ok {
		switch key.String() {
		case "enter":
			// Chama o Save com apenas UM argumento (o conteúdo)
			_ = m.FlakeGen.Save(m.Preview)
			_ = envrc.SaveEnvrc()
			return m, tea.Quit
		case "m":
			m.State = StateManual
			return m, nil
		case "esc":
			m.State = StatePreset
			return m, nil
		case "q":
			return m, tea.Quit
		}
	}
	return m, nil
}

// Helpers de utilidade permanecem no final do arquivo
func unique(list []string) []string {
	seen := make(map[string]bool)
	var out []string
	for _, v := range list {
		if v != "" && !seen[v] {
			seen[v] = true
			out = append(out, v)
		}
	}
	return out
}

func removeFromSlice(slice []string, val string) []string {
	for i, v := range slice {
		if v == val {
			return append(slice[:i], slice[i+1:]...)
		}
	}
	return slice
}
