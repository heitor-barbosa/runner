# Status da Sprint 1

## Situacao

A Sprint 1 esta concluida conforme o `PlanejamentoFinal.md`.

O foco da sprint foi estabelecer a base do CLI `assinatura`, a automacao de build e o processo de publicacao de artefatos versionados e verificaveis.

## Historias concluidas

### US-01.1 - Estrutura base do CLI em Go

- [x] Projeto Go inicializado em `projetos/assinador`.
- [x] CLI estruturado com Cobra.
- [x] Comando de versao disponivel.
- [x] Estrutura de pacotes documentada.
- [x] Aplicacao preparada para compilar em Windows, Linux e macOS.
- [x] `assinatura version` exibe a versao atual do CLI.

### US-05.1 - Pipeline CI/CD multiplataforma

- [x] GitHub Actions configurado.
- [x] Cross-compilation para `windows/amd64`, `linux/amd64` e `darwin/amd64`.
- [x] Build executado a cada push na branch principal.
- [x] Artefatos publicados como artifacts do workflow.

### US-05.2 - Publicacao de releases com SemVer

- [x] Tags seguem SemVer, por exemplo `v0.1.0`.
- [x] Workflow de release gera binarios nomeados por plataforma.
- [x] Binarios publicados automaticamente no GitHub Releases ao criar tag.
- [x] Nome dos artefatos segue `assinatura-<versao>-<os>-<arch>`.

### US-05.3 - Checksums SHA256 e assinatura com Cosign

- [x] Releases incluem `SHA256SUMS.txt`.
- [x] Artefatos assinados com Cosign.
- [x] Cada artefato possui `.sig` e `.pem`.
- [x] Processo de assinatura automatizado no pipeline.
- [x] Documentacao de verificacao com `cosign verify-blob`.

## Resultado entregue

- CLI base funcional.
- Pipeline de CI/CD multiplataforma.
- Releases versionadas.
- Checksums e assinaturas de artefatos.
- Documentacao inicial de uso e verificacao.

## Como validar

No modulo do CLI:

```bash
cd projetos/assinador
go test ./...
go run . version
go run . --help
```