package artifact

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"
)

const jarName = "simulador.jar"

type JarResult struct {
	Path string
}

func DownloadJar(sourceURL string) (*JarResult, error) {
	sourceURL = strings.TrimSpace(sourceURL)
	if sourceURL == "" {
		return nil, fmt.Errorf("URL do simulador.jar nao informada")
	}

	destination, err := cacheJarPath()
	if err != nil {
		return nil, err
	}
	if err := os.MkdirAll(filepath.Dir(destination), 0o755); err != nil {
		return nil, fmt.Errorf("falha ao criar cache do simulador.jar: %w", err)
	}

	if err := downloadToFile(sourceURL, destination); err != nil {
		return nil, err
	}
	return &JarResult{Path: destination}, nil
}

func ResolveLocalJar() (*JarResult, error) {
	for _, path := range jarCandidates() {
		if _, err := os.Stat(path); err == nil {
			return &JarResult{Path: path}, nil
		}
	}

	return nil, fmt.Errorf("%s nao encontrado", jarName)
}

// ResolveJarWithFallback attempts to resolve JAR locally, then falls back to GitHub Releases
func ResolveJarWithFallback() (*JarResult, error) {
	// Try local resolution first
	if jar, err := ResolveLocalJar(); err == nil {
		return jar, nil
	}

	// Fall back to GitHub Releases
	return DownloadAndVerifyJarFromRelease()
}

func downloadToFile(sourceURL string, destination string) error {
	client := &http.Client{Timeout: 30 * time.Second}
	response, err := client.Get(sourceURL)
	if err != nil {
		return fmt.Errorf("falha ao baixar simulador.jar: %w", err)
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		return fmt.Errorf("falha ao baixar simulador.jar: HTTP %s", response.Status)
	}

	file, err := os.OpenFile(destination, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0o600)
	if err != nil {
		return fmt.Errorf("falha ao criar simulador.jar em cache: %w", err)
	}
	defer file.Close()

	if _, err := io.Copy(file, response.Body); err != nil {
		return fmt.Errorf("falha ao salvar simulador.jar em cache: %w", err)
	}
	return nil
}

func jarCandidates() []string {
	var candidates []string

	if exe, err := os.Executable(); err == nil {
		candidates = append(candidates, filepath.Join(filepath.Dir(exe), jarName))
	}
	if cache, err := cacheJarPath(); err == nil {
		candidates = append(candidates, cache)
	}
	if wd, err := os.Getwd(); err == nil {
		candidates = append(candidates, filepath.Join(wd, jarName))
	}

	return candidates
}

func cacheJarPath() (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", fmt.Errorf("nao foi possivel determinar o diretorio home: %w", err)
	}
	return filepath.Join(home, ".hubsaude", jarName), nil
}
