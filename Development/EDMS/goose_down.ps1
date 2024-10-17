# PowerShell script: goose_down.ps1

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

# Run Goose down command
Write-Output "Running goose down..."
& $GOOSE_CMD -dir $MIGRATION_DIR postgres "user=$env:DB_USER password=$env:DB_PASSWORD dbname=$env:DB_NAME host=$env:DB_HOST port=$env:DB_PORT sslmode=disable" down

# Check if Goose down command was successful
if ($LASTEXITCODE -eq 0) {
    Write-Output "Goose down command executed successfully."

    # Delete file named "seedcomplete" from "./internal" directory
    $filePath = Join-Path -Path "./internal" -ChildPath "seed_complete"
    if (Test-Path $filePath) {
        Remove-Item -Path $filePath -Force
        Write-Output "File 'seed_complete' deleted successfully."
    } else {
        Write-Output "File 'seed_complete' does not exist."
    }
} else {
    Write-Output "Goose down command failed."
}

Write-Output "Done."