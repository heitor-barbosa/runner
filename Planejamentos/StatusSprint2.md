# Status da Sprint 2

## Sprint 2 Finalizada

O arquivo `cmd/jdk.go` foi removido pois não era parte do planejamento da Sprint 2 (não há comandos manuais no escopo).

### Status da Sprint 2:

✅ **US-02.1** — Simulação de criação de assinatura digital  
✅ **US-02.2** — Validação de parâmetros de criação de assinatura  
✅ **US-02.3** — Validação de parâmetros de validação de assinatura  
✅ **US-01.2** — Parsing de comandos e parâmetros no CLI  
✅ **US-01.3** — Invocação do assinador.jar no modo local  
✅ **US-01.4** — Exibição legível de resultados  
✅ **US-04.1** — Detecção e provisionamento automático do JDK  

### O que foi implementado:

- **CLI Go**: Comandos `sign` e `validate` funcionais
- **Assinador Java**: Simulação completa com validações rigorosas
- **Provisionamento JDK**: Download automático do JDK 21 (Temurin) quando necessário
- **Integração**: O runner detecta Java automaticamente e instala se ausente

### Para testar:

1. Construa o `assinador.jar` com Maven: `mvn package` (instale Maven se necessário)
2. Coloque o jar em `~/.hubsaude/assinador.jar` ou no diretório do executável
3. Execute `assinatura sign --bundle '{"resourceType":"Bundle"}' --provenance '{"resourceType":"Provenance"}' --credential-content "test" --certificate-chain '["cert"]' --timestamp 1751328001`

O provisionamento automático do JDK será acionado se Java não estiver disponível no sistema.

## Comandos para Teste

### 1. Construir o executável Go
```bash
cd projetos/assinador
go build -o assinatura.exe
```

### 2. Instalar Maven (se necessário)
```bash
# Opção 1: Usando Chocolatey (recomendado)
choco install maven -y

# Opção 2: Download manual
mkdir C:\temp
cd C:\temp
Invoke-WebRequest -Uri "https://archive.apache.org/dist/maven/maven-3/3.9.6/binaries/apache-maven-3.9.6-bin.zip" -OutFile "maven.zip"
Expand-Archive -Path "maven.zip" -DestinationPath "C:\maven"
# Adicionar C:\maven\apache-maven-3.9.6\bin ao PATH
```

### 3. Construir o projeto Java
```bash
cd projetos/assinador-java
mvn clean package
```

### 4. Copiar o JAR para o local esperado
```bash
# Copiar para o diretório do executável
copy target\assinador-1.0-SNAPSHOT.jar ..\assinador\assinador.jar

# OU copiar para o diretório global
copy target\assinador-1.0-SNAPSHOT.jar %USERPROFILE%\.hubsaude\assinador.jar
```

### 5. Testar o comando sign
```bash
cd projetos/assinador
.\assinatura.exe sign ^
  --bundle "{\"resourceType\":\"Bundle\",\"entry\":[]}" ^
  --provenance "{\"resourceType\":\"Provenance\",\"target\":[]}" ^
  --credential-content "test-key-content" ^
  --certificate-chain "[\"base64-cert1\",\"base64-cert2\"]" ^
  --timestamp 1751328001
```

### 6. Testar o comando validate
```bash
.\assinatura.exe validate ^
  --signature "eyJ0eXAiOiJKV1QiLCJhbGciOiJSUzI1NiJ9..." ^
  --certificate-chain "[\"base64-cert1\",\"base64-cert2\"]" ^
  --timestamp 1751328001
```

### 7. Verificar provisionamento automático do JDK
```bash
# Remover Java do PATH temporariamente para testar
# O sistema deve baixar e instalar automaticamente
.\assinatura.exe sign --bundle "{}" --provenance "{}" --credential-content "test" --certificate-chain "[]" --timestamp 1751328001
```