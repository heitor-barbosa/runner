package cmd

import (
	"fmt"
	"os"

	"github.com/heitor-barbosa/runner/projetos/assinador/internal/runner"
	"github.com/spf13/cobra"
)

var signCmd = &cobra.Command{
	Use:   "sign",
	Short: "Cria uma assinatura digital simulada",
	Long: `Cria uma assinatura digital simulada invocando o assinador.jar via HTTP quando houver servidor ativo.

Se o servidor nao estiver disponivel, o CLI faz fallback para java -jar. Use --local para forcar o modo local.
O CLI valida as flags obrigatorias e o assinador.jar valida o conteudo recebido.
Em caso de sucesso, exibe o valor de Signature.data (base64) pronto para uso em FHIR.

Exemplos:
  assinatura sign \
    --bundle '{"resourceType":"Bundle","entry":[...]}' \
    --provenance '{"resourceType":"Provenance","target":[...]}' \
    --credential-type PEM \
    --credential-content "$(cat chave.pem)" \
    --certificate-chain '["CERT1_B64","CERT2_B64"]' \
    --timestamp 1751328001 \
    --strategy iat \
    --policy "https://fhir.saude.go.gov.br/r4/seguranca/ImplementationGuide/br.go.ses.seguranca|0.1.2"`,
	RunE: runSign,
}

// flags do comando sign
var (
	signBundle       string
	signProvenance   string
	signCredType     string
	signCredContent  string
	signCredPassword string
	signCredAlias    string
	signPKCS11Config string
	signTokenLabel   string
	signCertChain    string
	signTimestamp    int64
	signStrategy     string
	signPolicy       string
	signLocal        bool
	signPort         int
)

func init() {
	rootCmd.AddCommand(signCmd)

	signCmd.Flags().StringVar(&signBundle, "bundle", "", "Instância Bundle FHIR R4 em JSON (obrigatório)")
	signCmd.Flags().StringVar(&signProvenance, "provenance", "", "Instância Provenance FHIR R4 em JSON (obrigatório)")
	signCmd.Flags().StringVar(&signCredType, "credential-type", "PEM", "Tipo de credencial: PEM, PKCS12, SMARTCARD, TOKEN")
	signCmd.Flags().StringVar(&signCredContent, "credential-content", "", "Conteúdo da credencial (chave PEM ou base64 PKCS12) (obrigatório)")
	signCmd.Flags().StringVar(&signCredPassword, "credential-password", "", "Senha para PEM criptografado ou PKCS12")
	signCmd.Flags().StringVar(&signCredAlias, "credential-alias", "", "Alias da chave no PKCS12 (obrigatório para PKCS12)")
	signCmd.Flags().StringVar(&signCertChain, "certificate-chain", "", "Cadeia de certificados X.509 DER em base64, formato JSON array (obrigatório)")
	signCmd.Flags().Int64Var(&signTimestamp, "timestamp", 0, "Timestamp Unix UTC de referência em segundos (obrigatório)")
	signCmd.Flags().StringVar(&signStrategy, "strategy", "iat", "Estratégia de timestamp: iat ou tsa")
	signCmd.Flags().StringVar(&signPolicy, "policy", "https://fhir.saude.go.gov.br/r4/seguranca/ImplementationGuide/br.go.ses.seguranca|0.1.2", "URI da política de assinatura")
	signCmd.Flags().BoolVar(&signLocal, "local", false, "Força a invocação local via java -jar, ignorando servidor HTTP ativo")
	signCmd.Flags().IntVar(&signPort, "port", 8080, "Porta do servidor HTTP do assinador.jar")

	signCmd.Flags().StringVar(&signPKCS11Config, "pkcs11-config", "", "Caminho do arquivo de configuracao SunPKCS11 para TOKEN/SMARTCARD")
	signCmd.Flags().StringVar(&signTokenLabel, "token-label", "", "Rotulo do token criptografico para TOKEN/SMARTCARD")

	_ = signCmd.MarkFlagRequired("bundle")
	_ = signCmd.MarkFlagRequired("provenance")
	_ = signCmd.MarkFlagRequired("credential-content")
	_ = signCmd.MarkFlagRequired("certificate-chain")
	_ = signCmd.MarkFlagRequired("timestamp")
}

func runSign(cmd *cobra.Command, args []string) error {
	// Monta payload para o assinador.jar
	payload := map[string]interface{}{
		"bundle":             signBundle,
		"provenance":         signProvenance,
		"credentialType":     signCredType,
		"credentialContent":  signCredContent,
		"certificateChain":   signCertChain,
		"referenceTimestamp": signTimestamp,
		"strategy":           signStrategy,
		"policyUri":          signPolicy,
	}
	if signCredPassword != "" {
		payload["credentialPassword"] = signCredPassword
	}
	if signCredAlias != "" {
		payload["credentialAlias"] = signCredAlias
	}
	if signPKCS11Config != "" {
		payload["pkcs11ConfigPath"] = signPKCS11Config
	}
	if signTokenLabel != "" {
		payload["tokenLabel"] = signTokenLabel
	}

	resp, err := runner.InvokeSignWithOptions(payload, runner.InvokeOptions{
		Local: signLocal,
		Port:  signPort,
	})
	if err != nil {
		return fmt.Errorf("erro ao invocar assinador.jar: %w", err)
	}

	printSignResult(resp)
	if !resp.Success {
		os.Exit(1)
	}
	return nil
}

func printSignResult(resp *runner.Response) {
	if resp.Success {
		fmt.Println()
		fmt.Println("✓ Assinatura criada com sucesso!")
		fmt.Println()
		fmt.Println("Signature.data (base64):")
		fmt.Println(resp.Data)
	} else {
		fmt.Fprintf(os.Stderr, "\n✗ Falha na criação da assinatura\n")
		fmt.Fprintf(os.Stderr, "  Código:    %s\n", resp.ErrorCode)
		fmt.Fprintf(os.Stderr, "  Mensagem:  %s\n", resp.ErrorMessage)
	}
}
