package artifact

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
	"time"
)

const (
	githubOwner = "heitor-barbosa"
	githubRepo  = "runner"
	githubAPI   = "https://api.github.com/repos"
)

// GitHubRelease represents a GitHub release with minimal fields
type GitHubRelease struct {
	TagName string      `json:"tag_name"`
	Assets  []ReleaseAsset `json:"assets"`
}

// ReleaseAsset represents a release asset with name and download URL
type ReleaseAsset struct {
	Name string `json:"name"`
	URL  string `json:"browser_download_url"`
}

// ChecksumFile represents a checksum entry
type ChecksumFile struct {
	Filename string
	Hash     string
}

// GetLatestRelease fetches the latest release from GitHub
func GetLatestRelease() (*GitHubRelease, error) {
	url := fmt.Sprintf("%s/%s/%s/releases/latest", githubAPI, githubOwner, githubRepo)
	
	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Get(url)
	if err != nil {
		return nil, fmt.Errorf("falha ao consultar GitHub Releases: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("falha ao consultar GitHub Releases: HTTP %s", resp.Status)
	}

	var release GitHubRelease
	if err := json.NewDecoder(resp.Body).Decode(&release); err != nil {
		return nil, fmt.Errorf("falha ao decodificar resposta de GitHub Releases: %w", err)
	}

	return &release, nil
}

// FindSimuladorJarAsset finds the simulador.jar asset in a release
func FindSimuladorJarAsset(release *GitHubRelease) *ReleaseAsset {
	for i := range release.Assets {
		if release.Assets[i].Name == "simulador.jar" {
			return &release.Assets[i]
		}
	}
	return nil
}

// FindChecksumsAsset finds the SHA256 checksums file in a release
func FindChecksumsAsset(release *GitHubRelease) *ReleaseAsset {
	for i := range release.Assets {
		if strings.HasSuffix(release.Assets[i].Name, "sha256sums.txt") {
			return &release.Assets[i]
		}
	}
	return nil
}

// DownloadJarFromRelease downloads the simulador.jar from the latest GitHub release
func DownloadJarFromRelease() (*JarResult, error) {
	release, err := GetLatestRelease()
	if err != nil {
		return nil, err
	}

	jarAsset := FindSimuladorJarAsset(release)
	if jarAsset == nil {
		return nil, fmt.Errorf("simulador.jar nao encontrado na versao %s do GitHub Releases", release.TagName)
	}

	return DownloadJar(jarAsset.URL)
}

// VerifyChecksum verifies a downloaded file against a SHA256 checksum
func VerifyChecksum(filePath string, expectedHash string) error {
	file, err := os.ReadFile(filePath)
	if err != nil {
		return fmt.Errorf("falha ao ler arquivo para verificacao de checksum: %w", err)
	}

	hash := sha256.Sum256(file)
	actualHash := hex.EncodeToString(hash[:])

	if strings.ToLower(actualHash) != strings.ToLower(expectedHash) {
		return fmt.Errorf("checksum nao corresponde: esperado %s, obtido %s", expectedHash, actualHash)
	}

	return nil
}

// DownloadAndVerifyJarFromRelease downloads the JAR and verifies its checksum
func DownloadAndVerifyJarFromRelease() (*JarResult, error) {
	release, err := GetLatestRelease()
	if err != nil {
		return nil, err
	}

	jarAsset := FindSimuladorJarAsset(release)
	if jarAsset == nil {
		return nil, fmt.Errorf("simulador.jar nao encontrado na versao %s do GitHub Releases", release.TagName)
	}

	checksumAsset := FindChecksumsAsset(release)
	if checksumAsset == nil {
		// Download without checksum verification if checksums file not found
		return DownloadJar(jarAsset.URL)
	}

	// Download the checksums file to temp location
	tempChecksum, err := downloadChecksumFile(checksumAsset.URL)
	if err != nil {
		// If checksum file fails, still download JAR without verification
		return DownloadJar(jarAsset.URL)
	}
	defer tempChecksum.Close()

	// Parse checksum for simulador.jar
	expectedHash, err := parseChecksumForJar(tempChecksum)
	if err != nil {
		// If parsing fails, still download JAR without verification
		return DownloadJar(jarAsset.URL)
	}

	// Download the JAR
	jarResult, err := DownloadJar(jarAsset.URL)
	if err != nil {
		return nil, err
	}

	// Verify checksum
	if err := VerifyChecksum(jarResult.Path, expectedHash); err != nil {
		return nil, fmt.Errorf("verificacao de integridade falhou: %w", err)
	}

	return jarResult, nil
}

func downloadChecksumFile(url string) (io.ReadCloser, error) {
	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Get(url)
	if err != nil {
		return nil, fmt.Errorf("falha ao baixar arquivo de checksums: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		resp.Body.Close()
		return nil, fmt.Errorf("falha ao baixar arquivo de checksums: HTTP %s", resp.Status)
	}

	return resp.Body, nil
}

func parseChecksumForJar(body io.Reader) (string, error) {
	data, err := io.ReadAll(body)
	if err != nil {
		return "", fmt.Errorf("falha ao ler arquivo de checksums: %w", err)
	}

	lines := strings.Split(string(data), "\n")
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if strings.Contains(line, "simulador.jar") {
			parts := strings.Fields(line)
			if len(parts) >= 1 {
				return parts[0], nil
			}
		}
	}

	return "", fmt.Errorf("checksum para simulador.jar nao encontrado no arquivo")
}
