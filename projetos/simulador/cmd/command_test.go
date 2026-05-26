package cmd

import (
	"bytes"
	"testing"
)

func TestLifecycleCommandsAreRegistered(t *testing.T) {
	for _, command := range []string{"start", "stop", "status"} {
		if _, _, err := rootCmd.Find([]string{command}); err != nil {
			t.Fatalf("%s command not registered: %v", command, err)
		}
	}
}

func TestVersionCommandIsRegistered(t *testing.T) {
	if _, _, err := rootCmd.Find([]string{"version"}); err != nil {
		t.Fatalf("version command not registered: %v", err)
	}
}

func TestVersionCommandOutput(t *testing.T) {
	var output bytes.Buffer
	versionCmd.SetOut(&output)

	versionCmd.Run(versionCmd, nil)

	if got, want := output.String(), "simulador v0.1.0\n"; got != want {
		t.Fatalf("version output = %q, want %q", got, want)
	}
}
