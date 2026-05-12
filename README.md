# Sistema Runner

## 1. Visão Geral

O **Sistema Runner** é um trabalho prático desenvolvido para a disciplina de **Implementação e Integração de Software** do Bacharelado em Engenharia de Software (2026) da **Universidade Federal de Goiás (UFG)**. Este projeto é de interesse real da **Secretaria de Estado de Saúde de Goiás (SES)** e da UFG, que realizam um esforço conjunto na construção da plataforma **HubSaúde**, voltada à interoperabilidade de dados em saúde.

O objetivo principal do sistema é **facilitar a execução de aplicações Java via linha de comandos**, permitindo que usuários utilizem essas aplicações sem a necessidade de conhecer detalhes técnicos de configuração, instalação ou execução do ambiente Java.

---

## 2. Componentes do Sistema

O Sistema Runner é composto por três elementos principais:

- **Assinatura (CLI)**  
  Interface de linha de comando, desenvolvida em Go, simples e multiplataforma (Windows, Linux e macOS), responsável por orquestrar a execução das aplicações Java.

- **Assinador (Java)**  
  Aplicação `assinador.jar` responsável por realizar a **simulação de assinatura digital** e a **validação rigorosa de parâmetros**, conforme especificações definidas.

- **Simulador do HubSaúde**  
  Aplicação `simulador.jar`, cujo ciclo de vida (iniciar, parar, monitorar) é gerenciado pelo CLI do sistema.

---

## 3. Principais Funcionalidades

### Execução Flexível
O CLI permite duas formas de invocação do assinador:
- **Modo Local (CLI)**: execução direta via `java -jar`
- **Modo Servidor (HTTP)**: comunicação via API HTTP, reduzindo latência

### Provisionamento Automático de JDK
O sistema:
- Detecta se o Java está instalado
- Baixa automaticamente o JDK necessário caso não esteja presente
- Configura o ambiente sem intervenção do usuário

### Simulação de Assinatura Digital
O assinador:
- Valida rigorosamente os parâmetros de entrada
- Simula a criação de assinaturas digitais
- Simula a validação de assinaturas
- Retorna mensagens claras em caso de erro

### Gerenciamento do Simulador
O CLI permite:
- Iniciar o simulador
- Parar o simulador
- Consultar status
- Baixar automaticamente o `simulador.jar` do GitHub Releases

### Segurança e Integridade
- Binários distribuídos com **checksums SHA256**
- Assinatura criptográfica via **Cosign (Sigstore)**
- Garantia de autenticidade e integridade dos artefatos

---

## 4. Arquitetura do Sistema
```bash
Usuário
↓
CLI (assinatura / simulador)
↓
Assinador (Java)
↓
Resposta (assinatura ou validação)
```

---

## 5. Estado Atual

Até a Sprint 2, o projeto já entrega:

- CLI `assinatura` com comandos `version`, `sign` e `validate`
- `assinador.jar` em Java 21 com simulação de assinatura e validação
- Invocação local do Java via `java -jar`
- Validação de parâmetros e mensagens de erro estruturadas
- Detecção/provisionamento automático de JDK 21
- Testes Go, testes Java e integração CLI → JAR no CI

Ficam para as próximas sprints:

- Execução em modo servidor HTTP
- Gerenciamento de processos (`start`, `stop`, `status`)
- Detecção automática de instância ativa
- Timeout automático por inatividade
- Integração com material criptográfico e simulador dedicado

---

## 6. Como Usar

O fluxo local da Sprint 2 está documentado em `projetos/assinador/README.md` e em `Planejamentos/StatusSprint2.md`.

Resumo do uso:

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
## 7. Contexto Acadêmico

| Campo | Informação |
|------|--------|
| Instituição | Universidade Federal de Goiás (UFG) |
| Unidade | Instituto de Informática |
| Curso | Bacharelado em Engenharia de Software |
| Disciplina | Implementação e Integração de Software (INF0466) |
| Professor | Fabio Nogueira de Lucena |
| Semestre | 2026/1 |

---

## 11. Equipe

- Brenner Rodrigues Sardinha  
- Heitor Barbosa Souza

---

## 12. Observações

Este projeto **não implementa criptografia real**, focando exclusivamente na **simulação e validação de parâmetros**, conforme definido no escopo da disciplina.
