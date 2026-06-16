# Checklist de Avaliação do Projeto

Este checklist serve para verificar se o projeto atende aos critérios fornecidos no documento de orientações.

## A. Princípios transversais

- [x] Rastreabilidade: spec → issue/PR → commit → código → teste
- [x] Referência à especificação via link fixo (commit/tag) e sem duplicar conteúdo upstream
- [x] Reprodutibilidade: clone + comando único geram build e testes verdes
- [x] Falhar bem: erros explicativos, códigos de saída coerentes, mensagens claras
- [x] Decisões registradas em ADRs curtos para escolhas não óbvias

## B. Organização do repositório

- [x] Estrutura coerente com projeto multi-módulo (CLI + JAR)
- [x] `.gitignore` adequado para stack e sem artefatos versionados (corrigido)
- [x] `LICENSE` presente e compatível com dependências
- [x] Sem documentos de especificação genéricos; documentação específica da implementação
- [x] Nomenclatura consistente em idioma, sem acentos, espaços ou misturas inconsistentes


## C. Documentação

- [ ] `README.md` claro: descrição, build, execução, testes, contribuição, status
- [ ] Referência à especificação com link fixo (commit/tag)
- [ ] ADRs curtos para decisões relevantes
- [ ] `plano.md`/`roadmap.md` só se refletirem trabalho real com datas e issues

## D. Qualidade de código

- [ ] Código claro e simples: funções curtas, responsabilidade única
- [ ] Transporte, domínio e interface com baixo acoplamento
- [ ] Contratos explícitos entre CLI e JAR documentados e testados
- [ ] Estilo da linguagem seguido e validado em CI
- [ ] Tipagem usada com intenção, não apenas decorativa
- [ ] Sem captura genérica de erros (`catch (Throwable)`, etc.)
- [ ] Logs estruturados em vez de `print`/`System.out`
- [ ] Sem segredos, caminhos absolutos, IPs ou portas hardcoded fora de configuração
- [ ] UTF-8 declarado e controle de line endings (`.gitattributes`)

## E. Requisitos funcionais e de integração

### E1. Invocação local do `assinador.jar`
- [ ] Executáveis funcionam independentemente do diretório atual
- [ ] Argumentos preservam espaços, acentos e aspas
- [ ] `stdout` separado de `stderr` e exit code propagado

### E2. Invocação via HTTP (modo servidor)
- [ ] Start idempotente com health check real para instância viva
- [ ] Porta padrão configurável e erro claro se porta ocupada
- [ ] Shutdown controlado por endpoint ou sinal, em qualquer porta indicada
- [ ] Auto-shutdown por inatividade com janela configurável e reinício do timer por requisição
- [ ] Modo servidor padrão; modo local explícito
- [ ] Tratamento claro de timeout, conexão recusada e resposta malformada

### E3. Validação de parâmetros
- [ ] Validação feita no `assinador.jar` como autoridade única
- [ ] Mensagens distinguem erro do usuário de erro do sistema
- [ ] Códigos de saída diferentes para cada tipo de erro

### E4. Simulador do HubSaúde
- [ ] Ciclo de vida start/stop/status com health check e readiness
- [ ] Pronto para receber requisição é testado, não apenas processo iniciado

### E5. Simulador PKCS11
- [ ] Testes de integração comprovam chamadas PKCS11 reais ao simulador

### E6. Portabilidade real
- [ ] Funciona em Windows e Linux, comprovado em CI

## F. Build, dependências, supply chain

- [ ] Build reproduzível
- [ ] Versões mínimas declaradas e verificadas em runtime com erro amigável
- [ ] Dependências mínimas e justificadas
- [ ] Distribuição do JAR como artefato único com `Main-Class` correto

## G. Testes

- [ ] Pirâmide de testes saudável: unitários, integração, end-to-end
- [ ] Testes de contrato CLI ↔ JAR com subprocesso real e HTTP real
- [ ] Cenários negativos cobertos: porta ocupada, JAR ausente, JVM ausente, timeout, payload inválido, race conditions
- [ ] Sem testes flaky tolerados ou isolados e marcados
- [ ] Relatório de cobertura publicado como sinal, não meta

## H. Engenharia de processo (Git/GitHub)

- [ ] Commits atômicos com mensagens no imperativo
- [ ] PRs pequenos, revisáveis e ligados a issues relevantes
- [ ] CI obrigatório: lint + build + testes em Windows e Linux
- [ ] Tags/releases semânticas coerentes com `release.json`
- [ ] Changelog gerado, não escrito à mão

## I. Operabilidade

- [ ] `--help` com exemplos úteis
- [ ] `--version` acessível e rastreável (tag + SHA curto)
- [ ] Logs com níveis ajustáveis (`--verbose`, `--quiet`)
