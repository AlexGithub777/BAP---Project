@echo off
setlocal

set GO111MODULE=on
set GOFLAGS=-mod=vendor

:: Run go mod vendor
echo Running go mod vendor...
go mod vendor
if %ERRORLEVEL% neq 0 (
    echo Error: Failed to run go mod vendor.
    exit /b %ERRORLEVEL%
)

:: Build the application from the cmd directory
echo Building Go application...
cd cmd
<<<<<<< HEAD
go build -ldflags="-s -w" -o ../edms.exe .
=======
go build -ldflags="-s -w" -o ../safety-device-app.exe .
>>>>>>> d3f1aef86552b4414e16af4df61e0e15859fe0b5
if %ERRORLEVEL% neq 0 (
    echo Error: Failed to build Go application.
    exit /b %ERRORLEVEL%
)

echo Build successful!
cd ..
endlocal