# Status da Sprint 3

## Situacao

A Sprint 3 esta parcialmente concluida.

O estado real atual e: o servidor HTTP do `assinador.jar` existe, o CLI consegue iniciar e reutilizar uma instancia ativa, e os comandos `sign` e `validate` ja usam HTTP quando o servidor esta disponivel, com fallback para modo local. Ainda faltam `stop`, timeout por inatividade e integracao PKCS#11.

## Historias concluidas

### US-02.4 - Endpoints HTTP do assinador.jar

- [x] `SignatureController` implementado com `POST /sign` e `POST /validate`.
- [x] Endpoints reutilizam `FakeSignatureService`, `SignRequestValidator` e `ValidateRequestValidator`.
- [x] Respostas HTTP seguem a estrutura `success`, `data`, `errorCode` e `errorMessage`.
- [x] Testes de integracao validam sucesso, falha de validacao e metodo HTTP invalido.

### US-01.5 - Iniciar assinador.jar no modo servidor

- [x] Comando `assinatura start` inicia o `assinador.jar` em background.
- [x] Porta padrao `8080` suportada.
- [x] PID, porta, caminho do Java e caminho do JAR registrados em `~/.hubsaude/`.
- [x] Feedback exibido ao usuario quando o servidor inicia ou quando uma instancia ativa e reutilizada.
- [x] Parametro `--port` permite personalizar a porta.

### US-01.6 - Invocar assinador.jar via HTTP

- [x] CLI envia requisicoes HTTP para `/sign` e `/validate`.
- [x] Modo servidor e usado por padrao quando ha instancia ativa.
- [x] Fallback automatico para modo local quando o servidor nao esta disponivel.
- [x] Flag `--local` permite forcar a invocacao via `java -jar`.
- [x] Testes cobrem HTTP ativo, fallback local e bypass do HTTP com `--local`.

### US-01.7 - Detectar instancia do assinador.jar em execucao

- [x] CLI consulta o estado em `~/.hubsaude/`.
- [x] Health check HTTP em `/health` confirma se a instancia responde.
- [x] Instancia ativa e reutilizada pelo `start`.
- [x] Registro sem resposta e tratado como inativo no fluxo de inicializacao.

## Historias pendentes

### US-01.8 - Interromper execucao do assinador.jar

- [ ] Comando `assinatura stop`.
- [ ] Parametro `--port` para escolher a porta a encerrar.
- [ ] Feedback de encerramento.
- [ ] Atualizacao/remocao do registro em `~/.hubsaude/`.

### US-01.9 - Timeout por inatividade

- [ ] Parametro `--timeout <minutos>`.
- [ ] Encerramento automatico apos periodo sem requisicoes.
- [ ] Documentacao do timeout no help do CLI.

### US-02.5 - Integracao PKCS#11

- [ ] Integracao com `SunPKCS11`.
- [ ] Testes com SoftHSM2 ou simulador equivalente.
- [ ] Mensagem clara quando dispositivo criptografico nao esta disponivel.
- [ ] Documentacao de setup para token/smart card.

## Resultado entregue ate agora

- `assinador.jar` funcionando em modo servidor HTTP.
- Endpoints `/health`, `/sign` e `/validate`.
- `assinatura start` com registro de estado em `~/.hubsaude/`.
- Reuso de instancia ativa por health check.
- `assinatura sign` e `assinatura validate` usando HTTP quando possivel.
- Fallback para modo local e flag `--local`.

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
go run . start --port 8080
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

### 6. Executar validacao usando HTTP quando servidor estiver ativo

```bash
go run . validate --port 8080 \
  --signature-data '<valor-base64>' \
  --timestamp "${timestamp}"
```

Para confirmar o fallback local, pare o processo Java manualmente e execute os mesmos comandos, ou use `--local` para forcar o modo local.
