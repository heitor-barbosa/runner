package artifact

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/exec"
	"strings"
	"time"
)

const (
	githubOwner = "heitor-barbosa"
	githubRepo  = "runner"
	githubAPI   = "https://api.github.com/repos"
)

var (
	getLatestReleaseFunc = getLatestReleaseFromGitHub
	verifyCosignBlobFunc = defaultVerifyCosignBlob
)

// GitHubRelease represents a GitHub release with minimal fields
type GitHubRelease struct {
	TagName string         `json:"tag_name"`
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
	return getLatestReleaseFunc()
}

func getLatestReleaseFromGitHub() (*GitHubRelease, error) {
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
		if strings.EqualFold(release.Assets[i].Name, "SHA256SUMS.txt") ||
			strings.HasSuffix(strings.ToLower(release.Assets[i].Name), "sha256sums.txt") {
			return &release.Assets[i]
		}
	}
	return nil
}

// FindReleaseAsset finds a release asset by exact name.
func FindReleaseAsset(release *GitHubRelease, name string) *ReleaseAsset {
	for i := range release.Assets {
		if release.Assets[i].Name == name {
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

// DownloadAndVerifyJarFromRelease downloads the JAR and verifies checksum and Cosign signature.
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
		return nil, fmt.Errorf("SHA256SUMS.txt nao encontrado na versao %s do GitHub Releases", release.TagName)
	}

	tempChecksum, err := downloadChecksumFile(checksumAsset.URL)
	if err != nil {
		return nil, err
	}
	defer tempChecksum.Close()

	expectedHash, err := parseChecksumForJar(tempChecksum)
	if err != nil {
		return nil, err
	}

	jarResult, err := DownloadJar(jarAsset.URL)
	if err != nil {
		return nil, err
	}

	if err := VerifyChecksum(jarResult.Path, expectedHash); err != nil {
		_ = os.Remove(jarResult.Path)
		return nil, fmt.Errorf("verificacao de integridade falhou: %w", err)
	}

	signatureAsset := FindReleaseAsset(release, "simulador.jar.sig")
	if signatureAsset == nil {
		return nil, fmt.Errorf("assinatura Cosign simulador.jar.sig nao encontrada na versao %s do GitHub Releases", release.TagName)
	}

	certificateAsset := FindReleaseAsset(release, "simulador.jar.pem")
	if certificateAsset == nil {
		return nil, fmt.Errorf("certificado Cosign simulador.jar.pem nao encontrado na versao %s do GitHub Releases", release.TagName)
	}

	signaturePath, err := downloadAssetToTempFile(signatureAsset.URL, "simulador-*.sig")
	if err != nil {
		return nil, err
	}
	defer os.Remove(signaturePath)

	certificatePath, err := downloadAssetToTempFile(certificateAsset.URL, "simulador-*.pem")
	if err != nil {
		return nil, err
	}
	defer os.Remove(certificatePath)

	if err := verifyCosignBlobFunc(jarResult.Path, certificatePath, signaturePath); err != nil {
		_ = os.Remove(jarResult.Path)
		return nil, fmt.Errorf("verificacao Cosign falhou: %w", err)
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

func downloadAssetToTempFile(url string, pattern string) (string, error) {
	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Get(url)
	if err != nil {
		return "", fmt.Errorf("falha ao baixar artefato de verificacao: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("falha ao baixar artefato de verificacao: HTTP %s", resp.Status)
	}

	file, err := os.CreateTemp("", pattern)
	if err != nil {
		return "", fmt.Errorf("falha ao criar arquivo temporario de verificacao: %w", err)
	}
	defer file.Close()

	if _, err := io.Copy(file, resp.Body); err != nil {
		_ = os.Remove(file.Name())
		return "", fmt.Errorf("falha ao salvar artefato de verificacao: %w", err)
	}

	return file.Name(), nil
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

func defaultVerifyCosignBlob(filePath string, certificatePath string, signaturePath string) error {
	output, err := exec.Command(
		"cosign",
		"verify-blob",
		"--certificate", certificatePath,
		"--signature", signaturePath,
		filePath,
	).CombinedOutput()
	if err != nil {
		return fmt.Errorf("%w: %s", err, strings.TrimSpace(string(output)))
	}
	return nil
}
