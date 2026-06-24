package cmd

import (
	"fmt"

	"github.com/heitor-barbosa/runner/projetos/assinador/internal/logging"
	"github.com/heitor-barbosa/runner/projetos/assinador/internal/runner"
	"github.com/spf13/cobra"
)

var (
	startPort    int
	startTimeout int
)

var startCmd = &cobra.Command{
	Use:   "start",
	Short: "Inicia o assinador.jar em modo servidor HTTP",
	Long: `Inicia o assinador.jar em background no modo servidor HTTP.

Se ja existir uma instancia ativa na porta informada, o CLI reutiliza essa instancia
e nao inicia outro processo.`,
	RunE: runStart,
}

func init() {
	rootCmd.AddCommand(startCmd)
	startCmd.Flags().IntVar(&startPort, "port", 8080, "Porta do servidor HTTP do assinador.jar")
	startCmd.Flags().IntVar(&startTimeout, "timeout", 0, "Minutos de inatividade antes do encerramento automatico; 0 desativa")
}

func runStart(cmd *cobra.Command, args []string) error {
	logging.Debugf("iniciando assinador start na porta %d timeout=%d", startPort, startTimeout)
	state, err := runner.StartServer(startPort, startTimeout)
	if err != nil {
		return err
	}

	if state.Reused {
		if state.PID > 0 {
			fmt.Fprintf(cmd.OutOrStdout(), "Assinador HTTP ja esta em execucao na porta %d (PID %d)\n", state.Port, state.PID)
			return nil
		}
		fmt.Fprintf(cmd.OutOrStdout(), "Assinador HTTP ja esta em execucao na porta %d\n", state.Port)
		return nil
	}

	fmt.Fprintf(cmd.OutOrStdout(), "Assinador HTTP iniciado na porta %d (PID %d)\n", state.Port, state.PID)
	if state.TimeoutMinutes > 0 {
		fmt.Fprintf(cmd.OutOrStdout(), "Timeout por inatividade: %d minuto(s)\n", state.TimeoutMinutes)
	}
	return nil
}
