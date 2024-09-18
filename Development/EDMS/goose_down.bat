@echo off
setlocal

rem Define variables for Goose commands
set GOOSE_CMD=goose
set MIGRATION_DIR=internal/database/migrations

rem Check if Goose is installed
where %GOOSE_CMD% >nul 2>&1
if errorlevel 1 (
    echo Goose is not installed or not in PATH.
    exit /b 1
)

rem Run Goose down command
echo Running goose down...
%GOOSE_CMD% -dir %MIGRATION_DIR% postgres "user=postgres password=postgres dbname=postgres host=localhost port=5432 sslmode=disable" down

echo Done.
endlocal