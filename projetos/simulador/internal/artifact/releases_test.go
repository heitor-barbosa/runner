package artifact

import (
	"strings"
	"testing"
)

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
			{Name: "sha256sums.txt", URL: "https://example.com/sha256sums.txt"},
		},
	}

	asset := FindChecksumsAsset(release)
	if asset == nil {
		t.Fatal("expected to find checksums asset")
	}
	if !strings.HasSuffix(asset.Name, "sha256sums.txt") {
		t.Fatalf("expected sha256sums.txt, got %s", asset.Name)
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
