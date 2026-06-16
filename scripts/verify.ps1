$ErrorActionPreference = 'Stop'
$scriptPath = Split-Path -Parent $MyInvocation.MyCommand.Path
$repoRoot = Resolve-Path (Join-Path $scriptPath '..')

Write-Host '=== Running Go tests: assinatura ==='
Push-Location (Join-Path $repoRoot 'projetos\assinador')
try {
    go test ./...
}
finally {
    Pop-Location
}

Write-Host '=== Running Go tests: simulador ==='
Push-Location (Join-Path $repoRoot 'projetos\simulador')
try {
    go test ./...
}
finally {
    Pop-Location
}

Write-Host '=== Running Java tests and packaging assinador.jar ==='
Push-Location (Join-Path $repoRoot 'projetos\assinador-java')
try {
    mvn --batch-mode clean verify
}
finally {
    Pop-Location
}

Write-Host '=== Verification complete ==='
