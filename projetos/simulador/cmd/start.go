package cmd

import (
	"fmt"

	"github.com/heitor-barbosa/runner/projetos/simulador/internal/artifact"
	"github.com/heitor-barbosa/runner/projetos/simulador/internal/lifecycle"
	"github.com/heitor-barbosa/runner/projetos/simulador/internal/logging"
	"github.com/spf13/cobra"
)

var (
	startPort          int
	startSource        string
	resolveJarFunc     = artifact.ResolveJarWithFallback
	downloadJarFunc    = artifact.DownloadJar
	startSimulatorFunc = lifecycle.StartSimulator
)

var startCmd = &cobra.Command{
	Use:   "start",
	Short: "Inicia o simulador.jar",
	Long: `Prepara a inicializacao do simulador.jar.

O comando localiza o simulador.jar nos caminhos esperados ou baixa o artefato
para ~/.hubsaude/ quando uma URL for informada por --source.
Quando o artefato nao existe localmente, o CLI tenta buscar a versao mais
recente no GitHub Releases.`,
	RunE: runStart,
}

func init() {
	rootCmd.AddCommand(startCmd)
	startCmd.Flags().IntVar(&startPort, "port", 8081, "Porta do Simulador do HubSaude")
	startCmd.Flags().StringVar(&startSource, "source", "", "URL para baixar o simulador.jar quando ele nao existir localmente")
}

func runStart(cmd *cobra.Command, args []string) error {
	logging.Debugf("iniciando simulador start na porta %d source=%q", startPort, startSource)
	jar, err := artifact.ResolveLocalJar()
	if err != nil {
		if startSource != "" {
			jar, err = downloadJarFunc(startSource)
			if err != nil {
				return err
			}
			fmt.Fprintf(cmd.OutOrStdout(), "simulador.jar baixado para %s\n", jar.Path)
		} else {
			jar, err = resolveJarFunc()
			if err != nil {
				return err
			}
			fmt.Fprintf(cmd.OutOrStdout(), "simulador.jar baixado e verificado para %s\n", jar.Path)
		}
	} else {
		fmt.Fprintf(cmd.OutOrStdout(), "simulador.jar encontrado em %s\n", jar.Path)
	}

	state, err := startSimulatorFunc(jar.Path, startPort)
	if err != nil {
		return err
	}
	fmt.Fprintf(cmd.OutOrStdout(), "simulador.jar iniciado em PID %d na porta %d\n", state.PID, state.Port)
	fmt.Fprintf(cmd.OutOrStdout(), "Comando simulador start preparado para a porta %d\n", startPort)
	return nil
}
