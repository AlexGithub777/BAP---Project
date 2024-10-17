@echo off
setlocal

set GO111MODULE=on
set GOFLAGS=-mod=mod

:: Run go mod download
echo Running go mod download...
go mod download
if %ERRORLEVEL% neq 0 (
    echo Error: Failed to run go mod download.
    exit /b %ERRORLEVEL%
)

:: Build the application from the cmd directory
echo Building Go application...
cd cmd
go build -ldflags="-s -w" -o ../edms.exe .
if %ERRORLEVEL% neq 0 (
    echo Error: Failed to build Go application.
    exit /b %ERRORLEVEL%
)

echo Build successful!
cd ..
endlocal