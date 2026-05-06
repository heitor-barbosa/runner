package jdk

import (
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

// FindJava localiza o executável do Java no sistema.
func FindJava() (string, error) {
	// Tenta encontrar java no PATH
	javaPath, err := exec.LookPath("java")
	if err == nil {
		return javaPath, nil
	}

	// Tenta JAVA_HOME
	javaHome := os.Getenv("JAVA_HOME")
	if javaHome != "" {
		binPath := filepath.Join(javaHome, "bin", "java.exe")
		if _, err := os.Stat(binPath); err == nil {
			return binPath, nil
		}
		binPath = filepath.Join(javaHome, "bin", "java")
		if _, err := os.Stat(binPath); err == nil {
			return binPath, nil
		}
	}

	return "", errors.New("java não encontrado no PATH nem em JAVA_HOME")
}

// GetJDKVersion retorna a versão do JDK instalado.
func GetJDKVersion(javaPath string) (string, error) {
	cmd := exec.Command(javaPath, "-version")
	output, err := cmd.CombinedOutput()
	if err != nil {
		return "", fmt.Errorf("erro ao obter versão do Java: %w", err)
	}

	lines := strings.Split(string(output), "\n")
	if len(lines) > 0 {
		return strings.TrimSpace(lines[0]), nil
	}

	return "", errors.New("não foi possível determinar a versão do Java")
}

// InstallJDK instala um JDK (placeholder para implementação futura).
func InstallJDK() error {
	return errors.New("instalação de JDK ainda não implementada")
}
