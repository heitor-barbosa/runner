# Assinatura CLI

CLI em Go do Sistema Runner. A partir da Sprint 2, este modulo tambem executa o fluxo local de assinatura e validacao simuladas invocando o `assinador.jar`.

## Estrutura

- `main.go`: ponto de entrada da aplicacao.
- `cmd/root.go`: comando raiz `assinatura`.
- `cmd/sign.go`: comando `assinatura sign`.
- `cmd/validate.go`: comando `assinatura validate`.
- `cmd/version.go`: comando `assinatura version` e variavel de versao usada no build.
- `internal/jdk`: deteccao e provisionamento automatico do JDK 21.
- `internal/runner`: invocacao local ou HTTP do `assinador.jar`.
- `.github/workflows/assinatura.yml`: workflow de CI/CD mantido na raiz do repositorio.

## Requisitos

- Go 1.25 ou superior.
- Java/JDK 21, detectado automaticamente ou provisionado pelo CLI.
- `assinador.jar` em um dos locais lidos pelo runner.

## Uso local

No diretorio do modulo:

```bash
cd projetos/assinador
go run . --help
go run . version
go run . --version
go run . sign --help
go run . validate --help
go run . start --help
go run . stop --help
```

Saida esperada do comando de versao em desenvolvimento:

```text
assinatura v0.1.0
```

Para executar a suite de testes:

```bash
go test ./...
```

Para gerar um binario local:

```bash
go build -o assinatura .
```

## Fluxo local com o `assinador.jar`

O runner procura o jar:

- ao lado do executavel `assinatura`;
- em `~/.hubsaude/assinador.jar`;
- no diretorio atual.

Nas releases, o uso esperado e baixar `assinatura-<versao>-<os>-<arch>` e `assinador.jar`,
colocando ambos na mesma pasta antes de executar o CLI.

Com o jar disponivel, um fluxo de assinatura pode ser executado assim:

```bash
timestamp="$(date +%s)"

go run . sign \
  --bundle '{"resourceType":"Bundle","entry":[{}]}' \
  --provenance '{"resourceType":"Provenance","target":[{"reference":"urn:uuid:abc"}]}' \
  --credential-content 'test-key' \
  --certificate-chain '["cert1","cert2"]' \
  --timestamp "${timestamp}"
```

Depois, use o valor impresso em `Signature.data (base64)`:

```bash
go run . validate \
  --signature-data '<valor-base64>' \
  --timestamp "${timestamp}"
```

No Windows PowerShell, para preservar as aspas do JSON inline, gere o binario e use `--%`:

```powershell
$timestamp = [DateTimeOffset]::UtcNow.ToUnixTimeSeconds()
$timestamp

go build -o assinatura.exe .
.\assinatura.exe --% sign --bundle "{\"resourceType\":\"Bundle\",\"entry\":[{}]}" --provenance "{\"resourceType\":\"Provenance\",\"target\":[{\"reference\":\"urn:uuid:abc\"}]}" --credential-content test-key --certificate-chain "[\"cert1\",\"cert2\"]" --timestamp <cole-aqui-o-timestamp-exibido>
```

Depois:

```powershell
.\assinatura.exe --% validate --signature-data <valor-base64> --timestamp <mesmo-timestamp-usado-no-sign>
```

## Modo servidor HTTP

A partir da Sprint 3, o CLI pode iniciar o `assinador.jar` em modo servidor:

```bash
go run . start --port 8080
```

O comando registra PID, porta, caminho do Java e caminho do JAR em `~/.hubsaude/`.
Se uma instancia ja estiver respondendo em `/health` na porta informada, o CLI reutiliza
essa instancia e nao inicia outro processo.

Para encerrar automaticamente apos um periodo sem requisicoes:

```bash
go run . start --port 8080 --timeout 15
```

Para encerrar a instancia registrada:

```bash
go run . stop --port 8080
```

Com o servidor ativo, os comandos `sign` e `validate` usam HTTP por padrao:

```bash
go run . sign --port 8080 ...
go run . validate --port 8080 ...
```

Se o servidor nao estiver disponivel, o CLI faz fallback automatico para o modo local
via `java -jar`. Para forcar o modo local mesmo com servidor ativo:

```bash
go run . sign --local ...
go run . validate --local ...
```

## PKCS#11

Para simular assinatura usando token ou smart card, informe o tipo de credencial
`TOKEN` ou `SMARTCARD` e a configuracao SunPKCS11:

```bash
go run . sign \
  --credential-type TOKEN \
  --credential-content token \
  --credential-alias assinatura \
  --pkcs11-config ./pkcs11.cfg \
  --token-label token-a \
  --bundle '{"resourceType":"Bundle","entry":[{}]}' \
  --provenance '{"resourceType":"Provenance","target":[{"reference":"urn:uuid:abc"}]}' \
  --certificate-chain '["cert1","cert2"]' \
  --timestamp "${timestamp}"
```

Se o provider `SunPKCS11`, o arquivo de configuracao ou o dispositivo nao estiverem
disponiveis, o assinador retorna erro estruturado `PKCS11.DEVICE-UNAVAILABLE`.

## CI/CD

O workflow `.github/workflows/assinatura.yml` executa automaticamente em:

- `pull_request`;
- `push` na branch `main`;
- tags SemVer no formato `v*.*.*`;
- execucao manual por `workflow_dispatch`.

Em cada execucao, o workflow roda os testes Go, valida o modulo Java, empacota o `assinador.jar`, executa um teste de integracao local e gera binarios para:

- `linux/amd64`;
- `windows/amd64`;
- `darwin/amd64`.

Os artefatos de build seguem a convencao:

```text
assinatura-<versao>-<os>-<arch>
assinatura-<versao>-windows-amd64.exe
```

Em builds que nao sao tags, a versao usada no nome do artefato segue o formato `dev-<commit-curto>`.

## Releases

Para publicar uma release, crie e envie uma tag SemVer:

```bash
git tag v0.1.0
git push origin v0.1.0
```

Ao receber a tag, o workflow cria uma GitHub Release com:

- binarios multiplataforma;
- `assinador.jar`;
- `SHA256SUMS.txt`;
- arquivos `.sig` gerados pelo Cosign;
- arquivos `.pem` gerados pelo Cosign.

## Verificacao dos artefatos

Depois de baixar os arquivos de uma release, verifique a integridade com:

```bash
sha256sum -c SHA256SUMS.txt
```

Verifique a assinatura de um binario com:

```bash
cosign verify-blob \
  --certificate assinatura-v0.1.0-linux-amd64.pem \
  --signature assinatura-v0.1.0-linux-amd64.sig \
  assinatura-v0.1.0-linux-amd64
```
