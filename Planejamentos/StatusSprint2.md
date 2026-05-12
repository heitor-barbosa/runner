# Status da Sprint 2

## Situacao

A Sprint 2 entrega o fluxo local completo:

- `assinador.jar` em Java 21 com `sign` e `validate`;
- validacao rigorosa de parametros no Java;
- CLI Go com comandos `assinatura sign` e `assinatura validate`;
- invocacao local via `java -jar assinador.jar`;
- deteccao de JDK 21 e provisionamento automatico quando necessario;
- testes Go, testes Java e teste de integracao no GitHub Actions.

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

### 2. Disponibilizar o jar para o CLI

Copie o arquivo para um local lido pelo runner:

```bash
cp target/assinador.jar ../assinador/assinador.jar
```

No Windows PowerShell:

```powershell
Copy-Item target\assinador.jar ..\assinador\assinador.jar
```

### 3. Testar o CLI Go

```bash
cd ../assinador
go test ./...
```

### 4. Assinar localmente

Use um timestamp atual para respeitar a janela de tolerancia do validador:

```bash
timestamp="$(date +%s)"

go run . sign \
  --bundle '{"resourceType":"Bundle","entry":[{}]}' \
  --provenance '{"resourceType":"Provenance","target":[{"reference":"urn:uuid:abc"}]}' \
  --credential-content 'test-key' \
  --certificate-chain '["cert1","cert2"]' \
  --timestamp "${timestamp}"
```

No Windows PowerShell:

```powershell
$timestamp = [DateTimeOffset]::UtcNow.ToUnixTimeSeconds()
$timestamp

go build -o assinatura.exe .
.\assinatura.exe --% sign --bundle "{\"resourceType\":\"Bundle\",\"entry\":[{}]}" --provenance "{\"resourceType\":\"Provenance\",\"target\":[{\"reference\":\"urn:uuid:abc\"}]}" --credential-content test-key --certificate-chain "[\"cert1\",\"cert2\"]" --timestamp <cole-aqui-o-timestamp-exibido>
```

### 5. Validar a assinatura

Copie o valor exibido em `Signature.data (base64)` e execute:

```bash
go run . validate \
  --signature-data '<valor-base64>' \
  --timestamp "${timestamp}"
```

No Windows PowerShell:

```powershell
.\assinatura.exe --% validate --signature-data <valor-base64> --timestamp <mesmo-timestamp-usado-no-sign>
```

No PowerShell, `--%` impede que as aspas internas do JSON sejam reescritas pelo shell antes de chegarem ao executavel.

## Provisionamento automatico do JDK

O CLI tenta usar, nesta ordem:

1. JDK 21 no `PATH`;
2. JDK 21 em `JAVA_HOME`;
3. JDK 21 provisionado em `~/.hubsaude/jdk/`.

Se nenhum estiver disponivel, o runner baixa um Temurin JDK 21 compativel com a plataforma e salva em `~/.hubsaude/jdk/`. Instalacoes locais validas sao reutilizadas.

## Validacao automatica no CI

O workflow `.github/workflows/assinatura.yml` executa:

- `go test ./...` em Windows, Linux e macOS;
- `mvn clean verify`;
- empacotamento de `assinador.jar`;
- teste de integracao `assinatura sign` seguido de `assinatura validate`;
- builds multiplataforma do CLI.
