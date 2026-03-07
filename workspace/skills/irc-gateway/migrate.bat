@echo off
REM IRC Gateway Migration Script: v1.0.0 to v1.0.1
REM Windows version

echo =========================================
echo   IRC Gateway Migration: v1.0.0 to v1.0.1
echo =========================================
echo.

REM Check if running from correct directory
if not exist "gateway.py" (
    echo [ERROR] Please run this script from workspace\skills\irc-gateway\
    pause
    exit /b 1
)

echo Step 1: Checking current setup...
echo -----------------------------------

REM Check for .env file
set ENV_FILE=..\..\..\..\env
if exist "%ENV_FILE%" (
    echo [OK] Found .env file in repo root
    set HAS_ENV=1
) else (
    echo [WARN] No .env file found
    set HAS_ENV=0
)

REM Check for PicoClaw config
set PICOCLAW_CONFIG=%USERPROFILE%\.picoclaw\config.json
if exist "%PICOCLAW_CONFIG%" (
    echo [OK] Found PicoClaw config: %PICOCLAW_CONFIG%
    
    REM Check if Telegram is configured
    findstr /C:"telegram" "%PICOCLAW_CONFIG%" >nul 2>&1
    if %ERRORLEVEL% EQU 0 (
        echo [OK] Telegram configured in PicoClaw
        set HAS_PICOCLAW_TOKEN=1
    ) else (
        echo [WARN] Telegram not configured in PicoClaw
        set HAS_PICOCLAW_TOKEN=0
    )
) else (
    echo [WARN] PicoClaw config not found
    set HAS_PICOCLAW_TOKEN=0
)

echo.
echo Step 2: Determining migration path...
echo -----------------------------------

if %HAS_PICOCLAW_TOKEN% EQU 1 (
    echo [OK] Recommended: Use PicoClaw Telegram config
    echo    Gateway will automatically use token from PicoClaw
    echo.
    set /p REMOVE_ENV="Remove .env file to use PicoClaw config? (y/N): "
    
    if /i "%REMOVE_ENV%"=="y" (
        if %HAS_ENV% EQU 1 (
            move "%ENV_FILE%" "%ENV_FILE%.backup" >nul 2>&1
            echo [OK] Backed up .env to .env.backup
        )
        set MIGRATION_PATH=picoclaw
    ) else (
        echo [WARN] Keeping .env file (higher priority than PicoClaw config)
        set MIGRATION_PATH=env
    )
) else if %HAS_ENV% EQU 1 (
    echo [OK] Using existing .env file
    set MIGRATION_PATH=env
) else (
    echo [ERROR] No token source found!
    echo.
    echo Please configure token in one of these locations:
    echo 1. %USERPROFILE%\.picoclaw\config.json (recommended)
    echo 2. .env file in repo root
    echo 3. Environment variable TELEGRAM_BOT_TOKEN
    echo.
    pause
    exit /b 1
)

echo.
echo Step 3: Backing up current gateway.py...
echo -----------------------------------

if exist "gateway.py" (
    copy gateway.py gateway.py.backup >nul 2>&1
    echo [OK] Backed up to gateway.py.backup
)

echo.
echo Step 4: Testing new gateway...
echo -----------------------------------

python -c "import sys; sys.path.insert(0, '.'); import gateway" 2>nul
if %ERRORLEVEL% EQU 0 (
    echo [OK] Gateway imports successfully
) else (
    echo [ERROR] Gateway import failed
    echo    Restoring backup...
    move gateway.py.backup gateway.py >nul 2>&1
    pause
    exit /b 1
)

echo.
echo Step 5: Running tests...
echo -----------------------------------

if exist "test_simple.py" (
    python test_simple.py
    if %ERRORLEVEL% EQU 0 (
        echo [OK] All tests passed
    ) else (
        echo [ERROR] Tests failed
        echo    Restoring backup...
        move gateway.py.backup gateway.py >nul 2>&1
        pause
        exit /b 1
    )
) else (
    echo [WARN] test_simple.py not found, skipping tests
)

echo.
echo =========================================
echo   Migration Complete!
echo =========================================
echo.

if "%MIGRATION_PATH%"=="picoclaw" (
    echo [OK] Gateway will use token from PicoClaw config
    echo    Location: %PICOCLAW_CONFIG%
    echo.
    echo To start gateway:
    echo   python gateway.py
) else (
    echo [OK] Gateway will use token from .env file
    echo    Location: %ENV_FILE%
    echo.
    echo To start gateway:
    echo   python gateway.py
)

echo.
echo Verify token source in startup logs:
echo   [OK] Loaded Telegram token from ^<source^>
echo.
echo Backup files created:
echo   - gateway.py.backup
if exist "%ENV_FILE%.backup" (
    echo   - %ENV_FILE%.backup
)
echo.
echo To rollback:
echo   move gateway.py.backup gateway.py
if exist "%ENV_FILE%.backup" (
    echo   move %ENV_FILE%.backup %ENV_FILE%
)
echo.
echo For more details, see MIGRATION.md
echo.
pause
