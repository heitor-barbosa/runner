package jdk

import (
	"archive/tar"
	"archive/zip"
	"compress/gzip"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
)

const (
	jdkMajorVersion     = "21"
	hubSaudeDirName     = ".hubsaude"
	localJDKBaseDirName = "jdk"
	downloadBaseURL     = "https://github.com/adoptium/temurin21-binaries/releases/download/jdk-21.0.4%2B7"
)

// FindJava localiza o executável do Java no sistema.
func FindJava() (string, error) {
	// 1. Tenta encontrar java no PATH
	javaPath, err := exec.LookPath("java")
	if err == nil {
		return javaPath, nil
	}

	// 2. Tenta JAVA_HOME
	javaHome := os.Getenv("JAVA_HOME")
	if javaHome != "" {
		javaPath = filepath.Join(javaHome, "bin", javaExecutableName())
		if _, err := os.Stat(javaPath); err == nil {
			return javaPath, nil
		}
	}

	// 3. Tenta JDK provisionado em ~/.hubsaude/jdk
	if userHome, err := os.UserHomeDir(); err == nil {
		if javaPath, err := findJavaInLocalJDK(userHome); err == nil {
			return javaPath, nil
		}
	}

	return "", errors.New("java não encontrado no PATH, em JAVA_HOME ou em ~/.hubsaude/jdk")
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

// InstallJDK instala automaticamente um JDK local em ~/.hubsaude/jdk.
func InstallJDK() error {
	userHome, err := os.UserHomeDir()
	if err != nil {
		return fmt.Errorf("não foi possível determinar o diretório home do usuário: %w", err)
	}

	installRoot := filepath.Join(userHome, hubSaudeDirName, localJDKBaseDirName)
	if err := os.MkdirAll(installRoot, 0o755); err != nil {
		return fmt.Errorf("falha ao criar diretório de instalação do JDK: %w", err)
	}

	archiveURL, archiveType, archiveName, err := jdkDownloadInfo()
	if err != nil {
		return err
	}

	tempDir, err := os.MkdirTemp("", "hubsaude-jdk-download-*")
	if err != nil {
		return fmt.Errorf("não foi possível criar diretório temporário de download: %w", err)
	}
	defer os.RemoveAll(tempDir)

	archivePath := filepath.Join(tempDir, archiveName)
	if err := downloadFile(archiveURL, archivePath); err != nil {
		return err
	}

	if err := extractArchive(archivePath, installRoot, archiveType); err != nil {
		return err
	}

	return nil
}

func javaExecutableName() string {
	if runtime.GOOS == "windows" {
		return "java.exe"
	}
	return "java"
}

func findJavaInLocalJDK(userHome string) (string, error) {
	root := filepath.Join(userHome, hubSaudeDirName, localJDKBaseDirName)
	entries, err := os.ReadDir(root)
	if err != nil {
		return "", err
	}

	for _, entry := range entries {
		if !entry.IsDir() {
			continue
		}

		candidate := filepath.Join(root, entry.Name(), "bin", javaExecutableName())
		if _, err := os.Stat(candidate); err == nil {
			return candidate, nil
		}
	}

	return "", errors.New("nenhum JDK válido encontrado em ~/.hubsaude/jdk")
}

func jdkDownloadInfo() (downloadURL, archiveType, archiveName string, err error) {
	switch runtime.GOOS {
	case "windows":
		if runtime.GOARCH != "amd64" {
			return "", "", "", fmt.Errorf("plataforma não suportada: windows/%s", runtime.GOARCH)
		}
		archiveName = "OpenJDK21U-jdk_x64_windows_hotspot_21.0.4_7.zip"
		archiveType = "zip"
	case "linux":
		switch runtime.GOARCH {
		case "amd64":
			archiveName = "OpenJDK21U-jdk_x64_linux_hotspot_21.0.4_7.tar.gz"
		case "arm64":
			archiveName = "OpenJDK21U-jdk_aarch64_linux_hotspot_21.0.4_7.tar.gz"
		default:
			return "", "", "", fmt.Errorf("plataforma não suportada: linux/%s", runtime.GOARCH)
		}
		archiveType = "tar.gz"
	case "darwin":
		switch runtime.GOARCH {
		case "amd64":
			archiveName = "OpenJDK21U-jdk_x64_mac_hotspot_21.0.4_7.tar.gz"
		case "arm64":
			archiveName = "OpenJDK21U-jdk_aarch64_mac_hotspot_21.0.4_7.tar.gz"
		default:
			return "", "", "", fmt.Errorf("plataforma não suportada: darwin/%s", runtime.GOARCH)
		}
		archiveType = "tar.gz"
	default:
		return "", "", "", fmt.Errorf("plataforma não suportada: %s/%s", runtime.GOOS, runtime.GOARCH)
	}

	downloadURL = fmt.Sprintf("%s/%s", downloadBaseURL, archiveName)
	return downloadURL, archiveType, archiveName, nil
}

func downloadFile(url, dest string) error {
	resp, err := http.Get(url)
	if err != nil {
		return fmt.Errorf("falha ao baixar JDK: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("falha ao baixar JDK: status %d", resp.StatusCode)
	}

	out, err := os.Create(dest)
	if err != nil {
		return fmt.Errorf("falha ao criar arquivo de download: %w", err)
	}
	defer out.Close()

	if _, err := io.Copy(out, resp.Body); err != nil {
		return fmt.Errorf("falha ao salvar arquivo de download: %w", err)
	}

	return nil
}

func extractArchive(archivePath, dest, archiveType string) error {
	switch archiveType {
	case "zip":
		return extractZipFile(archivePath, dest)
	case "tar.gz":
		return extractTarGzFile(archivePath, dest)
	default:
		return fmt.Errorf("tipo de arquivo desconhecido: %s", archiveType)
	}
}

func extractZipFile(archivePath, dest string) error {
	reader, err := zip.OpenReader(archivePath)
	if err != nil {
		return fmt.Errorf("falha ao abrir arquivo ZIP: %w", err)
	}
	defer reader.Close()

	for _, file := range reader.File {
		path, err := sanitizeExtractPath(dest, file.Name)
		if err != nil {
			return err
		}

		if file.FileInfo().IsDir() {
			if err := os.MkdirAll(path, file.Mode()); err != nil {
				return err
			}
			continue
		}

		if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
			return err
		}

		dstFile, err := os.OpenFile(path, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, file.Mode())
		if err != nil {
			return err
		}

		srcFile, err := file.Open()
		if err != nil {
			dstFile.Close()
			return err
		}

		if _, err := io.Copy(dstFile, srcFile); err != nil {
			srcFile.Close()
			dstFile.Close()
			return err
		}

		srcFile.Close()
		dstFile.Close()
	}

	return nil
}

func extractTarGzFile(archivePath, dest string) error {
	file, err := os.Open(archivePath)
	if err != nil {
		return fmt.Errorf("falha ao abrir arquivo TAR.GZ: %w", err)
	}
	defer file.Close()

	gzReader, err := gzip.NewReader(file)
	if err != nil {
		return fmt.Errorf("falha ao ler gzip: %w", err)
	}
	defer gzReader.Close()

	tarReader := tar.NewReader(gzReader)
	for {
		header, err := tarReader.Next()
		if err == io.EOF {
			break
		}
		if err != nil {
			return fmt.Errorf("falha ao ler TAR: %w", err)
		}

		path, err := sanitizeExtractPath(dest, header.Name)
		if err != nil {
			return err
		}

		switch header.Typeflag {
		case tar.TypeDir:
			if err := os.MkdirAll(path, os.FileMode(header.Mode)); err != nil {
				return err
			}
		case tar.TypeReg:
			if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
				return err
			}
			fileOut, err := os.OpenFile(path, os.O_CREATE|os.O_RDWR|os.O_TRUNC, os.FileMode(header.Mode))
			if err != nil {
				return err
			}
			if _, err := io.Copy(fileOut, tarReader); err != nil {
				fileOut.Close()
				return err
			}
			fileOut.Close()
		case tar.TypeSymlink:
			if runtime.GOOS == "windows" {
				continue
			}
			if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
				return err
			}
			if err := os.Symlink(header.Linkname, path); err != nil {
				return err
			}
		default:
			// Ignora outros tipos de entrada
		}
	}

	return nil
}

func sanitizeExtractPath(dest, filePath string) (string, error) {
	cleaned := filepath.Clean(filepath.Join(dest, filePath))
	destClean := filepath.Clean(dest)
	if !strings.HasPrefix(cleaned, destClean+string(os.PathSeparator)) && cleaned != destClean {
		return "", fmt.Errorf("caminho de extração inválido: %s", filePath)
	}
	return cleaned, nil
}
