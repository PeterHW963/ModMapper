param(
    [ValidateSet("up", "down", "status")]
    [string]$Action = "up"
)

$ErrorActionPreference = "Stop"

$previousGooseDriver = $env:GOOSE_DRIVER
$previousGooseDbString = $env:GOOSE_DBSTRING

Push-Location $PSScriptRoot

try {
    if (-not (Get-Command goose -ErrorAction SilentlyContinue)) {
        throw "Install Goose with: go install github.com/pressly/goose/v3/cmd/goose@latest. Ensure its installation directory is on PATH."
    }

    # Read the root .env as data, never as executable PowerShell.
    $envPath = Join-Path (Split-Path $PSScriptRoot -Parent) ".env"
    $settings = @{}
    try {
        foreach ($line in Get-Content -LiteralPath $envPath -ErrorAction Stop) {
            $line = $line.Trim()
            if (-not $line -or $line.StartsWith("#")) {
                continue
            }

            $parts = $line.Split('=', 2)
            if ($parts.Count -ne 2) {
                throw "Invalid .env variable assignment"
            }

            $key = $parts[0].Trim()
            $val = $parts[1].Trim()

            $settings[$key] = $val
        }
    } catch {
        throw "Could not read .env. Use one KEY=value entry per line."
    }

    $required = @('POSTGRES_HOST', 'POSTGRES_PORT', 'POSTGRES_DB',
                  'POSTGRES_USER', 'POSTGRES_PASSWORD', 'POSTGRES_SSLMODE')
    foreach ($key in $required) {
        if ([string]::IsNullOrEmpty($settings[$key])) {
            throw "Missing required setting in .env: $key"
        }
    }

    $dbHost = $settings.POSTGRES_HOST
    $dbPort = $settings.POSTGRES_PORT
    $dbUser = [Uri]::EscapeDataString($settings.POSTGRES_USER)
    $dbPassword = [Uri]::EscapeDataString($settings.POSTGRES_PASSWORD)
    $dbName = [Uri]::EscapeDataString($settings.POSTGRES_DB)
    $dbSSLMode = [Uri]::EscapeDataString($settings.POSTGRES_SSLMODE)

    $env:GOOSE_DRIVER = "postgres"
    $env:GOOSE_DBSTRING = "postgres://${dbUser}:${dbPassword}@${dbHost}:${dbPort}/${dbName}?sslmode=${dbSSLMode}"

    Write-Host "Running migration action: $Action"
    goose -dir migrations $Action

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
