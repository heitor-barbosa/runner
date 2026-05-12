package jdk

import (
	"errors"
	"path/filepath"
	"testing"
)

func TestParseJavaMajorVersion(t *testing.T) {
	tests := map[string]int{
		`openjdk version "21.0.8" 2025-07-15 LTS`: 21,
		`java version "21" 2023-09-19 LTS`:        21,
		`java version "1.8.0_402"`:                8,
	}

	for input, expected := range tests {
		actual, err := parseJavaMajorVersion(input)
		if err != nil {
			t.Fatalf("parseJavaMajorVersion(%q) returned error: %v", input, err)
		}
		if actual != expected {
			t.Fatalf("parseJavaMajorVersion(%q) = %d, want %d", input, actual, expected)
		}
	}
}

func TestParseJavaMajorVersionRejectsUnknownOutput(t *testing.T) {
	if _, err := parseJavaMajorVersion("version unavailable"); err == nil {
		t.Fatal("expected parseJavaMajorVersion to reject unknown output")
	}
}

func TestDedupePaths(t *testing.T) {
	input := []string{
		filepath.Join("tmp", "jdk", "bin", "java"),
		filepath.Join("tmp", "jdk", "bin", "java"),
	}
	paths := dedupePaths(input)
	if len(paths) != 1 {
		t.Fatalf("dedupePaths returned %d paths, want 1", len(paths))
	}
}

func TestSanitizeExtractPath(t *testing.T) {
	dest := filepath.Join("tmp", "jdk")
	if _, err := sanitizeExtractPath(dest, filepath.Join("jdk-21", "bin", "java")); err != nil {
		t.Fatalf("expected safe archive path, got error: %v", err)
	}
	if _, err := sanitizeExtractPath(dest, filepath.Join("..", "escape")); err == nil {
		t.Fatal("expected path traversal to be rejected")
	}
}

func TestJavaCandidatesFromSourcesIncludesAvailableSources(t *testing.T) {
	paths := javaCandidatesFromSources(
		func() (string, error) {
			return filepath.Join("path", "java"), nil
		},
		filepath.Join("home", "jdk"),
		filepath.Join("users", "runner"),
		func(userHome string) (string, error) {
			return filepath.Join(userHome, ".hubsaude", "jdk", "temurin", "bin", javaExecutableName()), nil
		},
	)

	expected := []string{
		filepath.Join("path", "java"),
		filepath.Join("home", "jdk", "bin", javaExecutableName()),
		filepath.Join("users", "runner", ".hubsaude", "jdk", "temurin", "bin", javaExecutableName()),
	}
	if len(paths) != len(expected) {
		t.Fatalf("javaCandidatesFromSources returned %d paths, want %d", len(paths), len(expected))
	}
	for index, path := range expected {
		if paths[index] != filepath.Clean(path) {
			t.Fatalf("javaCandidatesFromSources[%d] = %q, want %q", index, paths[index], filepath.Clean(path))
		}
	}
}

func TestJavaCandidatesFromSourcesSkipsUnavailableSources(t *testing.T) {
	paths := javaCandidatesFromSources(
		func() (string, error) {
			return "", errors.New("java not found")
		},
		"",
		filepath.Join("users", "runner"),
		func(string) (string, error) {
			return "", errors.New("local JDK not found")
		},
	)

	if len(paths) != 0 {
		t.Fatalf("javaCandidatesFromSources returned %d paths, want 0", len(paths))
	}
}
