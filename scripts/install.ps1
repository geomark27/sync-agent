# ================================================================
#  Sync Agent - Instalador para Windows
#  Uso: iwr -useb https://raw.githubusercontent.com/geomark27/sync-agent/main/scripts/install.ps1 | iex
# ================================================================

$ErrorActionPreference = "Stop"

$repo       = "geomark27/sync-agent"
$asset      = "sync-agent-windows-amd64.exe"
$installDir = "$env:LOCALAPPDATA\Programs\sync-agent"
$dest       = "$installDir\sync-agent.exe"

Write-Host ""
Write-Host "  Sync Agent - Instalador" -ForegroundColor Cyan
Write-Host "  ────────────────────────────────────" -ForegroundColor Cyan
Write-Host ""

# 1. Obtener última versión desde GitHub API
Write-Host "[1/3] Buscando ultima version..." -ForegroundColor Yellow
$release = Invoke-RestMethod "https://api.github.com/repos/$repo/releases/latest"
$version = $release.tag_name
Write-Host "      version: $version" -ForegroundColor Green

# 2. Crear directorio y descargar binario
if (-not (Test-Path $installDir)) {
    New-Item -ItemType Directory -Path $installDir | Out-Null
}
Write-Host "[2/3] Descargando $asset..." -ForegroundColor Yellow
$url = "https://github.com/$repo/releases/download/$version/$asset"
Invoke-WebRequest -Uri $url -OutFile $dest
Write-Host "      descargado en $dest" -ForegroundColor Green

# 3. Agregar al PATH si no esta
Write-Host "[3/3] Verificando PATH..." -ForegroundColor Yellow
$userPath = [Environment]::GetEnvironmentVariable("PATH", "User")
if ($userPath -notlike "*$installDir*") {
    [Environment]::SetEnvironmentVariable("PATH", "$userPath;$installDir", "User")
    Write-Host "      agregado al PATH" -ForegroundColor Green
} else {
    Write-Host "      ya estaba en PATH" -ForegroundColor Green
}

# Resultado final
Write-Host ""
Write-Host "  Sync Agent $version instalado correctamente" -ForegroundColor Green
Write-Host ""
Write-Host "  Proximos pasos:" -ForegroundColor Cyan
Write-Host "    1. Reinicia tu terminal"
Write-Host "    2. Ejecuta: sync-agent init"
Write-Host "    3. Edita el config.json (gist_token, gist_id, paths)"
Write-Host "    4. Ejecuta: sync-agent"
Write-Host ""
