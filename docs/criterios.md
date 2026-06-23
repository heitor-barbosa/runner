# Checklist de Avaliacao do Projeto

Este checklist registra a aderencia do projeto aos criterios de avaliacao e
separa o que ja possui evidencia do que permanece como limitacao conhecida.

Legenda:

- `[x]` Atendido com evidencia no repositorio.
- `[ ]` Ainda nao atendido ou mantido como limitacao assumida.

## A. Principios transversais

- [x] Rastreabilidade: especificacao -> PR/commit -> codigo -> teste.
- [x] Referencia a especificacao via link fixo e sem duplicar conteudo upstream.
- [x] Reprodutibilidade por scripts de verificacao e CI.
- [x] Falhar bem: erros explicativos e mensagens claras nos fluxos principais.
- [x] Decisoes registradas em ADR curto para escolhas nao obvias.

Evidencias: `docs/rastreabilidade.md`, `docs/ADR/0001-architecture-and-ci-decisions.md`,
`scripts/verify.sh`, `scripts/verify.ps1` e `.github/workflows/assinatura.yml`.

## B. Organizacao do repositorio

- [x] Estrutura coerente com projeto multi-modulo: CLI Go, JAR Java e CLI do simulador.
- [x] `.gitignore` adequado para stack e sem artefatos de build versionados.
- [x] `LICENSE` presente.
- [x] Documentacao especifica da implementacao.
- [x] Nomenclatura consistente nos diretorios e arquivos principais.

## C. Documentacao

- [x] `README.md` cobre descricao, build, execucao, testes, contribuicao, status e limitacoes.
- [x] Referencia a especificacao com link fixo.
- [x] ADR curto para decisoes relevantes.
- [x] Planejamento e status das sprints refletem o trabalho realizado.

Evidencias: `README.md`, `docs/Planejamentos/PlanejamentoFinal.md`,
`docs/StatusSprints/` e `docs/rastreabilidade.md`.

## D. Qualidade de codigo

- [x] Codigo organizado por responsabilidades principais.
- [x] Transporte, dominio e interface separados nos modulos principais.
- [x] Contratos principais entre CLI e JAR documentados e testados.
- [x] Estilo da linguagem validado por build/testes em CI.
- [x] UTF-8 declarado no `pom.xml` e controle de line endings em `.gitattributes`.
- [ ] Logs estruturados com niveis ajustaveis.

Observacao: a operabilidade por logs estruturados nao e o foco atual; os CLIs
priorizam saida direta e mensagens de erro orientativas.

## E. Requisitos funcionais e de integracao

### E1. Invocacao local do `assinador.jar`

- [x] CLI localiza o JAR em caminhos esperados.
- [x] Argumentos sao mapeados para payload estruturado.
- [x] Saida e erros do JAR sao tratados nos fluxos principais.
- [x] Integracao CLI -> JAR e exercitada no CI.

### E2. Invocacao via HTTP

- [x] `assinatura start` reutiliza instancia ativa por health check.
- [x] Porta padrao configuravel por flag.
- [x] Modo servidor usado por padrao quando ha instancia ativa.
- [x] Modo local explicito por `--local`.
- [x] Auto-shutdown por inatividade implementado.
- [x] Integracao HTTP real exercitada no CI.
- [x] Shutdown controlado por endpoint HTTP.
- [x] Teste especifico para reinicio do timer de inatividade.
- [x] Teste especifico para race condition em starts simultaneos.

### E3. Validacao de parametros

- [x] Validacao centralizada no `assinador.jar`.
- [x] Mensagens indicam parametro invalido e motivo.
- [x] Testes Java cobrem cenarios de sucesso e falha.
- [ ] Codigos de saida distintos para todas as categorias de erro.

### E4. Simulador do HubSaude

- [x] Ciclo de vida `start`, `stop` e `status` implementado.
- [x] Porta verificada antes do start.
- [x] Estado salvo em `~/.hubsaude/`.
- [x] Download/cache do `simulador.jar` implementado.
- [x] Fluxo padrao por release valida checksum SHA-256 e Cosign.
- [ ] Readiness real do `simulador.jar` apos start.
- [ ] Verificacao completa de checksum/Cosign para URL alternativa via `--source`.

### E5. PKCS#11

- [x] Suporte a credenciais `TOKEN` e `SMARTCARD`.
- [x] Provider PKCS#11 encapsulado em `Pkcs11ProviderLoader`.
- [x] Teste versionado cobre provider simulado.
- [ ] Teste de integracao com SoftHSM2 real.

### E6. Portabilidade

- [x] Go test executado em Linux, Windows e macOS no CI.
- [x] Java test/package executado em Linux e Windows no CI.
- [x] Binarios gerados para Linux, Windows e macOS.

## F. Build, dependencias e supply chain

- [x] Build reproduzivel por scripts e CI.
- [x] Versoes principais declaradas: Go 1.25 e Java 21.
- [x] Dependencias reduzidas e explicitas nos modulos.
- [x] `assinador.jar` distribuido como artefato unico com `Main-Class`.
- [x] Releases incluem checksums SHA-256.
- [x] Releases incluem assinaturas Cosign.

## G. Testes

- [x] Testes unitarios Go.
- [x] Testes unitarios Java.
- [x] Testes de integracao CLI -> JAR.
- [x] Testes de integracao CLI -> HTTP -> JAR no CI.
- [x] Testes do CLI `simulador` para comandos, lifecycle e artifact resolution.
- [ ] Cobertura publicada como relatorio.
- [ ] Cobertura completa de cenarios negativos: race condition, readiness e SoftHSM2 real.

## H. Engenharia de processo

- [x] CI obrigatorio com testes e build multiplataforma.
- [x] Releases SemVer por tags `v*.*.*`.
- [x] Artefatos de release nomeados por plataforma.
- [ ] Changelog gerado automaticamente.
- [ ] Arquivo `release.json`.

## I. Operabilidade

- [x] `--help` disponivel nos comandos principais.
- [x] `version` e `--version` disponiveis nos CLIs.
- [x] Estado operacional registrado em `~/.hubsaude/`.
- [ ] Logs com niveis `--verbose` e `--quiet`.

## Resumo

O projeto atende ao escopo funcional das sprints planejadas. As pendencias
restantes sao reforcos de robustez, observabilidade e testes de integracao mais
proximos de ambiente real.
