package flake

import (
	"os"
	"strings"
)

// Generator define o que um gerador de flake deve fazer
type Generator interface {
	GenerateString(deps []string) string
	Save(content string) error
}

// NixGenerator é a implementação concreta que a TUI vai usar
type NixGenerator struct{}

// NewNixGenerator cria uma nova instância do gerador
func NewNixGenerator() Generator {
	return &NixGenerator{}
}

// GenerateString é a sua função BuildFlake convertida em método
func (n *NixGenerator) GenerateString(deps []string) string {
	var b strings.Builder

	b.WriteString("{\n")
	b.WriteString("  description = \"Flake gerado pelo Floco\";\n\n")
	b.WriteString("  inputs = {\n")
	b.WriteString("    nixpkgs.url = \"github:NixOS/nixpkgs/nixos-unstable\";\n")
	b.WriteString("    flake-utils.url = \"github:numtide/flake-utils\";\n")
	b.WriteString("  };\n\n")
	b.WriteString("  outputs = { self, nixpkgs, flake-utils }:\n")
	b.WriteString("    flake-utils.lib.eachDefaultSystem (system:\n")
	b.WriteString("      let\n")
	b.WriteString("        pkgs = import nixpkgs { inherit system; };\n")
	b.WriteString("      in {\n")
	b.WriteString("        devShells.default = pkgs.mkShell {\n")
	b.WriteString("          buildInputs = [\n")

	for _, d := range deps {
		if d != "" {
			b.WriteString("            pkgs." + strings.TrimPrefix(d, "pkgs.") + "\n")
		}
	}

	b.WriteString("          ];\n\n")
	b.WriteString("          shellHook = ''\n")
	b.WriteString("            echo \"❄️ Floco: ambiente Nix carregado\"\n")
	b.WriteString("            export PNPM_HOME=\"$PWD/.pnpm\"\n")
	b.WriteString("            export PATH=\"$PNPM_HOME:$PATH\"\n")
	b.WriteString("            mkdir -p \"$PNPM_HOME\"\n")
	b.WriteString("          '';\n")
	b.WriteString("        };\n")
	b.WriteString("      });\n")
	b.WriteString("}\n")

	return b.String()
}

// Save é a sua função SaveFlake convertida em método
func (n *NixGenerator) Save(content string) error {
	return os.WriteFile("flake.nix", []byte(content), 0644)
}
