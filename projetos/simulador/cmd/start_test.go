package cmd

import (
	"bytes"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"testing"

	"github.com/heitor-barbosa/runner/projetos/simulador/internal/artifact"
	"github.com/heitor-barbosa/runner/projetos/simulador/internal/lifecycle"
	"github.com/spf13/cobra"
)

func TestRunStartUsesSourceURLWhenLocalJarMissing(t *testing.T) {
	oldResolve := resolveJarFunc
	oldDownload := downloadJarFunc
	oldStartSimulator := startSimulatorFunc
	defer func() {
		resolveJarFunc = oldResolve
		downloadJarFunc = oldDownload
		startSimulatorFunc = oldStartSimulator
	}()

	downloadJarFunc = func(sourceURL string) (*artifact.JarResult, error) {
		if sourceURL != "https://example.com/simulador.jar" {
			return nil, errors.New("unexpected source URL")
		}
		return &artifact.JarResult{Path: "/tmp/.hubsaude/simulador.jar"}, nil
	}

	resolveJarFunc = func() (*artifact.JarResult, error) {
		return nil, errors.New("should not call resolveJarFunc when source is provided")
	}

	startSimulatorFunc = func(jarPath string, port int) (*lifecycle.SimulatorState, error) {
		if jarPath != "/tmp/.hubsaude/simulador.jar" {
			return nil, errors.New("unexpected jar path")
		}
		if port != 8081 {
			return nil, errors.New("unexpected port")
		}
		return &lifecycle.SimulatorState{PID: 1234, Port: 8081, JarPath: jarPath}, nil
	}

	oldStartSource := startSource
	oldStartPort := startPort
	defer func() {
		startSource = oldStartSource
		startPort = oldStartPort
	}()

	startSource = "https://example.com/simulador.jar"
	startPort = 8081

	var output bytes.Buffer
	cmd := &cobra.Command{}
	cmd.SetOut(&output)

	if err := runStart(cmd, nil); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	got := output.String()
	if got != "simulador.jar baixado para /tmp/.hubsaude/simulador.jar\nsimulador.jar iniciado em PID 1234 na porta 8081\nComando simulador start preparado para a porta 8081\n" {
		t.Fatalf("unexpected output: %q", got)
	}
}

func TestRunStartUsesLocalJarWhenAvailable(t *testing.T) {
	oldStartSource := startSource
	oldStartPort := startPort
	oldStartSimulator := startSimulatorFunc
	defer func() {
		startSource = oldStartSource
		startPort = oldStartPort
		startSimulatorFunc = oldStartSimulator
	}()

	tmpDir := t.TempDir()
	oldWd, err := os.Getwd()
	if err != nil {
		t.Fatalf("failed to get current working directory: %v", err)
	}
	defer func() {
		_ = os.Chdir(oldWd)
	}()
	if err := os.Chdir(tmpDir); err != nil {
		t.Fatalf("failed to change working directory: %v", err)
	}

	localJarPath := filepath.Join(tmpDir, "simulador.jar")
	if err := os.WriteFile(localJarPath, []byte("dummy"), 0o600); err != nil {
		t.Fatalf("failed to create local jar: %v", err)
	}

	startSource = ""
	startPort = 8081
	startSimulatorFunc = func(jarPath string, port int) (*lifecycle.SimulatorState, error) {
		if jarPath != localJarPath {
			return nil, errors.New("unexpected jar path")
		}
		if port != 8081 {
			return nil, errors.New("unexpected port")
		}
		return &lifecycle.SimulatorState{PID: 1234, Port: port, JarPath: jarPath}, nil
	}

	var output bytes.Buffer
	cmd := &cobra.Command{}
	cmd.SetOut(&output)

	if err := runStart(cmd, nil); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	want := fmt.Sprintf("simulador.jar encontrado em %s\nsimulador.jar iniciado em PID 1234 na porta 8081\nComando simulador start preparado para a porta 8081\n", localJarPath)
	if got := output.String(); got != want {
		t.Fatalf("unexpected output: %q", got)
	}
}
