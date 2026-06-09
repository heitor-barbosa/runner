package lifecycle

import (
	"net"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
	"time"
)

func TestHelperProcess(t *testing.T) {
	if os.Getenv("GO_WANT_HELPER_PROCESS") != "1" {
		return
	}
	// Keep the helper process alive until killed by the test.
	select {}
}

func fakeCommand(name string, arg ...string) *exec.Cmd {
	cs := []string{"-test.run=TestHelperProcess", "--"}
	cs = append(cs, arg...)
	cmd := exec.Command(os.Args[0], cs...)
	cmd.Env = append(os.Environ(), "GO_WANT_HELPER_PROCESS=1")
	return cmd
}

func TestIsPortAvailableReturnsFalseWhenOccupied(t *testing.T) {
	listener, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatalf("failed to open test listener: %v", err)
	}
	defer listener.Close()

	port := listener.Addr().(*net.TCPAddr).Port
	if isPortAvailable(port) {
		t.Fatalf("expected port %d to be unavailable", port)
	}
}

func TestWriteSimulatorStateCreatesStateFile(t *testing.T) {
	tmpHome := t.TempDir()
	oldHome := userHomeDir
	userHomeDir = func() (string, error) { return tmpHome, nil }
	defer func() { userHomeDir = oldHome }()

	state := &SimulatorState{
		PID:       42,
		Port:      8081,
		JarPath:   "/tmp/simulador.jar",
		StartedAt: time.Now().UTC(),
	}

	if err := writeSimulatorState(state); err != nil {
		t.Fatalf("expected writeSimulatorState to succeed, got %v", err)
	}

	path := simulatorStatePath(state.Port)
	if _, err := os.Stat(path); err != nil {
		t.Fatalf("expected state file to exist at %s, got %v", path, err)
	}

	if !filepath.IsAbs(path) {
		t.Fatalf("expected state file path to be absolute, got %s", path)
	}
}

func TestStartSimulatorWritesStateAndStartsCommand(t *testing.T) {
	tmpHome := t.TempDir()
	oldHome := userHomeDir
	oldCommand := newCommand
	userHomeDir = func() (string, error) { return tmpHome, nil }
	newCommand = fakeCommand
	defer func() {
		userHomeDir = oldHome
		newCommand = oldCommand
	}()

	state, err := StartSimulator("/tmp/simulador.jar", 8082)
	if err != nil {
		t.Fatalf("expected StartSimulator to succeed, got %v", err)
	}
	defer func() {
		proc, err := os.FindProcess(state.PID)
		if err == nil {
			_ = proc.Kill()
		}
	}()

	path := simulatorStatePath(state.Port)
	if _, err := os.Stat(path); err != nil {
		t.Fatalf("expected state file to exist at %s, got %v", path, err)
	}

	content, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("failed to read state file: %v", err)
	}

	if !strings.Contains(string(content), "\"port\": 8082") {
		t.Fatalf("expected state file to contain port 8082, got %s", content)
	}
}

func TestStatusSimulatorReportsActiveProcess(t *testing.T) {
	tmpHome := t.TempDir()
	oldHome := userHomeDir
	oldProcessActive := processActive
	userHomeDir = func() (string, error) { return tmpHome, nil }
	processActive = func(pid int) bool { return pid == 1234 }
	defer func() {
		userHomeDir = oldHome
		processActive = oldProcessActive
	}()

	state := &SimulatorState{
		PID:       1234,
		Port:      8083,
		JarPath:   "/tmp/simulador.jar",
		StartedAt: time.Now().UTC(),
	}
	if err := writeSimulatorState(state); err != nil {
		t.Fatalf("expected writeSimulatorState to succeed, got %v", err)
	}

	status, err := StatusSimulator(8083)
	if err != nil {
		t.Fatalf("expected StatusSimulator to succeed, got %v", err)
	}
	if !status.Active {
		t.Fatalf("expected simulator to be active, got %+v", status)
	}
	if status.State == nil || status.State.PID != state.PID {
		t.Fatalf("expected state PID %d, got %+v", state.PID, status.State)
	}
}

func TestStatusSimulatorReportsMissingStateAsInactive(t *testing.T) {
	tmpHome := t.TempDir()
	oldHome := userHomeDir
	userHomeDir = func() (string, error) { return tmpHome, nil }
	defer func() { userHomeDir = oldHome }()

	status, err := StatusSimulator(8084)
	if err != nil {
		t.Fatalf("expected StatusSimulator to succeed, got %v", err)
	}
	if status.Active {
		t.Fatalf("expected simulator to be inactive, got %+v", status)
	}
	if status.State != nil {
		t.Fatalf("expected no state for missing registration, got %+v", status.State)
	}
}

func TestStopSimulatorKillsProcessAndRemovesState(t *testing.T) {
	tmpHome := t.TempDir()
	oldHome := userHomeDir
	oldCommand := newCommand
	userHomeDir = func() (string, error) { return tmpHome, nil }
	newCommand = fakeCommand
	defer func() {
		userHomeDir = oldHome
		newCommand = oldCommand
	}()

	state, err := StartSimulator("/tmp/simulador.jar", 8085)
	if err != nil {
		t.Fatalf("expected StartSimulator to succeed, got %v", err)
	}
	defer func() {
		proc, err := os.FindProcess(state.PID)
		if err == nil {
			_ = proc.Kill()
		}
	}()

	stopped, err := StopSimulator(8085)
	if err != nil {
		t.Fatalf("expected StopSimulator to succeed, got %v", err)
	}
	if stopped.PID != state.PID {
		t.Fatalf("expected stopped PID %d, got %d", state.PID, stopped.PID)
	}

	if _, err := os.Stat(simulatorStatePath(8085)); !os.IsNotExist(err) {
		t.Fatalf("expected state file to be removed, got %v", err)
	}
}
