package presets

import (
	"embed"
	"os"
	"path/filepath"

	"github.com/BurntSushi/toml"
)

//go:embed default_presets.toml
var embeddedFS embed.FS

// Preset representa uma stack de desenvolvimento
type Preset struct {
	Name     string   `toml:"name"`
	Packages []string `toml:"packages"`
}

// Config é a estrutura raiz do arquivo TOML
type Config struct {
	Presets []Preset `toml:"presets"`
}

// Métodos para satisfazer a interface list.Item do BubbleTea
func (p Preset) FilterValue() string { return p.Name }
func (p Preset) Title() string       { return p.Name }
func (p Preset) Description() string {
	if len(p.Packages) == 0 {
		return "Configure suas próprias dependências"
	}
	return "Packages: " + p.Name
}

// LoadPresets carrega os presets do disco (se existirem) ou do binário embutido
func LoadPresets() ([]Preset, error) {
	var defaultConfig Config
	var customConfig Config

	// 1. Carregar sempre o embutido primeiro (Garante que Stacks novas como Flutter apareçam)
	data, err := embeddedFS.ReadFile("default_presets.toml")
	if err == nil {
		toml.Decode(string(data), &defaultConfig)
	}

	// 2. Tentar carregar o arquivo do usuário
	home, _ := os.UserHomeDir()
	customPath := filepath.Join(home, ".config", "floco", "presets.toml")

	if _, err := os.Stat(customPath); err == nil {
		if _, err := toml.DecodeFile(customPath, &customConfig); err == nil {
			// Adiciona os presets customizados aos padrões
			return append(defaultConfig.Presets, customConfig.Presets...), nil
		}
	}

	return defaultConfig.Presets, nil
}
