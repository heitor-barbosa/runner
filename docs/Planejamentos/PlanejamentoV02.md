# Planejamento — Iteração 1

## Objetivo da Iteração
Ter um **CLI em Go funcional** que:
- inicia e para a aplicação Java (`assinador.jar`)
- se comunica com uma API simples de assinatura (fake)
- gerencia processo + configuração local

---

#  Semana 1 — Entendimento + Arquitetura

## Atividades
- Entender o fluxo completo:
  CLI (Go) → Java (assinador.jar) → serviço de assinatura (fake)
- Definir responsabilidades:
  - CLI → orquestra/processo
  - Java → lógica de assinatura + API HTTP
- Definir contrato da API

##  Entregáveis
- Documento com:
  - fluxo do sistema
  - endpoints:
    - POST /sign
    - POST /validate
- Estrutura inicial:
  /cli-go
  /assinador-java

---

#  Semana 2 — Base do CLI + Processo

##  Atividades
- Criar CLI em Go
- Implementar comandos:
  - start
  - stop
- Executar .jar via Go (os/exec)
- Capturar PID

##  Entregáveis
- CLI que:
  - inicia o .jar
  - salva PID
  - para o processo

##  Armazenamento local
~/.hubsaude/config.json

---

#  Semana 3 — Simulador Java

## Atividades
- Criar:
  - SignatureService (interface)
  - FakeSignatureService (implementação)

##  Entregáveis
- Métodos:
  - sign(message)
  - validate(message, signature)

---

#  Semana 4 — API HTTP

##  Atividades
- Criar SignatureController
- Subir servidor HTTP

##  Endpoints
- POST /sign
- POST /validate

---

# Semana 5 — Integração CLI ↔ Java

##  Atividades
- CLI inicia o servidor Java
- CLI verifica porta disponível
- CLI salva:
  - porta
  - PID
  - caminho do .jar

---

# Semana 6 — Gerenciamento de Porta

## Atividades
- Detectar portas livres
- Evitar conflitos
- Criar fallback automático

---

# Semana 7 — Download do .jar

## Atividades
- Implementar:
  cli start --source URL
- Baixar .jar
- Salvar localmente

---

# Semana 8 — Persistência Local

## Atividades
- Criar .hubsaude/config.json
- Salvar:
  - PID
  - porta
  - versão
  - caminho do jar

---

# Semana 9 — Startup Completo

##  Atividades
- Verificar ambiente:
  - Java instalado
- Baixar dependências automaticamente
- Iniciar aplicação

---

# Semana 10 — CI/CD

## Atividades
- Configurar GitHub Actions
- Automatizar:
  - build Go
  - build Java
  - testes básicos
- Gerar artefatos
