package cmd

import (
	"bytes"
	"errors"
	"testing"

	"github.com/heitor-barbosa/runner/projetos/simulador/internal/lifecycle"
	"github.com/spf13/cobra"
)

func TestRunStatusShowsActiveSimulatorDetails(t *testing.T) {
	oldStatus := statusSimulatorFunc
	oldPort := statusPort
	defer func() {
		statusSimulatorFunc = oldStatus
		statusPort = oldPort
	}()

	statusPort = 9090
	statusSimulatorFunc = func(port int) (*lifecycle.SimulatorStatus, error) {
		if port != 9090 {
			return nil, errors.New("unexpected port")
		}
		return &lifecycle.SimulatorStatus{
			Active:  true,
			Message: "Simulador em execucao na porta 9090 (PID 1234)",
			State:   &lifecycle.SimulatorState{PID: 1234, Port: 9090, JarPath: "/tmp/simulador.jar"},
		}, nil
	}

	var output bytes.Buffer
	cmd := &cobra.Command{}
	cmd.SetOut(&output)

	if err := runStatus(cmd, nil); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	want := "Simulador em execucao na porta 9090 (PID 1234)\nPID: 1234\nPorta: 9090\nJAR: /tmp/simulador.jar\n"
	if got := output.String(); got != want {
		t.Fatalf("unexpected output: %q", got)
	}
}

func TestRunStopStopsRegisteredSimulator(t *testing.T) {
	oldStop := stopSimulatorFunc
	oldPort := stopPort
	defer func() {
		stopSimulatorFunc = oldStop
		stopPort = oldPort
	}()

	stopPort = 9091
	stopSimulatorFunc = func(port int) (*lifecycle.SimulatorState, error) {
		if port != 9091 {
			return nil, errors.New("unexpected port")
		}
		return &lifecycle.SimulatorState{PID: 4321, Port: 9091}, nil
	}

	var output bytes.Buffer
	cmd := &cobra.Command{}
	cmd.SetOut(&output)

	if err := runStop(cmd, nil); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if got, want := output.String(), "Simulador encerrado na porta 9091 (PID 4321)\n"; got != want {
		t.Fatalf("unexpected output: %q", got)
	}
}
