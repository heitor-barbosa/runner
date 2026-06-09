package cmd

import (
	"bytes"
	"errors"
	"testing"

	"github.com/heitor-barbosa/runner/projetos/simulador/internal/artifact"
	"github.com/spf13/cobra"
)

func TestRunStartUsesSourceURLWhenLocalJarMissing(t *testing.T) {
	oldResolve := resolveJarFunc
	oldDownload := downloadJarFunc
	defer func() {
		resolveJarFunc = oldResolve
		downloadJarFunc = oldDownload
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
	if got != "simulador.jar baixado para /tmp/.hubsaude/simulador.jar\nComando simulador start preparado para a porta 8081\n" {
		t.Fatalf("unexpected output: %q", got)
	}
}
