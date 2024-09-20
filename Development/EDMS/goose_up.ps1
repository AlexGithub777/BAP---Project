# PowerShell script: goose_up.ps1

# Define variables for Goose commands
$GOOSE_CMD = "goose"
$MIGRATION_DIR = "internal/database/migrations"

# Load environment variables from .env file
Get-Content -Path .env | ForEach-Object {
    $key, $value = $_ -split '=', 2
    Set-Item -Path "env:$key" -Value $value
}

# Check if Goose is installed
if (-not (Get-Command $GOOSE_CMD -ErrorAction SilentlyContinue)) {
    Write-Output "Goose is not installed or not in PATH."
    exit 1
}

# Run Goose up command
Write-Output "Running goose up..."
& $GOOSE_CMD -dir $MIGRATION_DIR postgres "user=$env:DB_USER password=$env:DB_PASSWORD dbname=$env:DB_NAME host=$env:DB_HOST port=$env:DB_PORT sslmode=disable" up

Write-Output "Done."