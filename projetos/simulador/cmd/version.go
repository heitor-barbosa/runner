package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var Version = "v0.1.0"

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Exibe a versao atual do CLI",
	Long:  "Exibe a versao corrente da aplicacao simulador.",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Fprintf(cmd.OutOrStdout(), "simulador %s\n", Version)
	},
}

func init() {
	rootCmd.AddCommand(versionCmd)
}
