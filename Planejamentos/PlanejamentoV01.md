# Antes da Iteração 1
## O que estamos construindo?
Temos essencialmente três coisas:

- Dois CLIs (Command Line Interfaces)
- Uma aplicação Java (assinador/simulador)

O CLI não faz o trabalho principal — apenas chama as funções de outras aplicações (assinador, hubsaude).
As responsabilidades do CLI incluem gerenciar o ciclo de vida das aplicações e facilitar seu uso.

Enquanto isso, o assinador.jar:
- faz assinatura (simula)
- possivelmente roda como servidor HTTP

# Planejamento da Iteração 1

### 1. ENTENDER

- Compreender o problema geral
- Identificar que há um simulador no lugar do dispositivo real
- Entender o fluxo: CLI → aplicação Java → serviço de assinatura

---

### 2. CÓDIGO / INTEGRAÇÃO

- Definir classes principais
  - `SignatureService`
  - `FakeSignatureService`
- Planejar como integrar os componentes (CLI ↔ Java)

---

### 3. ENTRADAS / SAÍDAS

- Definir entradas:
  - mensagem a ser assinada
- Definir saídas:
  - mensagem assinada (fake)
- Avaliar formato:
  - CLI (flags)
  - API (JSON)

---

### 4. PROTÓTIPO GO

- Criar CLI básico em Go
- Comandos iniciais:
  - `start`
  - `stop`
- Gerenciar execução do processo Java

---

### 5. CLI + PROCESSO

- Iniciar aplicação Java via CLI
- Monitorar processo (PID)
- Preparar estrutura para parar o processo

---

### 6. SIMULADOR (CLASSE)

- Criar interface:
  - `SignatureService`
- Implementar:
  - `FakeSignatureService`
- Métodos:
  - `sign`
  - `validate`

---

### 7. API (CONTROLLER)

- Criar `SignatureController`
- Endpoints:
  - `/sign`
  - `/validate`

---

### 8. PORTAS (GO)

- Verificar portas disponíveis
- Definir porta para aplicação Java
- Evitar conflitos de execução

---

### 9. PARAR PROCESSO

- Implementar comando para parar aplicação
- Encerrar processo Java corretamente

---

### 10. DATABASE (LOCAL)

- Definir armazenamento local (ex: `.hubsaude`)
- Salvar:
  - porta usada
  - PID
  - runtime Java

---

### 11. DOWNLOAD

- Baixar aplicação Java (`.jar`)
- Permitir opção:
  - `--source URL`
- Armazenar localmente

---

### 12. STARTUP

- Fluxo de inicialização:
  - verificar ambiente
  - baixar dependências se necessário
  - iniciar aplicação
- Otimizar experiência do usuário

---

### 13. CI/CD

- Configurar GitHub Actions
- Build automático
- Testes
- Geração de artefatos (CLI + Java)
