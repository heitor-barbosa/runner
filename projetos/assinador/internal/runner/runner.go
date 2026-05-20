package runner

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"time"

	"github.com/heitor-barbosa/runner/projetos/assinador/internal/jdk"
)

const defaultServerPort = 8080

// Response representa a resposta JSON do assinador.jar.
type Response struct {
	Success      bool   `json:"success"`
	Data         string `json:"data"`
	ErrorCode    string `json:"errorCode"`
	ErrorMessage string `json:"errorMessage"`
}

// ServerState representa uma instancia do assinador.jar iniciada pelo CLI.
type ServerState struct {
	PID       int       `json:"pid"`
	Port      int       `json:"port"`
	JavaPath  string    `json:"javaPath"`
	JarPath   string    `json:"jarPath"`
	StartedAt time.Time `json:"startedAt"`
	Reused    bool      `json:"-"`
}

// InvokeOptions configura a estrategia de invocacao do assinador.jar.
type InvokeOptions struct {
	Local bool
	Port  int
}

var invokeLocalFunc = invokeLocal

// InvokeSign invoca o assinador.jar para criacao de assinatura.
func InvokeSign(payload map[string]interface{}) (*Response, error) {
	return InvokeSignWithOptions(payload, InvokeOptions{})
}

// InvokeSignWithOptions invoca o assinador.jar para criacao de assinatura com opcoes explicitas.
func InvokeSignWithOptions(payload map[string]interface{}, options InvokeOptions) (*Response, error) {
	return invokeWithFallback("sign", payload, options)
}

// InvokeValidate invoca o assinador.jar para validacao de assinatura.
func InvokeValidate(payload map[string]interface{}) (*Response, error) {
	return InvokeValidateWithOptions(payload, InvokeOptions{})
}

// InvokeValidateWithOptions invoca o assinador.jar para validacao de assinatura com opcoes explicitas.
func InvokeValidateWithOptions(payload map[string]interface{}, options InvokeOptions) (*Response, error) {
	return invokeWithFallback("validate", payload, options)
}

// StartServer inicia o assinador.jar em modo servidor, ou reutiliza a instancia ativa na porta.
func StartServer(port int) (*ServerState, error) {
	port = normalizePort(port)

	if state, err := readServerState(port); err == nil && isServerActive(port) {
		state.Reused = true
		return state, nil
	}

	if isServerActive(port) {
		return &ServerState{Port: port, Reused: true}, nil
	}

	javaPath, err := findJava()
	if err != nil {
		return nil, err
	}

	jarPath, err := findJar()
	if err != nil {
		return nil, err
	}

	cmd := exec.Command(javaPath, "-jar", jarPath, "server", "--port", strconv.Itoa(port)) //nolint:gosec
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Start(); err != nil {
		return nil, fmt.Errorf("falha ao iniciar assinador.jar em modo servidor: %w", err)
	}

	state := &ServerState{
		PID:       cmd.Process.Pid,
		Port:      port,
		JavaPath:  javaPath,
		JarPath:   jarPath,
		StartedAt: time.Now().UTC(),
	}

	if err := writeServerState(state); err != nil {
		_ = cmd.Process.Kill()
		return nil, err
	}

	if err := waitForServer(port, 10*time.Second); err != nil {
		_ = cmd.Process.Kill()
		_ = removeServerState(port)
		return nil, err
	}

	_ = cmd.Process.Release()
	return state, nil
}

// StopServer encerra a instancia registrada do assinador.jar na porta informada.
func StopServer(port int) (*ServerState, error) {
	port = normalizePort(port)

	state, err := readServerState(port)
	if err != nil {
		if isServerActive(port) {
			return nil, fmt.Errorf("assinador.jar esta ativo na porta %d, mas nao foi iniciado por este CLI ou nao possui registro local", port)
		}
		return nil, fmt.Errorf("nenhuma instancia do assinador.jar registrada na porta %d", port)
	}

	process, err := os.FindProcess(state.PID)
	if err != nil {
		_ = removeServerState(port)
		return state, fmt.Errorf("processo %d nao encontrado: %w", state.PID, err)
	}

	if err := process.Kill(); err != nil {
		_ = removeServerState(port)
		return state, fmt.Errorf("falha ao encerrar processo %d: %w", state.PID, err)
	}

	_ = removeServerState(port)
	return state, nil
}

func invokeWithFallback(command string, payload map[string]interface{}, options InvokeOptions) (*Response, error) {
	port := normalizePort(options.Port)
	if !options.Local && hasUsableHTTPServer(port) {
		response, err := invokeHTTP(command, payload, port)
		if err == nil {
			return response, nil
		}
	}

	return invokeLocalFunc(command, payload)
}

func invokeHTTP(command string, payload map[string]interface{}, port int) (*Response, error) {
	jsonBytes, err := json.Marshal(payload)
	if err != nil {
		return nil, fmt.Errorf("erro ao serializar parametros para HTTP: %w", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	url := fmt.Sprintf("http://127.0.0.1:%d/%s", normalizePort(port), command)
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, url, bytes.NewReader(jsonBytes))
	if err != nil {
		return nil, fmt.Errorf("erro ao criar requisicao HTTP: %w", err)
	}
	req.Header.Set("Content-Type", "application/json; charset=utf-8")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("erro ao chamar assinador.jar via HTTP: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("erro ao ler resposta HTTP do assinador.jar: %w", err)
	}

	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusBadRequest {
		return nil, fmt.Errorf("resposta HTTP inesperada do assinador.jar: %s\nSaida recebida: %s", resp.Status, string(body))
	}

	return parseResponse(body)
}

func invokeLocal(command string, payload map[string]interface{}) (*Response, error) {
	javaPath, err := findJava()
	if err != nil {
		return nil, err
	}

	jarPath, err := findJar()
	if err != nil {
		return nil, err
	}

	jsonBytes, err := json.Marshal(payload)
	if err != nil {
		return nil, fmt.Errorf("erro ao serializar parametros: %w", err)
	}

	args := []string{"-jar", jarPath, command, "--json", string(jsonBytes)}
	cmd := exec.Command(javaPath, args...) //nolint:gosec

	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	if err := cmd.Run(); err != nil {
		if stdout.Len() > 0 {
			return parseResponse(stdout.Bytes())
		}
		return nil, fmt.Errorf(
			"falha ao executar assinador.jar: %w\nSaida de erro: %s",
			err,
			stderr.String(),
		)
	}

	return parseResponse(stdout.Bytes())
}

func hasUsableHTTPServer(port int) bool {
	if state, err := readServerState(port); err == nil {
		return isServerActive(normalizePort(state.Port))
	}
	return isServerActive(port)
}

func findJava() (string, error) {
	javaPath, err := jdk.FindJava()
	if err == nil {
		return javaPath, nil
	}

	if installErr := jdk.InstallJDK(); installErr != nil {
		return "", fmt.Errorf(
			"java nao encontrado e instalacao automatica do JDK falhou: %v\nDetalhe: %w",
			installErr,
			err,
		)
	}

	javaPath, err = jdk.FindJava()
	if err != nil {
		return "", fmt.Errorf("java nao encontrado apos instalacao automatica do JDK: %w", err)
	}
	return javaPath, nil
}

func parseResponse(data []byte) (*Response, error) {
	var resp Response
	if err := json.Unmarshal(bytes.TrimSpace(data), &resp); err != nil {
		return nil, fmt.Errorf(
			"resposta do assinador.jar nao e JSON valido: %w\nSaida recebida: %s",
			err,
			string(data),
		)
	}
	return &resp, nil
}

// findJar procura o assinador.jar nas localizacoes esperadas:
// 1. Mesmo diretorio do executavel assinatura
// 2. ~/.hubsaude/
// 3. Diretorio atual
func findJar() (string, error) {
	candidates := jarCandidates()
	for _, path := range candidates {
		if _, err := os.Stat(path); err == nil {
			return path, nil
		}
	}
	return "", fmt.Errorf(
		"assinador.jar nao encontrado. Localizacoes verificadas:\n%s\n\nColoque o assinador.jar em um dos diretorios acima.",
		formatCandidates(candidates),
	)
}

func jarCandidates() []string {
	var candidates []string

	if exe, err := os.Executable(); err == nil {
		candidates = append(candidates, filepath.Join(filepath.Dir(exe), "assinador.jar"))
	}

	if home, err := os.UserHomeDir(); err == nil {
		candidates = append(candidates, filepath.Join(home, ".hubsaude", "assinador.jar"))
	}

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

func normalizePort(port int) int {
	if port <= 0 {
		return defaultServerPort
	}
	return port
}

func waitForServer(port int, timeout time.Duration) error {
	deadline := time.Now().Add(timeout)
	for time.Now().Before(deadline) {
		if isServerActive(port) {
			return nil
		}
		time.Sleep(200 * time.Millisecond)
	}
	return fmt.Errorf("assinador.jar nao respondeu em http://127.0.0.1:%d/health em %s", port, timeout)
}

func isServerActive(port int) bool {
	ctx, cancel := context.WithTimeout(context.Background(), 700*time.Millisecond)
	defer cancel()

	url := fmt.Sprintf("http://127.0.0.1:%d/health", normalizePort(port))
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return false
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return false
	}
	defer resp.Body.Close()
	return resp.StatusCode == http.StatusOK
}

func readServerState(port int) (*ServerState, error) {
	data, err := os.ReadFile(serverStatePath(port))
	if err != nil {
		return nil, err
	}

	var state ServerState
	if err := json.Unmarshal(data, &state); err != nil {
		return nil, err
	}
	return &state, nil
}

func writeServerState(state *ServerState) error {
	dir, err := hubSaudeDir()
	if err != nil {
		return err
	}
	if err := os.MkdirAll(dir, 0o755); err != nil {
		return fmt.Errorf("falha ao criar diretorio de estado: %w", err)
	}

	data, err := json.MarshalIndent(state, "", "  ")
	if err != nil {
		return err
	}
	if err := os.WriteFile(serverStatePath(state.Port), data, 0o600); err != nil {
		return fmt.Errorf("falha ao gravar estado do servidor: %w", err)
	}
	return nil
}

func removeServerState(port int) error {
	err := os.Remove(serverStatePath(port))
	if err != nil && !os.IsNotExist(err) {
		return err
	}
	return nil
}

func serverStatePath(port int) string {
	dir, err := hubSaudeDir()
	if err != nil {
		return filepath.Join(".hubsaude", fmt.Sprintf("assinador-server-%d.json", normalizePort(port)))
	}
	return filepath.Join(dir, fmt.Sprintf("assinador-server-%d.json", normalizePort(port)))
}

func hubSaudeDir() (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", fmt.Errorf("nao foi possivel determinar o diretorio home do usuario: %w", err)
	}
	return filepath.Join(home, ".hubsaude"), nil
}
