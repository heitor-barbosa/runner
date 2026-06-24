package cmd

import (
	"fmt"

	"github.com/heitor-barbosa/runner/projetos/simulador/internal/lifecycle"
	"github.com/heitor-barbosa/runner/projetos/simulador/internal/logging"
	"github.com/spf13/cobra"
)

var (
	statusPort          int
	statusSimulatorFunc = lifecycle.StatusSimulator
)

var statusCmd = &cobra.Command{
	Use:   "status",
	Short: "Consulta o status do simulador.jar",
	Long: `Consulta o status da instancia registrada do simulador.jar.

O comando usa o estado salvo em ~/.hubsaude e verifica se o processo registrado
ainda esta ativo.`,
	RunE: runStatus,
}

func init() {
	rootCmd.AddCommand(statusCmd)
	statusCmd.Flags().IntVar(&statusPort, "port", 8081, "Porta do Simulador do HubSaude")
}

func runStatus(cmd *cobra.Command, args []string) error {
	logging.Debugf("iniciando simulador status na porta %d", statusPort)
	status, err := statusSimulatorFunc(statusPort)
	if err != nil {
		return err
	}

	fmt.Fprintln(cmd.OutOrStdout(), status.Message)
	if status.State != nil {
		fmt.Fprintf(cmd.OutOrStdout(), "PID: %d\n", status.State.PID)
		fmt.Fprintf(cmd.OutOrStdout(), "Porta: %d\n", status.State.Port)
		fmt.Fprintf(cmd.OutOrStdout(), "JAR: %s\n", status.State.JarPath)
	}
	return nil
}
