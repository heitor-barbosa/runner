package artifact

import (
	"fmt"
	"os"
	"path/filepath"
)

const jarName = "simulador.jar"

type JarResult struct {
	Path string
}

func ResolveLocalJar() (*JarResult, error) {
	for _, path := range jarCandidates() {
		if _, err := os.Stat(path); err == nil {
			return &JarResult{Path: path}, nil
		}
	}

	return nil, fmt.Errorf("%s nao encontrado", jarName)
}

func jarCandidates() []string {
	var candidates []string

	if exe, err := os.Executable(); err == nil {
		candidates = append(candidates, filepath.Join(filepath.Dir(exe), jarName))
	}
	if home, err := os.UserHomeDir(); err == nil {
		candidates = append(candidates, filepath.Join(home, ".hubsaude", jarName))
	}
	if wd, err := os.Getwd(); err == nil {
		candidates = append(candidates, filepath.Join(wd, jarName))
	}

	return candidates
}
