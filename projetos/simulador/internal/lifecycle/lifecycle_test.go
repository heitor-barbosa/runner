package lifecycle

import (
	"net"
	"os"
	"path/filepath"
	"testing"
	"time"
)

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
