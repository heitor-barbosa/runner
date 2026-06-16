#!/usr/bin/env bash
set -euo pipefail
cd "$(dirname "$0")/.."

echo "=== Running Go tests ==="
go test ./projetos/assinador/... ./projetos/simulador/...

echo "=== Running Java tests and packaging assinador.jar ==="
cd projetos/assinador-java
mvn --batch-mode clean verify

echo "=== Verification complete ==="
