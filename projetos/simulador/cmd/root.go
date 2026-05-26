package cmd

import (
	"os"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:     "simulador",
	Short:   "CLI para gerenciar o Simulador do HubSaude",
	Long:    "O CLI simulador gerencia o ciclo de vida do simulador.jar do Sistema Runner.",
	Version: Version,
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}
