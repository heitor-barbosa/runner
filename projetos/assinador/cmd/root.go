package cmd

import (
	"os"

	"github.com/spf13/cobra"
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "assinatura",
	Short: "CLI para o Sistema Runner de assinatura digital",
	Long: `O Sistema Runner facilita o acesso à funcionalidade de execução de aplicações Java
via linha de comandos, focando na assinatura e validação de assinaturas digitais.`,

	// Run: func(cmd *cobra.Command, args []string) { },
}

func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	// rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.assinatura.yaml)")

	rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

