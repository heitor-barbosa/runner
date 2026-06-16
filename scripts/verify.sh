#!/usr/bin/env bash
set -euo pipefail
repo_root="$(cd "$(dirname "$0")/.." && pwd)"

echo "=== Running Go tests: assinatura ==="
(
  cd "${repo_root}/projetos/assinador"
  go test ./...
)

echo "=== Running Go tests: simulador ==="
(
  cd "${repo_root}/projetos/simulador"
  go test ./...
)

echo "=== Running Java tests and packaging assinador.jar ==="
(
  cd "${repo_root}/projetos/assinador-java"
  mvn --batch-mode clean verify
)

echo "=== Verification complete ==="
