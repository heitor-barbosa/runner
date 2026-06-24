package cmd

import (
	"fmt"
	"os"

	"github.com/heitor-barbosa/runner/projetos/simulador/internal/logging"
	"github.com/spf13/cobra"
)

var (
	verbose bool
	quiet   bool
)

var rootCmd = &cobra.Command{
	Use:     "simulador",
	Short:   "CLI para gerenciar o Simulador do HubSaude",
	Long:    "O CLI simulador gerencia o ciclo de vida do simulador.jar do Sistema Runner.",
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
