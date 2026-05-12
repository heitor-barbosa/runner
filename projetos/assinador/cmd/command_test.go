package cmd

import "testing"

func TestSignCommandRequiresMandatoryFlags(t *testing.T) {
	if err := signCmd.ValidateRequiredFlags(); err == nil {
		t.Fatal("expected sign command to reject missing required flags")
	}
}

func TestValidateCommandRequiresMandatoryFlags(t *testing.T) {
	if err := validateCmd.ValidateRequiredFlags(); err == nil {
		t.Fatal("expected validate command to reject missing required flags")
	}
}

func TestSignAndValidateCommandsAreRegistered(t *testing.T) {
	if _, _, err := rootCmd.Find([]string{"sign"}); err != nil {
		t.Fatalf("sign command not registered: %v", err)
	}
	if _, _, err := rootCmd.Find([]string{"validate"}); err != nil {
		t.Fatalf("validate command not registered: %v", err)
	}
}
