package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var stopPort int

var stopCmd = &cobra.Command{
	Use:   "stop",
	Short: "Interrompe o simulador.jar",
	Long: `Interrompe a instancia registrada do simulador.jar.

Esta estrutura inicial registra o comando no CLI. O encerramento real do processo
sera implementado na historia de controle de ciclo de vida do Simulador.`,
	Run: runStop,
}

func init() {
	rootCmd.AddCommand(stopCmd)
	stopCmd.Flags().IntVar(&stopPort, "port", 8081, "Porta do Simulador do HubSaude")
}

func runStop(cmd *cobra.Command, args []string) {
	fmt.Fprintf(cmd.OutOrStdout(), "Comando simulador stop registrado para a porta %d\n", stopPort)
}
