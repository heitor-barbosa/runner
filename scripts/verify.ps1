$ErrorActionPreference = 'Stop'
$scriptPath = Split-Path -Parent $MyInvocation.MyCommand.Path
Set-Location (Join-Path $scriptPath '..')

Write-Host '=== Running Go tests ==='
go test .\projetos\assinador\... .\projetos\simulador\...

Write-Host '=== Running Java tests and packaging assinador.jar ==='
Set-Location .\projetos\assinador-java
mvn --batch-mode clean verify

Write-Host '=== Verification complete ==='
