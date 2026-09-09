$ErrorActionPreference = "Stop"

$previousGooseDriver = $env:GOOSE_DRIVER
$previousGooseDbString = $env:GOOSE_DBSTRING

Push-Location $PSScriptRoot

try {
    if (-not (Get-Command goose -ErrorAction SilentlyContinue)) {
        throw "Install Goose with: go install github.com/pressly/goose/v3/cmd/goose@latest. Ensure its installation directory is on PATH."
    }

    $securePassword = Read-Host "Local PostgreSQL password for dev" -AsSecureString
    $credential = [System.Net.NetworkCredential]::new("dev", $securePassword)
    $encodedPassword = [Uri]::EscapeDataString($credential.Password)

    $env:GOOSE_DRIVER = "postgres"
    $env:GOOSE_DBSTRING = "postgres://dev:${encodedPassword}@localhost:5432/modmapper?sslmode=disable"

    Remove-Variable credential, securePassword, encodedPassword

    Write-Host "Applying pending migrations..."
    goose -dir migrations up

    if ($LASTEXITCODE -ne 0) {
        throw "Migration failed. See Goose's error output above."
    }

    Write-Host "Migrations completed."
}
finally {
    $env:GOOSE_DRIVER = $previousGooseDriver
    $env:GOOSE_DBSTRING = $previousGooseDbString

    Pop-Location
}