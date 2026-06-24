package cmd

import (
	"fmt"

	"github.com/heitor-barbosa/runner/projetos/simulador/internal/lifecycle"
	"github.com/heitor-barbosa/runner/projetos/simulador/internal/logging"
	"github.com/spf13/cobra"
)

var (
	stopPort          int
	stopSimulatorFunc = lifecycle.StopSimulator
)

var stopCmd = &cobra.Command{
	Use:   "stop",
	Short: "Interrompe o simulador.jar",
	Long: `Interrompe a instancia registrada do simulador.jar.

O comando usa o PID salvo em ~/.hubsaude para encerrar o processo da porta
informada e remove o arquivo de estado local apos o encerramento.`,
	RunE: runStop,
}

func init() {
	rootCmd.AddCommand(stopCmd)
	stopCmd.Flags().IntVar(&stopPort, "port", 8081, "Porta do Simulador do HubSaude")
}

func runStop(cmd *cobra.Command, args []string) error {
	logging.Debugf("iniciando simulador stop na porta %d", stopPort)
	state, err := stopSimulatorFunc(stopPort)
	if err != nil {
		return err
	}

	fmt.Fprintf(cmd.OutOrStdout(), "Simulador encerrado na porta %d (PID %d)\n", state.Port, state.PID)
	return nil
}
