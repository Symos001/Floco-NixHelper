package envrc

import (
	"fmt"
	"os"
)

func SaveEnvrc() error {
	if _, err := os.Stat(".envrc"); err == nil {
		fmt.Println("✔ .envrc já existe")
		return nil
	}

	if err := os.WriteFile(".envrc", []byte("use flake\n"), 0644); err != nil {
		return err
	}

	fmt.Println("✔ .envrc criado (lembre de rodar: direnv allow)")
	return nil
}
