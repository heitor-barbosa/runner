# Assinatura CLI

CLI em Go do Sistema Runner. Nesta sprint, este modulo entrega a base do comando `assinatura`, o comando de versao e o pipeline de build/release multiplataforma.

## Estrutura

- `main.go`: ponto de entrada da aplicacao.
- `cmd/root.go`: comando raiz `assinatura`.
- `cmd/version.go`: comando `assinatura version` e variavel de versao usada no build.
- `.github/workflows/assinatura.yml`: workflow de CI/CD mantido na raiz do repositorio.

## Requisitos

- Go 1.25 ou superior.

## Uso local

No diretorio do modulo:

```bash
cd projetos/assinador
go run . --help
go run . version
go run . --version
```

Saida esperada do comando de versao em desenvolvimento:

```text
assinatura v0.1.0
```

Para executar a suite de testes:

```bash
go test ./...
```

Para gerar um binario local:

```bash
go build -o assinatura .
```

## CI/CD

O workflow `.github/workflows/assinatura.yml` executa automaticamente em:

- `pull_request`;
- `push` na branch `main`;
- tags SemVer no formato `v*.*.*`;
- execucao manual por `workflow_dispatch`.

Em cada execucao, o workflow roda os testes e gera binarios para:

- `linux/amd64`;
- `windows/amd64`;
- `darwin/amd64`.

Os artefatos de build seguem a convencao:

```text
assinatura-<versao>-<os>-<arch>
assinatura-<versao>-windows-amd64.exe
```

Em builds que nao sao tags, a versao usada no nome do artefato segue o formato `dev-<commit-curto>`.

## Releases

Para publicar uma release, crie e envie uma tag SemVer:

```bash
git tag v0.1.0
git push origin v0.1.0
```

Ao receber a tag, o workflow cria uma GitHub Release com:

- binarios multiplataforma;
- `SHA256SUMS.txt`;
- arquivos `.sig` gerados pelo Cosign;
- arquivos `.pem` gerados pelo Cosign.

## Verificacao dos artefatos

Depois de baixar os arquivos de uma release, verifique a integridade com:

```bash
sha256sum -c SHA256SUMS.txt
```

Verifique a assinatura de um binario com:

```bash
cosign verify-blob \
  --certificate assinatura-v0.1.0-linux-amd64.pem \
  --signature assinatura-v0.1.0-linux-amd64.sig \
  assinatura-v0.1.0-linux-amd64
```
