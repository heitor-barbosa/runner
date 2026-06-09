package lifecycle

import (
	"encoding/json"
	"fmt"
	"net"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strconv"
	"strings"
	"time"
)

const defaultPort = 8081

var userHomeDir = os.UserHomeDir
var newCommand = exec.Command
var findProcess = os.FindProcess
var processActive = defaultProcessActive

// SimulatorState representa o estado do simulador iniciado pelo CLI.
type SimulatorState struct {
	PID       int       `json:"pid"`
	Port      int       `json:"port"`
	JarPath   string    `json:"jarPath"`
	StartedAt time.Time `json:"startedAt"`
}

// SimulatorStatus representa a visibilidade atual do processo registrado.
type SimulatorStatus struct {
	State   *SimulatorState
	Active  bool
	Message string
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

// StatusSimulator consulta o registro local e confirma se o processo ainda existe.
func StatusSimulator(port int) (*SimulatorStatus, error) {
	state, err := readSimulatorState(port)
	if err != nil {
		if os.IsNotExist(err) {
			return &SimulatorStatus{
				Active:  false,
				Message: fmt.Sprintf("Simulador nao esta em execucao na porta %d", normalizePort(port)),
			}, nil
		}
		return nil, err
	}

	active := processActive(state.PID)
	status := &SimulatorStatus{State: state, Active: active}
	if active {
		status.Message = fmt.Sprintf("Simulador em execucao na porta %d (PID %d)", state.Port, state.PID)
		return status, nil
	}

	status.Message = fmt.Sprintf("Simulador registrado na porta %d nao esta ativo (PID %d)", state.Port, state.PID)
	return status, nil
}

// StopSimulator encerra o processo registrado e remove o estado local.
func StopSimulator(port int) (*SimulatorState, error) {
	state, err := readSimulatorState(port)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, fmt.Errorf("simulador nao esta registrado na porta %d", normalizePort(port))
		}
		return nil, err
	}

	process, err := findProcess(state.PID)
	if err != nil {
		_ = removeSimulatorState(state.Port)
		return nil, fmt.Errorf("falha ao localizar processo do simulador (PID %d): %w", state.PID, err)
	}

	if err := process.Kill(); err != nil {
		_ = removeSimulatorState(state.Port)
		return nil, fmt.Errorf("falha ao encerrar simulador (PID %d): %w", state.PID, err)
	}

	if err := removeSimulatorState(state.Port); err != nil && !os.IsNotExist(err) {
		return nil, err
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

func readSimulatorState(port int) (*SimulatorState, error) {
	data, err := os.ReadFile(simulatorStatePath(port))
	if err != nil {
		return nil, err
	}

	var state SimulatorState
	if err := json.Unmarshal(data, &state); err != nil {
		return nil, fmt.Errorf("falha ao ler estado do simulador: %w", err)
	}
	return &state, nil
}

func removeSimulatorState(port int) error {
	if err := os.Remove(simulatorStatePath(port)); err != nil {
		return fmt.Errorf("falha ao remover estado do simulador: %w", err)
	}
	return nil
}

func defaultProcessActive(pid int) bool {
	if pid <= 0 {
		return false
	}

	if runtime.GOOS == "windows" {
		output, err := exec.Command("tasklist", "/FI", fmt.Sprintf("PID eq %d", pid), "/FO", "CSV", "/NH").Output()
		if err != nil {
			return false
		}
		return strings.Contains(string(output), strconv.Itoa(pid))
	}

	return exec.Command("kill", "-0", strconv.Itoa(pid)).Run() == nil
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
