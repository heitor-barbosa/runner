# Status da Sprint 3
## Situação

A Sprint 3 está funcionalmente concluída nos fluxos principais do `assinador.jar` em modo servidor.

O estado real atual é: o servidor HTTP do `assinador.jar` existe, o CLI consegue iniciar, reutilizar e parar uma instância ativa, e os comandos `sign` e `validate` usam HTTP quando o servidor está disponível, com fallback para modo local. Também existe suporte simulado a PKCS#11 e timeout por inatividade.

Ainda há pontos que merecem reforço para aderir melhor aos critérios de aceitação gerais: testes mais fortes para porta ocupada, race no start, reinício do timer de inatividade, execução Java em Windows no CI e integração PKCS#11 mais próxima de um simulador real como SoftHSM2.

## Histórias concluídas

### US-02.4 - Endpoints HTTP do `assinador.jar`

- [x] `SignatureController` implementado com `POST /sign` e `POST /validate`.
- [x] Endpoints reutilizam `FakeSignatureService`, `SignRequestValidator` e `ValidateRequestValidator`.
- [x] Respostas HTTP seguem a estrutura `success`, `data`, `errorCode` e `errorMessage`.
- [x] Endpoint `/health` implementado para health check.
- [x] Testes de integração validam sucesso, falha de validação, método HTTP inválido e health check.

### US-01.5 - Iniciar `assinador.jar` no modo servidor

- [x] Comando `assinatura start` inicia o `assinador.jar` em background.
- [x] Porta padrão `8080` suportada.
- [x] PID, porta, caminho do Java e caminho do JAR registrados em `~/.hubsaude/`.
- [x] Feedback exibido ao usuário quando o servidor inicia ou quando uma instância ativa é reutilizada.
- [x] Parâmetro `--port` permite personalizar a porta.
- [+] Falha por porta ocupada existe indiretamente pelo processo Java, mas ainda precisa de teste específico.

### US-01.6 - Invocar `assinador.jar` via HTTP

- [x] CLI envia requisições HTTP para `/sign` e `/validate`.
- [x] Modo servidor é usado por padrão quando há instância ativa.
- [x] Fallback automático para modo local quando o servidor não está disponível.
- [x] Flag `--local` permite forçar a invocação via `java -jar`.
- [x] Testes Go cobrem HTTP ativo, fallback local e bypass do HTTP com `--local`.
- [+] O CI exercita integração local CLI -> JAR, mas ainda pode reforçar integração HTTP real CLI -> JAR.

### US-01.7 - Detectar instância do `assinador.jar` em execução

- [x] CLI consulta o estado em `~/.hubsaude/`.
- [x] Health check HTTP em `/health` confirma se a instância responde.
- [x] Instância ativa é reutilizada pelo `start`.
- [x] Registro sem resposta é tratado como inativo no fluxo de inicialização.
- [+] Ainda falta teste específico para condição de corrida durante start simultâneo.

### US-01.8 - Interromper execução do `assinador.jar`

- [x] Comando `assinatura stop` implementado.
- [x] Parâmetro `--port` permite escolher a porta a encerrar.
- [x] Feedback de encerramento exibido ao usuário.
- [x] Registro em `~/.hubsaude/` é removido após encerramento.
- [x] Teste Go cobre encerramento de processo registrado e remoção do estado.
- [+] O encerramento atual usa PID salvo e `Kill`; não há endpoint HTTP de shutdown controlado.

### US-01.9 - Timeout por inatividade

- [x] Parâmetro `--timeout <minutos>` implementado no comando `assinatura start`.
- [x] Servidor Java recebe o timeout e encerra após período sem requisições.
- [x] O servidor atualiza a última interação em `/health`, `/sign` e `/validate`.
- [x] Mecanismo de timeout está documentado no help/comando e no README do módulo `assinador`.
- [+] Ainda falta teste específico comprovando que o timer reinicia a cada requisição.

### US-02.5 - Integração PKCS#11

- [x] Suporte a credenciais `TOKEN` e `SMARTCARD` no fluxo de assinatura.
- [x] Integração com provider PKCS#11 por `Pkcs11ProviderLoader`.
- [x] Mensagem estruturada para indisponibilidade de dispositivo: `PKCS11.DEVICE-UNAVAILABLE`.
- [x] Teste de integração cobre uso de provider PKCS#11 simulado.
- [x] Documentação do uso com `--pkcs11-config` e `--token-label` existe no README do módulo `assinador`.
- [+] O teste atual usa um provider simulado em memória; ainda não comprova execução com SoftHSM2 real.

## Resultado entregue

- `assinador.jar` funcionando em modo servidor HTTP.
- Endpoints `/health`, `/sign` e `/validate`.
- `assinatura start` com registro de estado em `~/.hubsaude/`.
- Reuso de instância ativa por health check.
- `assinatura sign` e `assinatura validate` usando HTTP quando possível.
- Fallback para modo local e flag `--local`.
- `assinatura stop` para encerrar instância registrada.
- Timeout por inatividade configurável com `--timeout`.
- Suporte simulado a PKCS#11 para `TOKEN` e `SMARTCARD`.
- Testes Go para comandos, runner e JDK.
- Testes Java versionados para validação, HTTP, assinatura fake e PKCS#11 simulado.
- CI com testes Go em Linux, Windows e macOS, testes Java em Ubuntu, integração local CLI -> JAR, build multiplataforma e release.

## Como validar localmente

### 1. Testar o Java

```bash
cd projetos/assinador-java
mvn test
```

### 2. Testar o CLI Go

```bash
cd projetos/assinador
go test ./...
```

### 3. Empacotar e disponibilizar o JAR

```bash
cd projetos/assinador-java
mvn clean package
cp target/assinador.jar ../assinador/assinador.jar
```

No Windows PowerShell:

```powershell
cd projetos\assinador-java
mvn clean package
Copy-Item target\assinador.jar ..\assinador\assinador.jar
```

### 4. Iniciar o servidor pelo CLI

```bash
cd ../assinador
go run . start --port 8080 --timeout 15
```

### 5. Executar assinatura usando HTTP quando servidor estiver ativo

```bash
timestamp="$(date +%s)"

go run . sign --port 8080 \
  --bundle '{"resourceType":"Bundle","entry":[{}]}' \
  --provenance '{"resourceType":"Provenance","target":[{"reference":"urn:uuid:abc"}]}' \
  --credential-content 'test-key' \
  --certificate-chain '["cert1","cert2"]' \
  --timestamp "${timestamp}"
```

### 6. Executar validação usando HTTP quando servidor estiver ativo

```bash
go run . validate --port 8080 \
  --signature-data '<valor-base64>' \
  --timestamp "${timestamp}"
```

### 7. Encerrar o servidor pelo CLI

```bash
go run . stop --port 8080
```

Para confirmar o fallback local, pare o servidor e execute os mesmos comandos, ou use `--local` para forçar o modo local.
irmar o fallback local, pare o processo Java manualmente e execute os mesmos comandos, ou use `--local` para forcar o modo local.
irmar o fallback local, pare o processo Java manualmente e execute os mesmos comandos, ou use `--local` para forcar o modo local.
