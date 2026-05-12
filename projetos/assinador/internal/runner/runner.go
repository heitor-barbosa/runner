package runner

import (
	"bytes"
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"

	"example.com/assinador/internal/jdk"
)

// Response representa a resposta JSON do assinador.jar.
type Response struct {
	Success      bool   `json:"success"`
	Data         string `json:"data"`
	ErrorCode    string `json:"errorCode"`
	ErrorMessage string `json:"errorMessage"`
}

// InvokeSign invoca o assinador.jar no modo local para criação de assinatura.
// Constrói o payload JSON, localiza java e assinador.jar, executa e retorna a resposta.
func InvokeSign(payload map[string]interface{}) (*Response, error) {
	return invoke("sign", payload)
}

// InvokeValidate invoca o assinador.jar no modo local para validação de assinatura.
func InvokeValidate(payload map[string]interface{}) (*Response, error) {
	return invoke("validate", payload)
}

// ── internos ──────────────────────────────────────────────────────────────────

func invoke(command string, payload map[string]interface{}) (*Response, error) {
	// 1. Localiza java
	javaPath, err := jdk.FindJava()
	if err != nil {
		if installErr := jdk.InstallJDK(); installErr != nil {
			return nil, fmt.Errorf(
				"java não encontrado e instalação automática do JDK falhou: %v\nDetalhe: %w",
				installErr,
				err,
			)
		}

		javaPath, err = jdk.FindJava()
		if err != nil {
			return nil, fmt.Errorf(
				"java não encontrado após instalação automática do JDK: %w",
				err,
			)
		}
	}

	// 2. Localiza assinador.jar
	jarPath, err := findJar()
	if err != nil {
		return nil, err
	}

	// 3. Serializa o payload para JSON
	jsonBytes, err := json.Marshal(payload)
	if err != nil {
		return nil, fmt.Errorf("erro ao serializar parâmetros: %w", err)
	}

	// 4. Monta e executa o comando: java -jar assinador.jar <command> --json '<payload>'
	args := []string{"-jar", jarPath, command, "--json", string(jsonBytes)}
	cmd := exec.Command(javaPath, args...) //nolint:gosec

	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	if err := cmd.Run(); err != nil {
		// Verifica se há saída de erro do assinador.jar (pode conter JSON de erro)
		if stdout.Len() > 0 {
			return parseResponse(stdout.Bytes())
		}
		return nil, fmt.Errorf(
			"falha ao executar assinador.jar: %w\nSaída de erro: %s",
			err, stderr.String(),
		)
	}

	return parseResponse(stdout.Bytes())
}

func parseResponse(data []byte) (*Response, error) {
	var resp Response
	if err := json.Unmarshal(bytes.TrimSpace(data), &resp); err != nil {
		return nil, fmt.Errorf(
			"resposta do assinador.jar não é JSON válido: %w\nSaída recebida: %s",
			err, string(data),
		)
	}
	return &resp, nil
}

// findJar procura o assinador.jar nas localizações esperadas:
// 1. Mesmo diretório do executável assinatura
// 2. ~/.hubsaude/
// 3. Diretório atual
func findJar() (string, error) {
	candidates := jarCandidates()
	for _, path := range candidates {
		if _, err := os.Stat(path); err == nil {
			return path, nil
		}
	}
	return "", fmt.Errorf(
		"assinador.jar não encontrado. Localizações verificadas:\n%s\n\nColoque o assinador.jar em um dos diretórios acima.",
		formatCandidates(candidates),
	)
}

func jarCandidates() []string {
	var candidates []string

	// 1. Diretório do executável
	if exe, err := os.Executable(); err == nil {
		candidates = append(candidates, filepath.Join(filepath.Dir(exe), "assinador.jar"))
	}

	// 2. ~/.hubsaude/
	if home, err := os.UserHomeDir(); err == nil {
		candidates = append(candidates, filepath.Join(home, ".hubsaude", "assinador.jar"))
	}

	// 3. Diretório atual
	if wd, err := os.Getwd(); err == nil {
		candidates = append(candidates, filepath.Join(wd, "assinador.jar"))
	}

	return candidates
}

func formatCandidates(paths []string) string {
	var sb bytes.Buffer
	for _, p := range paths {
		sb.WriteString("  - ")
		sb.WriteString(p)
		sb.WriteString("\n")
	}
	return sb.String()
}
