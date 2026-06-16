# Rastreabilidade da Seção A

Este documento registra a cadeia de rastreabilidade usada pelo projeto Runner:
`especificação → PR/commit → código → teste`.

## Fonte de requisitos

- Especificação principal: `docs/Especificacao/especificacao.md`
- Link fixo com tag estável: https://github.com/heitor-barbosa/runner/blob/v1.1.0/docs/Especificacao/especificacao.md

## Planejamento e status

- `docs/Planejamentos/PlanejamentoFinal.md`
- `docs/StatusSprints/StatusSprint1.md`
- `docs/StatusSprints/StatusSprint2.md`
- `docs/StatusSprints/StatusSprint3.md`

## PRs e commits relevantes

O repositório usa fluxo de PRs para mudanças significativas. Exemplos recentes:

- `4f7a970` – Merge PR #22: atualiza README e planejamento para `simulador.jar` e Cosign
- `f39f8a7` – Merge PR #21: release do simulador com Cosign e testes atualizados
- `7cc97cf` – docs de aceitação completas para stop/status do simulador
- `8a4c0f5` – implementação de stop/status com testes
- `166bc4a` – implementa timeout do servidor HTTP e adiciona testes correspondentes
- `8507f6d` – adiciona verificação de disponibilidade de porta no start
- `b1469b6` – integra HTTP e detecção de instância do `assinador.jar`
- `48eddbb` – adiciona comando stop para `assinador.jar`

## Código

Os artefatos principais associados à seção A incluem:

- CLI de assinatura:
  - `projetos/assinador/cmd/`
  - `projetos/assinador/internal/runner/runner.go`
- Aplicação Java do assinador:
  - `projetos/assinador-java/src/main/java/br/gov/go/ses/assinador/`
- Scripts de verificação e documentação de decisão:
  - `scripts/verify.sh`
  - `scripts/verify.ps1`
  - `docs/ADR/0001-architecture-and-ci-decisions.md`

## Testes

Os testes que atestam a cadeia são:

- `go test ./...` nos módulos `projetos/assinador` e `projetos/simulador`
- `mvn --batch-mode clean verify` em `projetos/assinador-java`
- Integração CLI → JAR e HTTP via workflow:
  - `.github/workflows/assinatura.yml`

## Observações

- O projeto possui um sistema de issues externo documentado no repositório, e o histórico de commits e PRs fornece o vínculo entre requisito e implementação.
- A especificação interna é tratada como fonte única de verdade para este trabalho prático.
