package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
)

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Exibe a versão atual do CLI",
	Long:  `Exibe a versão corrente da aplicação assinatura.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("assinatura versão v0.1.0")
	},
}

func init() {
	rootCmd.AddCommand(versionCmd)
}
