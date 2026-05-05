package cmd

import (
	"os"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:     "assinatura",
	Short:   "CLI para o Sistema Runner de assinatura digital",
	Long:    "O Sistema Runner facilita o acesso a aplicacoes Java via linha de comandos.",
	Version: Version,
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}
