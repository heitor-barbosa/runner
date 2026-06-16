# ADR 0001 — Arquitetura do projeto e validação contínua

Status: Aceito

## Contexto

O projeto Runner integra duas peças principais:

- um CLI multiplataforma para orquestrar aplicações Java;
- um `assinador.jar` em Java 21 que oferece simulação de assinatura digital e validação.

Também há necessidade de:

- manter a solução reproduzível para qualquer desenvolvedor;
- fornecer um fluxo de execução local e um modo servidor HTTP;
- garantir compatibilidade com Windows e Linux (além do macOS em CI);
- registrar as decisões não triviais para auditoria e manutenção.

## Decisão

Adotamos a seguinte arquitetura:

1. CLI em Go 1.25
   - por ser fácil de compilar para múltiplas plataformas;
   - por oferecer controle robusto de processos e pipes;
   - por gerar binários distribuíveis sem dependências externas.

2. `assinador.jar` em Java 21
   - para manter compatibilidade com o requisito de execução de JAR Java;
   - para permitir validação de parâmetros e simulação de assinatura no runtime correto.

3. Modo servidor HTTP como padrão e modo local como escolha explícita
   - o servidor reduz o overhead de inicialização em chamadas repetidas;
   - o modo local permanece disponível para execuções isoladas ou scripts.

4. GitHub Actions multiplataforma para validação contínua
   - roda testes Go em Ubuntu, Windows e macOS;
   - roda testes Java em Ubuntu e Windows;
   - executa integração CLI → JAR e HTTP em Ubuntu.

5. Documentação de rastreabilidade e decisão
   - `docs/Especificacao/especificacao.md` é a fonte única de requisitos;
   - decisões importantes são registradas em ADRs em `docs/ADR/`;
   - rastreabilidade é documentada em `docs/rastreabilidade.md`.

6. Referência de especificação fixa
   - a especificação principal é ancorada a uma tag de release estável (`v1.1.0`) para manter a rastreabilidade.

## Consequências

- A documentação e os controles do projeto tornam explícitos quais artefatos são a fonte da verdade.
- O fluxo de validação passa a ter um comando único seguido por CI multiplataforma.
- A manutenção futura se beneficia de histórico claro entre requisito, decisão e implementação.
