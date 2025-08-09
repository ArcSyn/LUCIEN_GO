@echo off
echo ========== LUCIEN CLI FINAL PRODUCTION TEST ==========
echo.

echo 1. Basic Commands:
echo pwd | .\build\lucien.exe --batch
echo.

echo 2. Operator Chaining:
echo echo step1 ^&^& echo step2 ^|^| echo backup ; echo final | .\build\lucien.exe --batch
echo.

echo 3. Quoted Operators:
echo echo "command ^&^& operator" | .\build\lucien.exe --batch
echo.

echo 4. Security Toggle:
echo :secure strict | .\build\lucien.exe --batch
echo.

echo 5. Home Command:
echo home | .\build\lucien.exe --batch
echo.

echo 6. Variables:
echo set TESTVAR=production ^& echo $TESTVAR | .\build\lucien.exe --batch
echo.

echo 7. Aliases:
echo alias test='echo alias works' ^& test | .\build\lucien.exe --batch
echo.

echo 8. Error Handling:
echo nonexistentcommand | .\build\lucien.exe --batch
echo.

echo ========== ALL CORE FUNCTIONALITY TESTED ==========