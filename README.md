# Sistema Runner

## 1. Visão Geral

O **Sistema Runner** é um trabalho prático desenvolvido para a disciplina de **Implementação e Integração de Software** do Bacharelado em Engenharia de Software da **Universidade Federal de Goiás (UFG)**.

O projeto tem interesse real para a **Secretaria de Estado de Saúde de Goiás (SES-GO)** e para a UFG, no contexto da plataforma **HubSaúde**, voltada à interoperabilidade de dados em saúde.

O objetivo principal é facilitar a execução de aplicações Java por linha de comando, permitindo que usuários utilizem essas aplicações sem conhecer detalhes técnicos de configuração, instalação ou execução do ambiente Java.

## 2. Componentes do Sistema

O Sistema Runner é composto por três elementos principais:

- **Assinatura CLI**: CLI em Go responsável por orquestrar a execução do `assinador.jar`.
- **Assinador Java**: aplicação `assinador.jar`, em Java 21, responsável por simular assinatura digital e validar parâmetros.
- **Simulador do HubSaúde**: aplicação `simulador.jar`, cujo ciclo de vida será gerenciado por um CLI próprio na Sprint 4.

## 3. Principais Funcionalidades

### Execução flexível

O CLI `assinatura` permite duas formas de invocação do assinador:

- **Modo local**: execução direta via `java -jar`.
- **Modo servidor HTTP**: comunicação com o `assinador.jar` via API HTTP.

### Provisionamento automático de JDK

O sistema:

- detecta se um JDK 21 está disponível;
- baixa automaticamente o JDK necessário caso não esteja presente;
- armazena o JDK em `~/.hubsaude/jdk/` para reuso.

### Simulação de assinatura digital

O assinador:

- valida rigorosamente os parâmetros de entrada;
- simula a criação de assinaturas digitais;
- simula a validação de assinaturas;
- retorna mensagens estruturadas em caso de erro.

### Modo servidor do assinador

O CLI `assinatura`:

- inicia o `assinador.jar` em modo servidor com `assinatura start`;
- reutiliza instância ativa quando o health check responde;
- invoca `/sign` e `/validate` por HTTP quando o servidor está disponível;
- permite fallback para modo local;
- encerra instância registrada com `assinatura stop`;
- permite timeout automático por inatividade com `--timeout`.

### Segurança e integridade

O pipeline de release gera:

- binários multiplataforma;
- `assinador.jar`;
- checksums SHA-256;
- assinaturas Cosign com Sigstore.

## 4. Arquitetura do Sistema

```text
Usuário
  |
  v
CLI assinatura / simulador
  |
  v
Aplicações Java
  |
  v
Resposta ao usuário
```

## 5. Estado Atual

Até a Sprint 3, o projeto já entrega:

- CLI `assinatura` com comandos `version`, `sign`, `validate`, `start` e `stop`;
- `assinador.jar` em Java 21 com simulação de assinatura e validação;
- invocação local do Java via `java -jar`;
- invocação HTTP para `sign` e `validate` quando o servidor está ativo;
- comando `assinatura start` para iniciar ou reutilizar o servidor HTTP;
- comando `assinatura stop` para encerrar uma instância registrada;
- timeout automático por inatividade via `assinatura start --timeout <minutos>`;
- health check HTTP em `/health`;
- validação de parâmetros e mensagens de erro estruturadas;
- suporte simulado a material criptográfico via PKCS#11;
- detecção e provisionamento automático de JDK 21;
- testes Go, testes Java e integração CLI -> JAR/HTTP no CI.

Ficam para a Sprint 4:

- implementação real do ciclo de vida do CLI `simulador`;
- download/cache do `simulador.jar`;
- health check/readiness do Simulador do HubSaúde;
- publicação do binário `simulador` junto aos artefatos de release.

## 6. Como Usar

Os fluxos das sprints estão documentados em `StatusSprints/`, e o uso detalhado do CLI está em `projetos/assinador/README.md`.

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

Com o servidor ativo, `sign` e `validate` usam HTTP por padrão. Para encerrar:

```bash
assinatura stop --port 8080
```

Para uso via release, o usuário precisa baixar o binário `assinatura` da sua plataforma e o arquivo `assinador.jar`, mantendo ambos na mesma pasta.

## 7. Contexto Acadêmico

| Campo | Informação |
| --- | --- |
| Instituição | Universidade Federal de Goiás (UFG) |
| Unidade | Instituto de Informática |
| Curso | Bacharelado em Engenharia de Software |
| Disciplina | Implementação e Integração de Software |
| Semestre | 2026/1 |

## 8. Equipe

- Brenner Rodrigues Sardinha
- Heitor Barbosa Souza

## 9. Observações

Este projeto **não implementa criptografia real**. O foco é simulação, integração, validação de parâmetros e gestão de execução conforme o escopo da disciplina.
