package cmd

import (
	"fmt"
	"os"

	"github.com/heitor-barbosa/runner/projetos/assinador/internal/logging"
	"github.com/spf13/cobra"
)

var (
	verbose bool
	quiet   bool
)

var rootCmd = &cobra.Command{
	Use:     "assinatura",
	Short:   "CLI para o Sistema Runner de assinatura digital",
	Long:    "O Sistema Runner facilita o acesso a aplicacoes Java via linha de comandos.",
	Version: Version,
	PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
		if verbose && quiet {
			return fmt.Errorf("as flags --verbose e --quiet sao mutuamente exclusivas")
		}
		logging.Init(verbose, quiet)
		return nil
	},
}

func init() {
	rootCmd.PersistentFlags().BoolVarP(&verbose, "verbose", "v", false, "Exibe logs detalhados")
	rootCmd.PersistentFlags().BoolVarP(&quiet, "quiet", "q", false, "Exibe apenas erros")
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}
