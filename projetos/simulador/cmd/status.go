package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var statusPort int

var statusCmd = &cobra.Command{
	Use:   "status",
	Short: "Consulta o status do simulador.jar",
	Long: `Consulta o status da instancia registrada do simulador.jar.

Esta estrutura inicial registra o comando no CLI. A consulta real de PID, porta e
processo ativo sera implementada na historia de monitoramento do Simulador.`,
	Run: runStatus,
}

func init() {
	rootCmd.AddCommand(statusCmd)
	statusCmd.Flags().IntVar(&statusPort, "port", 8081, "Porta do Simulador do HubSaude")
}

func runStatus(cmd *cobra.Command, args []string) {
	fmt.Fprintf(cmd.OutOrStdout(), "Comando simulador status registrado para a porta %d\n", statusPort)
}
