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
	"regexp"
	"runtime"
	"strconv"
	"strings"
	"time"
)

const (
	jdkMajorVersion     = 21
	hubSaudeDirName     = ".hubsaude"
	localJDKBaseDirName = "jdk"
	downloadBaseURL     = "https://github.com/adoptium/temurin21-binaries/releases/download/jdk-21.0.4%2B7"
)

var javaVersionPattern = regexp.MustCompile(`(?i)(?:java|openjdk) version "([^"]+)"`)

// FindJava returns a Java executable backed by a JDK 21 installation.
func FindJava() (string, error) {
	var validationErrors []string
	for _, candidate := range javaCandidates() {
		if candidate == "" {
			continue
		}
		if err := validateJDK21(candidate); err == nil {
			return candidate, nil
		} else {
			validationErrors = append(validationErrors, fmt.Sprintf("%s: %v", candidate, err))
		}
	}

	if len(validationErrors) == 0 {
		return "", errors.New("JDK 21 nao encontrado no PATH, em JAVA_HOME ou em ~/.hubsaude/jdk")
	}

	return "", fmt.Errorf(
		"JDK 21 nao encontrado. Candidatos rejeitados:\n  - %s",
		strings.Join(validationErrors, "\n  - "),
	)
}

// GetJDKVersion returns the first line emitted by `java -version`.
func GetJDKVersion(javaPath string) (string, error) {
	cmd := exec.Command(javaPath, "-version")
	output, err := cmd.CombinedOutput()
	if err != nil {
		return "", fmt.Errorf("erro ao obter versao do Java: %w", err)
	}

	for _, line := range strings.Split(string(output), "\n") {
		line = strings.TrimSpace(line)
		if line != "" {
			return line, nil
		}
	}

	return "", errors.New("nao foi possivel determinar a versao do Java")
}

// InstallJDK downloads and installs a local JDK 21 under ~/.hubsaude/jdk.
func InstallJDK() error {
	userHome, err := os.UserHomeDir()
	if err != nil {
		return fmt.Errorf("nao foi possivel determinar o diretorio home do usuario: %w", err)
	}

	if _, err := findJavaInLocalJDK(userHome); err == nil {
		return nil
	}

	installRoot := filepath.Join(userHome, hubSaudeDirName, localJDKBaseDirName)
	if err := os.MkdirAll(installRoot, 0o755); err != nil {
		return fmt.Errorf("falha ao criar diretorio de instalacao do JDK: %w", err)
	}

	archiveURL, archiveType, archiveName, err := jdkDownloadInfo()
	if err != nil {
		return err
	}

	tempDir, err := os.MkdirTemp("", "hubsaude-jdk-download-*")
	if err != nil {
		return fmt.Errorf("nao foi possivel criar diretorio temporario de download: %w", err)
	}
	defer os.RemoveAll(tempDir)

	archivePath := filepath.Join(tempDir, archiveName)
	if err := downloadFile(archiveURL, archivePath); err != nil {
		return err
	}

	if err := extractArchive(archivePath, installRoot, archiveType); err != nil {
		return err
	}

	if _, err := findJavaInLocalJDK(userHome); err != nil {
		return fmt.Errorf("JDK baixado, mas nenhum Java 21 valido foi encontrado: %w", err)
	}

	return nil
}

func javaCandidates() []string {
	userHome, _ := os.UserHomeDir()
	return javaCandidatesFromSources(
		func() (string, error) {
			return exec.LookPath(javaExecutableName())
		},
		os.Getenv("JAVA_HOME"),
		userHome,
		findJavaInLocalJDK,
	)
}

func javaCandidatesFromSources(
	lookPath func() (string, error),
	javaHome string,
	userHome string,
	findLocal func(string) (string, error),
) []string {
	candidates := make([]string, 0, 4)
	if javaPath, err := lookPath(); err == nil && javaPath != "" {
		candidates = append(candidates, javaPath)
	}
	if javaHome != "" {
		candidates = append(candidates, filepath.Join(javaHome, "bin", javaExecutableName()))
	}
	if userHome != "" {
		if javaPath, err := findLocal(userHome); err == nil && javaPath != "" {
			candidates = append(candidates, javaPath)
		}
	}
	return dedupePaths(candidates)
}

func dedupePaths(paths []string) []string {
	seen := make(map[string]struct{}, len(paths))
	result := make([]string, 0, len(paths))
	for _, path := range paths {
		normalized := filepath.Clean(path)
		key := strings.ToLower(normalized)
		if _, ok := seen[key]; ok {
			continue
		}
		seen[key] = struct{}{}
		result = append(result, normalized)
	}
	return result
}

func validateJDK21(javaPath string) error {
	if _, err := os.Stat(javaPath); err != nil {
		return err
	}

	versionLine, err := GetJDKVersion(javaPath)
	if err != nil {
		return err
	}
	majorVersion, err := parseJavaMajorVersion(versionLine)
	if err != nil {
		return err
	}
	if majorVersion != jdkMajorVersion {
		return fmt.Errorf("versao %d encontrada; versao exigida: %d", majorVersion, jdkMajorVersion)
	}

	javacPath := filepath.Join(filepath.Dir(javaPath), javacExecutableName())
	if _, err := os.Stat(javacPath); err != nil {
		return fmt.Errorf("javac nao encontrado ao lado de %s", javaPath)
	}

	return nil
}

func parseJavaMajorVersion(versionLine string) (int, error) {
	matches := javaVersionPattern.FindStringSubmatch(versionLine)
	if len(matches) != 2 {
		return 0, fmt.Errorf("saida de versao do Java nao reconhecida: %s", versionLine)
	}

	versionParts := strings.Split(matches[1], ".")
	if len(versionParts) == 0 {
		return 0, fmt.Errorf("versao do Java invalida: %s", matches[1])
	}
	if versionParts[0] == "1" && len(versionParts) > 1 {
		return strconv.Atoi(versionParts[1])
	}
	return strconv.Atoi(versionParts[0])
}

func javaExecutableName() string {
	if runtime.GOOS == "windows" {
		return "java.exe"
	}
	return "java"
}

func javacExecutableName() string {
	if runtime.GOOS == "windows" {
		return "javac.exe"
	}
	return "javac"
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
		if err := validateJDK21(candidate); err == nil {
			return candidate, nil
		}
	}

	return "", errors.New("nenhum JDK 21 valido encontrado em ~/.hubsaude/jdk")
}

func jdkDownloadInfo() (downloadURL, archiveType, archiveName string, err error) {
	switch runtime.GOOS {
	case "windows":
		if runtime.GOARCH != "amd64" {
			return "", "", "", fmt.Errorf("plataforma nao suportada: windows/%s", runtime.GOARCH)
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
			return "", "", "", fmt.Errorf("plataforma nao suportada: linux/%s", runtime.GOARCH)
		}
		archiveType = "tar.gz"
	case "darwin":
		switch runtime.GOARCH {
		case "amd64":
			archiveName = "OpenJDK21U-jdk_x64_mac_hotspot_21.0.4_7.tar.gz"
		case "arm64":
			archiveName = "OpenJDK21U-jdk_aarch64_mac_hotspot_21.0.4_7.tar.gz"
		default:
			return "", "", "", fmt.Errorf("plataforma nao suportada: darwin/%s", runtime.GOARCH)
		}
		archiveType = "tar.gz"
	default:
		return "", "", "", fmt.Errorf("plataforma nao suportada: %s/%s", runtime.GOOS, runtime.GOARCH)
	}

	downloadURL = fmt.Sprintf("%s/%s", downloadBaseURL, archiveName)
	return downloadURL, archiveType, archiveName, nil
}

func downloadFile(url, dest string) error {
	client := http.Client{Timeout: 10 * time.Minute}
	resp, err := client.Get(url)
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
		}
	}

	return nil
}

func sanitizeExtractPath(dest, filePath string) (string, error) {
	cleaned := filepath.Clean(filepath.Join(dest, filePath))
	destClean := filepath.Clean(dest)
	if !strings.HasPrefix(cleaned, destClean+string(os.PathSeparator)) && cleaned != destClean {
		return "", fmt.Errorf("caminho de extracao invalido: %s", filePath)
	}
	return cleaned, nil
}
