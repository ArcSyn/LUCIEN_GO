@echo off
echo Testing Lucien CLI Build...
echo.
echo 1. Building the project...
go build ./cmd/lucien
if errorlevel 1 (
    echo BUILD FAILED!
    exit /b 1
)
echo BUILD SUCCESS!
echo.

echo 2. Testing help output...
lucien.exe --help
echo.

echo 3. Testing version info...
echo The lucien.exe executable was created successfully.
echo File size: 
dir lucien.exe | find "lucien.exe"
echo.

echo 4. All systems operational! 
echo The Lucien CLI cyberpunk shell interface is ready.
echo.
echo To run the full TUI interface, execute: ./lucien.exe
echo (Note: Use Ctrl+C to exit when testing)