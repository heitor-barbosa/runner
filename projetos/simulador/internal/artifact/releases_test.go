package artifact

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func assertSameFile(t *testing.T, expectedPath string, actualPath string) {
	t.Helper()

	expectedInfo, err := os.Stat(expectedPath)
	if err != nil {
		t.Fatalf("failed to stat expected file %s: %v", expectedPath, err)
	}
	actualInfo, err := os.Stat(actualPath)
	if err != nil {
		t.Fatalf("failed to stat actual file %s: %v", actualPath, err)
	}
	if !os.SameFile(expectedInfo, actualInfo) {
		t.Fatalf("expected %s and %s to identify the same file", expectedPath, actualPath)
	}
}

func TestFindSimuladorJarAsset(t *testing.T) {
	release := &GitHubRelease{
		TagName: "v1.0.0",
		Assets: []ReleaseAsset{
			{Name: "simulador.jar", URL: "https://example.com/simulador.jar"},
			{Name: "assinatura-v1.0.0-linux-amd64", URL: "https://example.com/assinatura-linux"},
		},
	}

	asset := FindSimuladorJarAsset(release)
	if asset == nil {
		t.Fatal("expected to find simulador.jar asset")
	}
	if asset.Name != "simulador.jar" {
		t.Fatalf("expected simulador.jar, got %s", asset.Name)
	}
}

func TestFindSimuladorJarAssetNotFound(t *testing.T) {
	release := &GitHubRelease{
		TagName: "v1.0.0",
		Assets: []ReleaseAsset{
			{Name: "assinatura-v1.0.0-linux-amd64", URL: "https://example.com/assinatura-linux"},
		},
	}

	asset := FindSimuladorJarAsset(release)
	if asset != nil {
		t.Fatal("expected to not find simulador.jar asset")
	}
}

func TestFindChecksumsAsset(t *testing.T) {
	release := &GitHubRelease{
		TagName: "v1.0.0",
		Assets: []ReleaseAsset{
			{Name: "simulador.jar", URL: "https://example.com/simulador.jar"},
			{Name: "SHA256SUMS.txt", URL: "https://example.com/SHA256SUMS.txt"},
		},
	}

	asset := FindChecksumsAsset(release)
	if asset == nil {
		t.Fatal("expected to find checksums asset")
	}
	if asset.Name != "SHA256SUMS.txt" {
		t.Fatalf("expected SHA256SUMS.txt, got %s", asset.Name)
	}
}

func TestFindReleaseAsset(t *testing.T) {
	release := &GitHubRelease{
		TagName: "v1.0.0",
		Assets: []ReleaseAsset{
			{Name: "simulador.jar.sig", URL: "https://example.com/simulador.jar.sig"},
			{Name: "simulador.jar.pem", URL: "https://example.com/simulador.jar.pem"},
		},
	}

	asset := FindReleaseAsset(release, "simulador.jar.pem")
	if asset == nil {
		t.Fatal("expected to find simulador.jar.pem asset")
	}
	if asset.URL != "https://example.com/simulador.jar.pem" {
		t.Fatalf("unexpected asset URL: %s", asset.URL)
	}
}

func TestParseChecksumForJar(t *testing.T) {
	checksumContent := `
abc123def456  assinatura-v1.0.0-linux-amd64
def789ghi012  simulador.jar
jkl345mno678  assinatura-v1.0.0-windows-amd64
`
	reader := strings.NewReader(checksumContent)

	hash, err := parseChecksumForJar(reader)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if hash != "def789ghi012" {
		t.Fatalf("expected def789ghi012, got %s", hash)
	}
}

func TestDownloadAndVerifyJarFromReleaseRequiresChecksumAndCosign(t *testing.T) {
	tmpHome := t.TempDir()
	t.Setenv("HOME", tmpHome)
	t.Setenv("USERPROFILE", tmpHome)

	jarContent := []byte("fake simulador jar")
	sum := sha256.Sum256(jarContent)
	expectedHash := hex.EncodeToString(sum[:])

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Path {
		case "/simulador.jar":
			_, _ = w.Write(jarContent)
		case "/SHA256SUMS.txt":
			_, _ = fmt.Fprintf(w, "%s  simulador.jar\n", expectedHash)
		case "/simulador.jar.sig":
			_, _ = w.Write([]byte("signature"))
		case "/simulador.jar.pem":
			_, _ = w.Write([]byte("certificate"))
		default:
			http.NotFound(w, r)
		}
	}))
	defer server.Close()

	oldLatestRelease := getLatestReleaseFunc
	oldVerifyCosign := verifyCosignBlobFunc
	defer func() {
		getLatestReleaseFunc = oldLatestRelease
		verifyCosignBlobFunc = oldVerifyCosign
	}()

	getLatestReleaseFunc = func() (*GitHubRelease, error) {
		return &GitHubRelease{
			TagName: "v1.0.0",
			Assets: []ReleaseAsset{
				{Name: "simulador.jar", URL: server.URL + "/simulador.jar"},
				{Name: "SHA256SUMS.txt", URL: server.URL + "/SHA256SUMS.txt"},
				{Name: "simulador.jar.sig", URL: server.URL + "/simulador.jar.sig"},
				{Name: "simulador.jar.pem", URL: server.URL + "/simulador.jar.pem"},
			},
		}, nil
	}

	var cosignCalled bool
	verifyCosignBlobFunc = func(filePath string, certificatePath string, signaturePath string) error {
		cosignCalled = true
		if filepath.Base(filePath) != "simulador.jar" {
			t.Fatalf("expected simulador.jar, got %s", filePath)
		}
		if _, err := os.Stat(certificatePath); err != nil {
			t.Fatalf("expected certificate temp file: %v", err)
		}
		if _, err := os.Stat(signaturePath); err != nil {
			t.Fatalf("expected signature temp file: %v", err)
		}
		return nil
	}

	result, err := DownloadAndVerifyJarFromRelease()
	if err != nil {
		t.Fatalf("expected DownloadAndVerifyJarFromRelease to succeed, got %v", err)
	}
	if !cosignCalled {
		t.Fatal("expected Cosign verification to be called")
	}
	if _, err := os.Stat(result.Path); err != nil {
		t.Fatalf("expected downloaded jar to exist: %v", err)
	}
}

func TestDownloadAndVerifyJarFromReleaseRemovesJarWhenChecksumFails(t *testing.T) {
	tmpHome := t.TempDir()
	t.Setenv("HOME", tmpHome)
	t.Setenv("USERPROFILE", tmpHome)

	jarContent := []byte("fake simulador jar")
	wrongHash := strings.Repeat("0", 64)

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Path {
		case "/simulador.jar":
			_, _ = w.Write(jarContent)
		case "/SHA256SUMS.txt":
			_, _ = fmt.Fprintf(w, "%s  simulador.jar\n", wrongHash)
		default:
			http.NotFound(w, r)
		}
	}))
	defer server.Close()

	oldLatestRelease := getLatestReleaseFunc
	oldVerifyCosign := verifyCosignBlobFunc
	defer func() {
		getLatestReleaseFunc = oldLatestRelease
		verifyCosignBlobFunc = oldVerifyCosign
	}()

	getLatestReleaseFunc = func() (*GitHubRelease, error) {
		return &GitHubRelease{
			TagName: "v1.0.0",
			Assets: []ReleaseAsset{
				{Name: "simulador.jar", URL: server.URL + "/simulador.jar"},
				{Name: "SHA256SUMS.txt", URL: server.URL + "/SHA256SUMS.txt"},
				{Name: "simulador.jar.sig", URL: server.URL + "/simulador.jar.sig"},
				{Name: "simulador.jar.pem", URL: server.URL + "/simulador.jar.pem"},
			},
		}, nil
	}
	verifyCosignBlobFunc = func(filePath string, certificatePath string, signaturePath string) error {
		t.Fatal("Cosign verification should not run after checksum failure")
		return nil
	}

	result, err := DownloadAndVerifyJarFromRelease()
	if err == nil {
		t.Fatal("expected checksum failure")
	}
	if result != nil {
		t.Fatalf("expected nil result on checksum failure, got %+v", result)
	}
	if !strings.Contains(err.Error(), "verificacao de integridade falhou") {
		t.Fatalf("error = %q, want integrity failure", err)
	}

	cachePath, err := cacheJarPath()
	if err != nil {
		t.Fatalf("cacheJarPath returned error: %v", err)
	}
	if _, err := os.Stat(cachePath); !os.IsNotExist(err) {
		t.Fatalf("expected failed download to be removed, stat error = %v", err)
	}
}
func TestResolveLocalJarFindsJarInWorkingDirectory(t *testing.T) {
	tmpDir := t.TempDir()
	oldWd, err := os.Getwd()
	if err != nil {
		t.Fatalf("failed to get current working directory: %v", err)
	}
	defer func() {
		_ = os.Chdir(oldWd)
	}()

	if err := os.Chdir(tmpDir); err != nil {
		t.Fatalf("failed to change working directory: %v", err)
	}

	localPath := filepath.Join(tmpDir, jarName)
	if err := os.WriteFile(localPath, []byte("dummy"), 0o600); err != nil {
		t.Fatalf("failed to write local jar: %v", err)
	}

	jar, err := ResolveLocalJar()
	if err != nil {
		t.Fatalf("expected ResolveLocalJar to succeed, got %v", err)
	}
	assertSameFile(t, localPath, jar.Path)
}

func TestResolveJarWithFallbackUsesLocalJarWhenAvailable(t *testing.T) {
	tmpDir := t.TempDir()
	oldWd, err := os.Getwd()
	if err != nil {
		t.Fatalf("failed to get current working directory: %v", err)
	}
	defer func() {
		_ = os.Chdir(oldWd)
	}()

	if err := os.Chdir(tmpDir); err != nil {
		t.Fatalf("failed to change working directory: %v", err)
	}

	localPath := filepath.Join(tmpDir, jarName)
	if err := os.WriteFile(localPath, []byte("dummy"), 0o600); err != nil {
		t.Fatalf("failed to write local jar: %v", err)
	}

	oldLatestRelease := getLatestReleaseFunc
	defer func() { getLatestReleaseFunc = oldLatestRelease }()
	getLatestReleaseFunc = func() (*GitHubRelease, error) {
		t.Fatal("expected ResolveJarWithFallback to use local jar without calling release API")
		return nil, nil
	}

	jar, err := ResolveJarWithFallback()
	if err != nil {
		t.Fatalf("expected ResolveJarWithFallback to succeed, got %v", err)
	}
	assertSameFile(t, localPath, jar.Path)
}

func TestResolveJarWithFallbackDownloadsWhenLocalJarMissing(t *testing.T) {
	tmpHome := t.TempDir()
	t.Setenv("HOME", tmpHome)
	t.Setenv("USERPROFILE", tmpHome)

	jarContent := []byte("fake simulador jar")
	sum := sha256.Sum256(jarContent)
	expectedHash := hex.EncodeToString(sum[:])

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Path {
		case "/simulador.jar":
			_, _ = w.Write(jarContent)
		case "/SHA256SUMS.txt":
			_, _ = fmt.Fprintf(w, "%s  simulador.jar\n", expectedHash)
		case "/simulador.jar.sig":
			_, _ = w.Write([]byte("signature"))
		case "/simulador.jar.pem":
			_, _ = w.Write([]byte("certificate"))
		default:
			http.NotFound(w, r)
		}
	}))
	defer server.Close()

	oldLatestRelease := getLatestReleaseFunc
	oldVerifyCosign := verifyCosignBlobFunc
	defer func() {
		getLatestReleaseFunc = oldLatestRelease
		verifyCosignBlobFunc = oldVerifyCosign
	}()

	getLatestReleaseFunc = func() (*GitHubRelease, error) {
		return &GitHubRelease{
			TagName: "v1.0.0",
			Assets: []ReleaseAsset{
				{Name: "simulador.jar", URL: server.URL + "/simulador.jar"},
				{Name: "SHA256SUMS.txt", URL: server.URL + "/SHA256SUMS.txt"},
				{Name: "simulador.jar.sig", URL: server.URL + "/simulador.jar.sig"},
				{Name: "simulador.jar.pem", URL: server.URL + "/simulador.jar.pem"},
			},
		}, nil
	}
	verifyCosignBlobFunc = func(filePath string, certificatePath string, signaturePath string) error {
		return nil
	}

	jar, err := ResolveJarWithFallback()
	if err != nil {
		t.Fatalf("expected ResolveJarWithFallback to succeed, got %v", err)
	}
	if !strings.HasSuffix(jar.Path, jarName) {
		t.Fatalf("expected jar path to end with %s, got %s", jarName, jar.Path)
	}
}
