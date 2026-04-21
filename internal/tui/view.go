package tui

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/lipgloss"
)

// Estilos centralizados para consistência visual
var (
	styleTitle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#00ADD8")).
			Bold(true).
			MarginBottom(1)

	styleSelected = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#FF00FF")).
			Bold(true)

	styleChecked = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#00FF00"))

	stylePreview = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(lipgloss.Color("#888888")).
			Padding(1).
			MarginTop(1)

	styleHelp = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#6272a4")).
			Italic(true)
)

func (m Model) View() string {
	var s string

	// IMPORTANTE: Usando M.State (maiúsculo) conforme o novo Model
	switch m.State {
	case StateSplash:
		s = fmt.Sprintf("%s\n\n%s",
			SplashScreen(),
			styleHelp.Render("Pressione Enter para iniciar..."))

	case StatePreset:
		// Renderiza a lista de Stacks (Presets do TOML)
		s = m.PresetList.View()

	case StateDeps:
		var b strings.Builder
		b.WriteString(styleTitle.Render("📦 Customizar Dependências") + "\n")

		// Renderiza a lista de pacotes sugeridos
		for i, it := range m.DepsList.Items() {
			name := string(it.(item))
			cursor := "  "
			lineStyle := lipgloss.NewStyle()

			if i == m.DepsList.Index() {
				cursor = styleSelected.Render("➤ ")
				lineStyle = styleSelected
			}

			box := "☐"
			if m.Checked[name] {
				box = styleChecked.Render("✔")
			}

			b.WriteString(fmt.Sprintf("%s%s %s\n", cursor, box, lineStyle.Render(name)))
		}

		b.WriteString("\n" + styleHelp.Render("espaço: marca • enter: continuar"))
		s = b.String()

	case StateManual:
		s = fmt.Sprintf(
			"%s\n\n%s\n\n%s",
			styleTitle.Render("⌨  Adicionar Dependências Extras"),
			m.ManualInput.View(),
			styleHelp.Render("Digite pacotes separados por espaço (ex: git vim) e pressione Enter"),
		)

	case StatePreview:
		// Cabeçalho dinâmico com o nome da Stack vindo do TOML
		header := styleTitle.Render("📄 Resumo da Stack: " + m.SelectedPreset.Name)
		previewBox := stylePreview.Render(m.Preview)
		footer := styleHelp.Render("Enter: Gerar Arquivos • M: Manual • Esc: Trocar Stack • Q: Sair")

		s = fmt.Sprintf("%s\n%s\n\n%s", header, previewBox, footer)
	}

	// Aplica a margem global em todos os estados
	return lipgloss.NewStyle().Margin(1, 2).Render(s)
}
