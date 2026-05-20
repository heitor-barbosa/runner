package cmd

import (
	"fmt"

	"github.com/heitor-barbosa/runner/projetos/assinador/internal/runner"
	"github.com/spf13/cobra"
)

var stopPort int

var stopCmd = &cobra.Command{
	Use:   "stop",
	Short: "Interrompe o assinador.jar em modo servidor HTTP",
	Long: `Interrompe a instancia do assinador.jar registrada pelo CLI.

O comando usa o PID salvo em ~/.hubsaude para encerrar o processo da porta informada
e remove o arquivo de estado local apos o encerramento.`,
	RunE: runStop,
}

func init() {
	rootCmd.AddCommand(stopCmd)
	stopCmd.Flags().IntVar(&stopPort, "port", 8080, "Porta do servidor HTTP do assinador.jar")
}

func runStop(cmd *cobra.Command, args []string) error {
	state, err := runner.StopServer(stopPort)
	if err != nil {
		return err
	}

	fmt.Fprintf(cmd.OutOrStdout(), "Assinador HTTP encerrado na porta %d (PID %d)\n", state.Port, state.PID)
	return nil
}
