package lifecycle

import (
	"fmt"
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

	args := os.Args[1:]
	for i, arg := range args {
		if arg == "--" {
			args = args[i+1:]
			break
		}
	}

	var port string
	var delay time.Duration
	for i := 0; i < len(args); i++ {
		switch args[i] {
		case "--port":
			if i+1 < len(args) {
				port = args[i+1]
			}
		case "--listen-after":
			if i+1 < len(args) {
				d, err := time.ParseDuration(args[i+1])
				if err == nil {
					delay = d
				}
			}
		}
	}

	if port != "" {
		if delay > 0 {
			time.Sleep(delay)
		}
		ln, err := net.Listen("tcp", fmt.Sprintf("127.0.0.1:%s", port))
		if err != nil {
			fmt.Fprintf(os.Stderr, "failed to listen: %v\n", err)
			os.Exit(1)
		}
		defer ln.Close()
	}

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

func TestStartSimulatorRejectsOccupiedPortBeforeStartingProcess(t *testing.T) {
	listener, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatalf("failed to open test listener: %v", err)
	}
	defer listener.Close()

	tmpHome := t.TempDir()
	oldHome := userHomeDir
	oldCommand := newCommand
	var commandStarted bool
	userHomeDir = func() (string, error) { return tmpHome, nil }
	newCommand = func(name string, arg ...string) *exec.Cmd {
		commandStarted = true
		return fakeCommand(name, arg...)
	}
	defer func() {
		userHomeDir = oldHome
		newCommand = oldCommand
	}()

	port := listener.Addr().(*net.TCPAddr).Port
	_, err = StartSimulator("/tmp/simulador.jar", port)
	if err == nil {
		t.Fatal("expected StartSimulator to reject occupied port")
	}
	if !strings.Contains(err.Error(), "porta") || !strings.Contains(err.Error(), "indisponivel") {
		t.Fatalf("error = %q, want clear occupied-port message", err)
	}
	if commandStarted {
		t.Fatal("expected StartSimulator to fail before starting process")
	}
	if _, err := os.Stat(simulatorStatePath(port)); !os.IsNotExist(err) {
		t.Fatalf("expected no state file for rejected start, stat error = %v", err)
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

func TestWaitForSimulatorReadyReturnsAfterPortIsOpen(t *testing.T) {
	listener, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatalf("failed to listen on ephemeral port: %v", err)
	}
	port := listener.Addr().(*net.TCPAddr).Port
	if err := listener.Close(); err != nil {
		t.Fatalf("failed to close initial listener: %v", err)
	}

	readyCh := make(chan net.Listener, 1)
	errCh := make(chan error, 1)
	go func() {
		time.Sleep(100 * time.Millisecond)
		ln, err := net.Listen("tcp", fmt.Sprintf("127.0.0.1:%d", port))
		if err != nil {
			errCh <- err
			return
		}
		readyCh <- ln
	}()

	start := time.Now()
	if err := waitForSimulatorReady(port, 2*time.Second); err != nil {
		select {
		case err := <-errCh:
			t.Fatalf("failed to open port listener: %v", err)
		default:
			t.Fatalf("expected waitForSimulatorReady to succeed, got %v", err)
		}
	}
	elapsed := time.Since(start)
	if elapsed < 100*time.Millisecond {
		t.Fatalf("expected waitForSimulatorReady to wait for readiness, elapsed %s", elapsed)
	}

	ln := <-readyCh
	defer ln.Close()
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
	oldFindProcess := findProcess
	oldKillProcess := killProcess
	userHomeDir = func() (string, error) { return tmpHome, nil }
	defer func() {
		userHomeDir = oldHome
		findProcess = oldFindProcess
		killProcess = oldKillProcess
	}()

	state := &SimulatorState{
		PID:       4321,
		Port:      8085,
		JarPath:   "/tmp/simulador.jar",
		StartedAt: time.Now().UTC(),
	}
	if err := writeSimulatorState(state); err != nil {
		t.Fatalf("expected writeSimulatorState to succeed, got %v", err)
	}

	var foundPID int
	var killedPID int
	findProcess = func(pid int) (*os.Process, error) {
		foundPID = pid
		return &os.Process{Pid: pid}, nil
	}
	killProcess = func(process *os.Process) error {
		killedPID = process.Pid
		return nil
	}

	stopped, err := StopSimulator(state.Port)
	if err != nil {
		t.Fatalf("expected StopSimulator to succeed, got %v", err)
	}
	if foundPID != state.PID {
		t.Fatalf("findProcess PID = %d, want %d", foundPID, state.PID)
	}
	if killedPID != state.PID {
		t.Fatalf("killProcess PID = %d, want %d", killedPID, state.PID)
	}
	if stopped.PID != state.PID {
		t.Fatalf("expected stopped PID %d, got %d", state.PID, stopped.PID)
	}

	if _, err := os.Stat(simulatorStatePath(state.Port)); !os.IsNotExist(err) {
		t.Fatalf("expected state file to be removed, got %v", err)
	}
}
