# Status da Sprint 2

## Situacao

A Sprint 2 esta concluida conforme o `PlanejamentoFinal.md`.

A entrega principal foi o fluxo local completo de assinatura e validacao simuladas: o usuario executa o CLI Go, o CLI invoca o `assinador.jar` via `java -jar`, e o Java valida os parametros e retorna respostas estruturadas.

## Historias concluidas

### US-02.1 - Simulacao de criacao de assinatura digital

- [x] Projeto Java base criado em `projetos/assinador-java`.
- [x] Interface `SignatureService` definida com `sign` e `validate`.
- [x] `FakeSignatureService` retorna assinatura simulada para parametros validos.
- [x] Resposta simulada segue a estrutura esperada.
- [x] Testes unitarios cobrem o cenario de sucesso.

### US-02.2 - Validacao de parametros de assinatura

- [x] Parametros obrigatorios verificados.
- [x] Mensagens de erro indicam parametro invalido e motivo.
- [x] Parametros invalidos rejeitados antes do processamento.
- [x] Testes unitarios cobrem os cenarios de validacao.

### US-02.3 - Simulacao e validacao de parametros de validacao

- [x] Parametros de validacao verificados.
- [x] Resultado predeterminado retornado com base em criterios simples.
- [x] Mensagens de erro claras para parametros invalidos.
- [x] Testes unitarios cobrem sucesso e falha.

### US-01.2 - Parsing de comandos e parametros no CLI

- [x] Comando `assinatura sign` aceita os parametros necessarios.
- [x] Comando `assinatura validate` aceita os parametros necessarios.
- [x] Ajuda dos comandos documenta flags e uso.
- [x] Parametros ausentes geram erro orientativo.
- [x] Testes cobrem registro e validacao dos comandos.

### US-01.3 - Invocacao local do assinador.jar

- [x] CLI localiza Java disponivel ou provisionado.
- [x] CLI executa `java -jar assinador.jar` com parametros mapeados em JSON.
- [x] Saida do JAR e capturada e repassada ao usuario.
- [x] Erros de execucao sao tratados com mensagens claras.
- [x] Fluxo CLI para JAR validado por testes.

### US-01.4 - Exibicao legivel de resultados

- [x] Resultado de assinatura exibido de forma legivel.
- [x] Resultado de validacao indica se a assinatura e valida.
- [x] Erros exibem codigo e mensagem.
- [x] Saida adequada para terminal.

### US-04.1 - Deteccao e provisionamento automatico do JDK

- [x] Sistema procura JDK 21 no `PATH`, `JAVA_HOME` e diretorio gerenciado.
- [x] JDK pode ser baixado automaticamente quando ausente.
- [x] JDK e armazenado em `~/.hubsaude/jdk/`.
- [x] Download nao e repetido quando ja existe JDK valido.
- [x] Testes cobrem deteccao e provisionamento.

## Resultado entregue

- Fluxo local ponta-a-ponta funcional.
- `assinador.jar` com assinatura e validacao simuladas.
- CLI `assinatura sign` e `assinatura validate`.
- Provisionamento automatico de JDK 21.
- Testes Go e Java cobrindo a entrega.

## Como validar localmente

### 1. Testar e empacotar o Java

```bash
cd projetos/assinador-java
mvn clean verify
```

O JAR esperado sera:

```text
target/assinador.jar
```

### 2. Disponibilizar o JAR para o CLI

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

### 4. Executar assinatura local

```bash
timestamp="$(date +%s)"

go run . sign --local \
  --bundle '{"resourceType":"Bundle","entry":[{}]}' \
  --provenance '{"resourceType":"Provenance","target":[{"reference":"urn:uuid:abc"}]}' \
  --credential-content 'test-key' \
  --certificate-chain '["cert1","cert2"]' \
  --timestamp "${timestamp}"
```

### 5. Executar validacao local

```bash
go run . validate --local \
  --signature-data '<valor-base64>' \
  --timestamp "${timestamp}"
```

No PowerShell, use `--%` se precisar preservar aspas de JSON inline.
