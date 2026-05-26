package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var startPort int

var startCmd = &cobra.Command{
	Use:   "start",
	Short: "Inicia o simulador.jar",
	Long: `Inicia o simulador.jar em background.

Esta estrutura inicial registra o comando no CLI. A inicializacao real do processo,
validacao de portas e download automatico serao implementados nas proximas historias.`,
	Run: runStart,
}

func init() {
	rootCmd.AddCommand(startCmd)
	startCmd.Flags().IntVar(&startPort, "port", 8081, "Porta do Simulador do HubSaude")
}

func runStart(cmd *cobra.Command, args []string) {
	fmt.Fprintf(cmd.OutOrStdout(), "Comando simulador start registrado para a porta %d\n", startPort)
}
