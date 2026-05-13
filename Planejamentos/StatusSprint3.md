# Status da Sprint 3

## Situacao

A primeira fatia da Sprint 3 foi concluida com foco exclusivo na US-02.4:

- `assinador.jar` agora pode iniciar em modo servidor HTTP;
- endpoint `POST /sign` exposto;
- endpoint `POST /validate` exposto;
- endpoints reutilizam `FakeSignatureService`, `SignRequestValidator` e `ValidateRequestValidator`;
- respostas mantem a estrutura JSON usada no modo CLI: `success`, `data`, `errorCode` e `errorMessage`;
- testes de integracao validam sucesso, falha de validacao e metodo HTTP invalido.

As demais historias da Sprint 3 permanecem pendentes: PKCS#11, lifecycle pelo CLI, invocacao HTTP pelo CLI, deteccao de instancia, stop e timeout por inatividade.

## Como validar localmente

### 1. Testar e empacotar o Java

No modulo Java:

```bash
cd projetos/assinador-java
mvn clean verify
```

O jar executavel esperado sera:

```text
target/assinador.jar
```

### 2. Iniciar o servidor HTTP

```bash
java -jar target/assinador.jar server --port 8080
```

Se a porta nao for informada, o servidor usa a porta padrao `8080`.

### 3. Chamar o endpoint de assinatura

Use um timestamp atual para respeitar a janela de tolerancia do validador:

```bash
timestamp="$(date +%s)"

curl -X POST http://localhost:8080/sign \
  -H "Content-Type: application/json" \
  -d "{
    \"bundle\":\"{\\\"resourceType\\\":\\\"Bundle\\\",\\\"entry\\\":[{}]}\",
    \"provenance\":\"{\\\"resourceType\\\":\\\"Provenance\\\",\\\"target\\\":[{\\\"reference\\\":\\\"urn:uuid:abc\\\"}]}\",
    \"credentialType\":\"PEM\",
    \"credentialContent\":\"test-key\",
    \"certificateChain\":\"[\\\"cert1\\\",\\\"cert2\\\"]\",
    \"referenceTimestamp\":${timestamp},
    \"strategy\":\"iat\",
    \"policyUri\":\"https://fhir.saude.go.gov.br/r4/seguranca/ImplementationGuide/br.go.ses.seguranca|0.1.2\"
  }"
```

### 4. Chamar o endpoint de validacao

Copie o valor retornado no campo `data` do `/sign` e use em:

```bash
curl -X POST http://localhost:8080/validate \
  -H "Content-Type: application/json" \
  -d "{
    \"signatureData\":\"<valor-data-retornado-no-sign>\",
    \"referenceTimestamp\":${timestamp},
    \"policyUri\":\"https://fhir.saude.go.gov.br/r4/seguranca/ImplementationGuide/br.go.ses.seguranca|0.1.2\"
  }"
```

## Validacao automatica

O comando `mvn clean verify` executa:

- testes unitarios de servico e validadores;
- testes de integracao dos endpoints HTTP;
- empacotamento do `target/assinador.jar`.
