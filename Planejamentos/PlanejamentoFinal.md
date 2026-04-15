# Plano de Desenvolvimento — Sistema Runner

## Premissas

- CLI desenvolvido em **Go (multiplataforma)**
- Aplicações Java em **Java (JAR executável)**
- Arquitetura:
  - CLI → orquestra execução
  - Java → validação + simulação
- Persistência local em: ~/.hubsaude/
- Estratégia: **iterativa, incremental e orientada a requisitos**
- Organização em **5 Sprints (1 semana cada)**

---

# Objetivo

Facilitar a execução de aplicações Java via CLI, ocultando complexidade de ambiente e configuração.

---

# Rastreabilidade (Baseado na Especificação)

| US | Descrição |
|----|----------|
| US-01 | Invocar assinador.jar via CLI |
| US-02 | Simular assinatura digital |
| US-03 | Gerenciar simulador |
| US-04 | Provisionar JDK automaticamente |
| US-05 | Distribuição multiplataforma |

---

# Sprint 1 — Fundação do Projeto

## Objetivo
Criar base do sistema (CLI + Java + estrutura)

## Valor entregue
Projeto compilando e organizado

---

### US-03.1 — Estrutura do CLI

- [ ] Inicializar projeto Go (`go mod init`)
- [ ] Estrutura:
    /cli-assinatura
    /cli-simulador
- [ ] Comando `version`

---

### US-02.1 — Base do assinador.jar

- [ ] Projeto Java criado
- [ ] Interface `SignatureService`
- [ ] Classe `FakeSignatureService`
- [ ] Métodos:
- [ ] `sign`
- [ ] `validate`
- [ ] Build do `.jar`

---

### Setup do Projeto

- [ ] README inicial
- [ ] Estrutura do repositório definida
- [ ] Convenção de commits

---

# Sprint 2 — Assinatura (Modo Local)

## Objetivo
Implementar fluxo completo via CLI local

## Valor entregue
Usuário já consegue assinar e validar via terminal

---

### US-01.1 — Invocar assinador.jar (modo local)

- [ ] CLI executa: java -jar assinador.jar
- [ ] Passa parâmetros corretamente
- [ ] Captura resposta

---

### US-02.2 — Validação de parâmetros

- [ ] Validar todos os parâmetros de entrada
- [ ] Mensagens de erro claras
- [ ] Rejeitar entradas inválidas

---

### US-02.3 — Simulação de assinatura

- [ ] Retornar assinatura fake
- [ ] Estrutura padronizada de resposta

---

### US-02.4 — Simulação de validação

- [ ] Retornar resultado (válido/inválido)
- [ ] Lógica simples

---

### CLI (assinatura)

- [ ] `sign`
- [ ] `validate`
- [ ] `--help`
- [ ] Saída formatada

---

# Sprint 3 — Modo Servidor (HTTP)

## Objetivo
Implementar execução contínua via HTTP

## Valor entregue
Melhor performance (sem cold start)

---

### US-01.2 — API HTTP no assinador

- [ ] `POST /sign`
- [ ] `POST /validate`
- [ ] JSON entrada/saída

---

### US-01.3 — Start automático

- [ ] CLI inicia servidor
- [ ] Porta padrão configurada

---

### US-01.4 — Reutilizar instância

- [ ] Detectar instância ativa
- [ ] Reutilizar processo existente

---

### US-01.5 — CLI usa HTTP por padrão

- [ ] Enviar requisições HTTP
- [ ] Fallback para modo local

---

### US-01.6 — Stop do servidor

- [ ] `cli stop`
- [ ] Encerrar processo

---

### US-01.7 — Timeout automático

- [ ] Encerrar por inatividade
- [ ] Parâmetro de tempo

---

# Sprint 4 — Simulador + Ambiente

## Objetivo
Gerenciar simulador + ambiente Java

## Valor entregue
Sistema completo (assinador + simulador)

---

### US-03.1 — CLI simulador

- [ ] `simulador start`
- [ ] `simulador stop`
- [ ] `simulador status`

---

### US-03.2 — Download do simulador

- [ ] Buscar no GitHub Releases
- [ ] Baixar automaticamente
- [ ] Cache local

---

### US-03.3 — Gerenciar portas

- [ ] Verificar portas disponíveis
- [ ] Evitar conflitos

---

### US-04.1 — Provisionar JDK

- [ ] Detectar Java instalado
- [ ] Baixar JDK automaticamente
- [ ] Armazenar em: ~/.hubsaude/jdk/

---

### Persistência

- [ ] Criar: ~/.hubsaude/config.json
- [ ] Salvar:
- PID
- porta
- versão
- caminhos

---

# Sprint 5 — Distribuição, Segurança e Qualidade

## Objetivo
Preparar para entrega final

## Valor entregue
Sistema pronto para uso real

---

### US-05.1 — Build multiplataforma

- [ ] Windows (amd64)
- [ ] Linux (amd64)
- [ ] macOS (amd64)

---

### US-05.2 — Releases

- [ ] GitHub Releases
- [ ] Versionamento SemVer

---

### US-05.3 — Checksums

- [ ] Gerar SHA256
- [ ] Publicar com artefatos

---

### Segurança (Cosign)

- [ ] Assinar artefatos
- [ ] Gerar `.sig` e `.pem`
- [ ] Verificação funcional

---

### CI/CD

- [ ] GitHub Actions
- [ ] Build automático
- [ ] Testes

---

### Testes

- [ ] Unitários (Go e Java)
- [ ] Integração (CLI ↔ Java)
- [ ] Cenários de erro

---

### Documentação

- [ ] Manual de usuário
- [ ] Guia de instalação
- [ ] Exemplos de uso

---

# Fluxos Implementados

## Criação de Assinatura
```bash
Usuário → CLI → assinador.jar → CLI → Usuário
```
## Validação
```bash
Usuário → CLI → assinador.jar → CLI → Usuário
```
---

# Tratamento de Erros

- [ ] Capturar exceções
- [ ] Propagar erro estruturado
- [ ] Mensagens claras ao usuário

---

# Resumo dos Sprints

| Sprint | Entrega |
|--------|--------|
| 1 | Estrutura base |
| 2 | Assinatura local |
| 3 | Modo servidor |
| 4 | Simulador + JDK |
| 5 | Distribuição |

---

# Definition of Done

- [ ] CLI executa assinatura e validação
- [ ] Modo local e HTTP funcionando
- [ ] Simulador gerenciado via CLI
- [ ] JDK provisionado automaticamente
- [ ] Binários multiplataforma disponíveis
- [ ] Artefatos assinados com Cosign
- [ ] Testes cobrindo principais fluxos
- [ ] Documentação completa

---

# Fora do Escopo

- [ ] Assinatura criptográfica real
- [ ] Certificados digitais
- [ ] Interface gráfica
- [ ] Autenticação
