# Sistema Runner

## 1. Visao geral

O **Sistema Runner** e um trabalho pratico desenvolvido para a disciplina de
**Implementacao e Integracao de Software** do Bacharelado em Engenharia de
Software da **Universidade Federal de Goias (UFG)**.

O projeto tem interesse real para a **Secretaria de Estado de Saude de Goias
(SES-GO)** e para a UFG, no contexto da plataforma **HubSaude**, voltada a
interoperabilidade de dados em saude.

O objetivo principal e facilitar a execucao de aplicacoes Java por linha de
comando, permitindo que usuarios utilizem essas aplicacoes sem conhecer detalhes
de configuracao, instalacao ou execucao do ambiente Java.

## 2. Componentes do sistema

O Sistema Runner e composto por tres elementos principais:

- **Assinatura CLI**: CLI em Go responsavel por orquestrar a execucao do
  `assinador.jar`.
- **Assinador Java**: aplicacao `assinador.jar`, em Java 21, responsavel por
  simular assinatura digital, validar parametros e expor endpoints HTTP.
- **Simulador CLI**: CLI em Go responsavel por iniciar, parar, consultar status
  e obter dinamicamente o `simulador.jar`.

O `simulador.jar` em si nao faz parte do escopo de desenvolvimento deste
repositorio; o Runner gerencia sua obtencao e execucao.

## 3. Funcionalidades principais

### Assinatura e validacao

O CLI `assinatura` permite duas formas de invocacao do assinador:

- **Modo local**: execucao direta via `java -jar`.
- **Modo servidor HTTP**: comunicacao com o `assinador.jar` via API HTTP.

O assinador:

- valida os parametros de entrada;
- simula a criacao de assinaturas digitais;
- simula a validacao de assinaturas;
- retorna mensagens estruturadas em caso de erro.

### Modo servidor do assinador

O CLI `assinatura`:

- inicia o `assinador.jar` em modo servidor com `assinatura start`;
- reutiliza instancia ativa quando o health check responde;
- invoca `/sign` e `/validate` por HTTP quando o servidor esta disponivel;
- permite fallback para modo local;
- encerra instancia registrada com `assinatura stop`;
- permite timeout automatico por inatividade com `--timeout`.

### Provisionamento automatico de JDK

O sistema:

- detecta se um JDK 21 esta disponivel;
- baixa automaticamente o JDK necessario caso nao esteja presente;
- armazena o JDK em `~/.hubsaude/jdk/` para reuso.

### Simulador do HubSaude

O CLI `simulador`:

- inicia o `simulador.jar` com `simulador start`;
- verifica disponibilidade de porta antes de iniciar;
- registra PID, porta e caminho do JAR em `~/.hubsaude/`;
- consulta o estado com `simulador status`;
- encerra o processo registrado com `simulador stop`;
- baixa/cacheia o `simulador.jar` quando ele nao esta disponivel localmente;
- verifica checksum SHA-256 e assinatura Cosign no fluxo padrao por GitHub
  Releases.

### Seguranca e integridade

O pipeline de release gera:

- binarios multiplataforma dos CLIs `assinatura` e `simulador`;
- `assinador.jar`;
- `simulador.jar`, quando o artefato estiver disponivel para publicacao;
- checksums SHA-256;
- assinaturas Cosign com Sigstore.

## 4. Arquitetura

```text
Usuario
  |
  v
CLI assinatura / simulador
  |
  v
Aplicacoes Java
  |
  v
Resposta ao usuario
```

## 5. Estado atual

O projeto entrega:

- CLI `assinatura` com comandos `version`, `sign`, `validate`, `start` e `stop`;
- `assinador.jar` em Java 21 com simulacao de assinatura e validacao;
- invocacao local via `java -jar`;
- invocacao HTTP para `sign` e `validate` quando o servidor esta ativo;
- servidor HTTP do assinador com `/health`, `/sign` e `/validate`;
- timeout automatico por inatividade no servidor do assinador;
- validacao de parametros e mensagens de erro estruturadas;
- suporte simulado a material criptografico via PKCS#11;
- deteccao e provisionamento automatico de JDK 21;
- CLI `simulador` com `start`, `stop`, `status` e `version`;
- download/cache do `simulador.jar`;
- verificacao de checksum SHA-256 e assinatura Cosign no download por release;
- testes Go, testes Java e integracao CLI -> JAR/HTTP no CI;
- publicacao de binarios multiplataforma via GitHub Actions.

## 5.1. Qualidade, rastreabilidade e validacao

O projeto possui os seguintes artefatos de apoio:

- Especificacao unica de requisitos: `docs/Especificacao/especificacao.md`
- Referencia estavel ancorada em tag fixa:
  https://github.com/heitor-barbosa/runner/blob/v1.1.0/docs/Especificacao/especificacao.md
- Decisoes registradas em ADR:
  `docs/ADR/0001-architecture-and-ci-decisions.md`
- Rastreabilidade entre requisito, PR/commit, codigo e testes:
  `docs/rastreabilidade.md`
- Status das sprints:
  - `docs/StatusSprints/StatusSprint1.md`
  - `docs/StatusSprints/StatusSprint2.md`
  - `docs/StatusSprints/StatusSprint3.md`
  - `docs/StatusSprints/StatusSprint4.md`
- Build e verificacao por scripts:
  - Linux/macOS: `./scripts/verify.sh`
  - Windows PowerShell: `.\scripts\verify.ps1`
- Pipeline CI multiplataforma em `.github/workflows/assinatura.yml`

## 6. Como validar

### Verificacao completa

```bash
./scripts/verify.sh
```

No Windows PowerShell:

```powershell
.\scripts\verify.ps1
```

### Testes por modulo

```bash
cd projetos/assinador
go test ./...

cd ../simulador
go test ./...

cd ../assinador-java
mvn --batch-mode clean verify
```

## 7. Como usar

Os fluxos detalhados das sprints estao documentados em `docs/StatusSprints/`.
O uso detalhado do CLI `assinatura` esta em `projetos/assinador/README.md`.

Resumo do fluxo local:

```bash
assinatura sign \
  --bundle '{"resourceType":"Bundle","entry":[{}]}' \
  --provenance '{"resourceType":"Provenance","target":[{"reference":"urn:uuid:abc"}]}' \
  --credential-content 'test-key' \
  --certificate-chain '["cert1","cert2"]' \
  --timestamp <timestamp-atual>

assinatura validate \
  --signature-data '<valor-base64>' \
  --timestamp <mesmo-timestamp-usado-no-sign>
```

Para iniciar o `assinador.jar` em modo servidor HTTP:

```bash
assinatura start --port 8080 --timeout 15
```

Com o servidor ativo, `sign` e `validate` usam HTTP por padrao. Para encerrar:

```bash
assinatura stop --port 8080
```

Para gerenciar o simulador:

```bash
simulador start --port 8081
simulador status --port 8081
simulador stop --port 8081
```

Para uso via release, o usuario precisa baixar o binario da sua plataforma e o
JAR correspondente, mantendo ambos na mesma pasta quando aplicavel.

## 8. Contribuicao

O fluxo esperado de contribuicao e:

1. Criar uma branch a partir da branch principal.
2. Fazer commits pequenos e focados.
3. Rodar os testes locais ou `scripts/verify`.
4. Abrir PR com referencia ao requisito, sprint ou ajuste documental.
5. Aguardar o CI antes do merge.

## 9. Limitacoes conhecidas

- O projeto nao implementa criptografia real; o foco e simulacao, integracao,
  validacao de parametros e gestao de execucao.
- A integracao PKCS#11 e simulada nos testes versionados; execucao com SoftHSM2
  real pode ser adicionada como reforco futuro.
- O `simulador start` registra o processo iniciado, mas ainda nao aguarda um
  endpoint real de readiness do `simulador.jar`.
- O fluxo `simulador start --source` baixa uma URL alternativa diretamente; a
  verificacao completa de checksum e Cosign esta no fluxo padrao por GitHub
  Releases.

## 10. Contexto academico

| Campo | Informacao |
| --- | --- |
| Instituicao | Universidade Federal de Goias (UFG) |
| Unidade | Instituto de Informatica |
| Curso | Bacharelado em Engenharia de Software |
| Disciplina | Implementacao e Integracao de Software |
| Semestre | 2026/1 |

## 11. Equipe

- Brenner Rodrigues Sardinha
- Heitor Barbosa Souza
