# Status da Sprint 4

## Situacao

A Sprint 4 esta funcionalmente concluida conforme o `PlanejamentoFinal.md`.

A entrega principal foi o CLI `simulador`, responsavel por gerenciar o ciclo de
vida do `simulador.jar`, obter o artefato dinamicamente quando necessario e
participar do mesmo fluxo de CI/CD e release do CLI `assinatura`.

O estado real atual e: o projeto Go `projetos/simulador` existe, os comandos
`start`, `stop` e `status` estao implementados, o JAR pode ser resolvido por
arquivo local, cache em `~/.hubsaude/` ou GitHub Releases, e o fluxo de release
inclui binarios do `simulador` junto aos artefatos verificaveis por SHA-256 e
Cosign.

Ainda ha pontos que podem ser reforcados fora do escopo principal desta sprint:
teste end-to-end com um `simulador.jar` real respondendo a health/readiness,
validacao de prontidao apos o start e verificacao por checksum/Cosign tambem
para URLs alternativas informadas por `--source`.

## Historias concluidas

### US-03.1 - Iniciar o Simulador via CLI

- [x] Comando `simulador start` implementado.
- [x] O comando resolve o `simulador.jar` localmente ou baixa o artefato quando necessario.
- [x] A porta e validada antes da inicializacao.
- [x] O processo e iniciado em background via `java -jar`.
- [x] PID, porta, caminho do JAR e data de inicio sao registrados em `~/.hubsaude/`.
- [x] Feedback e exibido ao usuario apos a inicializacao.

### US-03.2 - Parar e monitorar o Simulador

- [x] Comando `simulador stop` implementado.
- [x] Comando `simulador status` implementado.
- [x] O status consulta o arquivo de estado em `~/.hubsaude/`.
- [x] O status confirma se o processo registrado ainda esta ativo.
- [x] O stop encerra o PID registrado e remove o estado local.
- [x] Testes cobrem status ativo, estado ausente e parada de processo registrado.

### US-03.3 - Estrutura base do CLI `simulador` em Go

- [x] Projeto Go criado em `projetos/simulador`.
- [x] Estrutura de comandos segue o padrao usado no CLI `assinatura`.
- [x] Comandos `start`, `stop`, `status` e `version` registrados.
- [x] Testes Go cobrem comandos e componentes internos.
- [x] Workflow de CI executa testes do modulo `simulador`.
- [x] Workflow de build gera binarios multiplataforma para o `simulador`.

### US-03.4 - Obter `simulador.jar` dinamicamente

- [x] O CLI consulta o GitHub Releases para identificar a release mais recente.
- [x] O artefato `simulador.jar` e localizado entre os assets da release.
- [x] O download automatico e usado quando o JAR nao esta disponivel localmente.
- [x] O cache local usa `~/.hubsaude/simulador.jar`.
- [x] A opcao `--source <url>` permite indicar URL alternativa para download.
- [x] O fluxo padrao por GitHub Releases valida checksum SHA-256.
- [x] O fluxo padrao por GitHub Releases valida assinatura Cosign com `.sig` e `.pem`.
- [x] Testes cobrem localizacao de assets, parsing de checksum e chamada de verificacao Cosign.

## Resultado entregue

- CLI `simulador` com comandos `start`, `stop`, `status` e `version`.
- Gerenciamento de ciclo de vida do `simulador.jar`.
- Registro de estado em `~/.hubsaude/`.
- Verificacao de disponibilidade de porta antes do start.
- Download/cache do `simulador.jar`.
- Consulta ao GitHub Releases para obter o artefato mais recente.
- Verificacao de checksum SHA-256 e assinatura Cosign no fluxo padrao de release.
- Testes Go para comandos, lifecycle e artifact resolution.
- CI/CD publicando binarios do `simulador` junto com o CLI `assinatura`.

## Evidencias

- Codigo do CLI:
  - `projetos/simulador/cmd/`
  - `projetos/simulador/internal/lifecycle/`
  - `projetos/simulador/internal/artifact/`
- Testes:
  - `projetos/simulador/cmd/start_test.go`
  - `projetos/simulador/cmd/lifecycle_commands_test.go`
  - `projetos/simulador/internal/lifecycle/lifecycle_test.go`
  - `projetos/simulador/internal/artifact/releases_test.go`
- CI/CD:
  - `.github/workflows/assinatura.yml`

## Como validar localmente

No modulo do simulador:

```bash
cd projetos/simulador
go test ./...
go run . --help
go run . version
go run . start --help
go run . status --help
go run . stop --help
```

Com um `simulador.jar` disponivel localmente ou em release:

```bash
go run . start --port 8081
go run . status --port 8081
go run . stop --port 8081
```

## Limitacoes conhecidas

- O start registra o processo iniciado, mas ainda nao aguarda um endpoint real
  de readiness do `simulador.jar`.
- O fluxo `--source` baixa uma URL alternativa diretamente; a verificacao
  completa de checksum e Cosign esta no fluxo padrao por GitHub Releases.
- O `simulador.jar` em si nao faz parte do escopo de desenvolvimento deste
  repositorio; o CLI apenas gerencia sua obtencao e execucao.
