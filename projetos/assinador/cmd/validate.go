package cmd

import (
	"fmt"
	"os"

	"example.com/assinador/internal/runner"
	"github.com/spf13/cobra"
)

var validateCmd = &cobra.Command{
	Use:   "validate",
	Short: "Valida uma assinatura digital simulada",
	Long: `Valida uma assinatura digital simulada invocando o assinador.jar localmente via java -jar.

O valor de --signature-data corresponde ao campo Signature.data de um recurso FHIR Signature,
produzido pelo comando 'sign'.

Exemplos:
  assinatura validate \
    --signature-data "<base64 de Signature.data>" \
    --timestamp 1751328001 \
    --policy "https://fhir.saude.go.gov.br/r4/seguranca/ImplementationGuide/br.go.ses.seguranca|0.1.2"`,
	RunE: runValidate,
}

// flags do comando validate
var (
	validateSignatureData string
	validateTimestamp     int64
	validatePolicy        string
)

func init() {
	rootCmd.AddCommand(validateCmd)

	validateCmd.Flags().StringVar(&validateSignatureData, "signature-data", "", "Valor de Signature.data em base64 (obrigatório)")
	validateCmd.Flags().Int64Var(&validateTimestamp, "timestamp", 0, "Timestamp Unix UTC de referência em segundos (obrigatório)")
	validateCmd.Flags().StringVar(&validatePolicy, "policy", "https://fhir.saude.go.gov.br/r4/seguranca/ImplementationGuide/br.go.ses.seguranca|0.1.2", "URI da política de assinatura")

	_ = validateCmd.MarkFlagRequired("signature-data")
	_ = validateCmd.MarkFlagRequired("timestamp")
}

func runValidate(cmd *cobra.Command, args []string) error {
	payload := map[string]interface{}{
		"signatureData":      validateSignatureData,
		"referenceTimestamp": validateTimestamp,
		"policyUri":          validatePolicy,
	}

	fmt.Fprintln(os.Stderr, "Invocando assinador.jar (modo local)...")

	resp, err := runner.InvokeValidate(payload)
	if err != nil {
		return fmt.Errorf("erro ao invocar assinador.jar: %w", err)
	}

	printValidateResult(resp)
	if !resp.Success {
		os.Exit(1)
	}
	return nil
}

func printValidateResult(resp *runner.Response) {
	if resp.Success {
		fmt.Println()
		fmt.Println("✓ Assinatura válida!")
		fmt.Println("  Resultado:", resp.Data)
	} else {
		fmt.Fprintf(os.Stderr, "\n✗ Assinatura inválida ou erro na validação\n")
		fmt.Fprintf(os.Stderr, "  Código:    %s\n", resp.ErrorCode)
		fmt.Fprintf(os.Stderr, "  Mensagem:  %s\n", resp.ErrorMessage)
	}
}
