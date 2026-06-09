package lifecycle

import (
	"encoding/json"
	"fmt"
	"net"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"time"
)

const defaultPort = 8081

var userHomeDir = os.UserHomeDir
var newCommand = exec.Command

// SimulatorState representa o estado do simulador iniciado pelo CLI.
type SimulatorState struct {
	PID       int       `json:"pid"`
	Port      int       `json:"port"`
	JarPath   string    `json:"jarPath"`
	StartedAt time.Time `json:"startedAt"`
}

// StartSimulator inicia o simulador.jar em segundo plano e grava o PID/porta em ~/.hubsaude/.
func StartSimulator(jarPath string, port int) (*SimulatorState, error) {
	port = normalizePort(port)
	if jarPath == "" {
		return nil, fmt.Errorf("caminho do simulador.jar nao informado")
	}

	if !isPortAvailable(port) {
		return nil, fmt.Errorf("porta %d indisponivel: ja existe um processo escutando", port)
	}

	devNull, err := os.OpenFile(os.DevNull, os.O_RDWR, 0)
	if err != nil {
		return nil, fmt.Errorf("falha ao abrir devnull: %w", err)
	}
	defer devNull.Close()

	cmd := newCommand("java", "-jar", jarPath, "--port", strconv.Itoa(port))
	cmd.Stdout = devNull
	cmd.Stderr = devNull

	if err := cmd.Start(); err != nil {
		return nil, fmt.Errorf("falha ao iniciar simulador.jar: %w", err)
	}

	state := &SimulatorState{
		PID:       cmd.Process.Pid,
		Port:      port,
		JarPath:   jarPath,
		StartedAt: time.Now().UTC(),
	}

	if err := writeSimulatorState(state); err != nil {
		_ = cmd.Process.Kill()
		return nil, err
	}

	if err := cmd.Process.Release(); err != nil {
		return nil, fmt.Errorf("falha ao liberar processo do simulador: %w", err)
	}

	return state, nil
}

func normalizePort(port int) int {
	if port <= 0 {
		return defaultPort
	}
	return port
}

func isPortAvailable(port int) bool {
	listener, err := net.Listen("tcp", fmt.Sprintf("127.0.0.1:%d", normalizePort(port)))
	if err != nil {
		return false
	}
	_ = listener.Close()
	return true
}

func writeSimulatorState(state *SimulatorState) error {
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

	if err := os.WriteFile(simulatorStatePath(state.Port), data, 0o600); err != nil {
		return fmt.Errorf("falha ao gravar estado do simulador: %w", err)
	}
	return nil
}

func simulatorStatePath(port int) string {
	dir, err := hubSaudeDir()
	if err != nil {
		return filepath.Join(".hubsaude", fmt.Sprintf("simulador-server-%d.json", normalizePort(port)))
	}
	return filepath.Join(dir, fmt.Sprintf("simulador-server-%d.json", normalizePort(port)))
}

func hubSaudeDir() (string, error) {
	home, err := userHomeDir()
	if err != nil {
		return "", fmt.Errorf("nao foi possivel determinar o diretorio home do usuario: %w", err)
	}
	return filepath.Join(home, ".hubsaude"), nil
}
