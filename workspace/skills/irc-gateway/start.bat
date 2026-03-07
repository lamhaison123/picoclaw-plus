@echo off
REM PicoClaw IRC Gateway - Quick Start Script
REM Double-click to run the gateway

echo ========================================
echo   PicoClaw IRC Gateway Launcher
echo ========================================
echo.

REM Check Python installation
python --version >nul 2>&1
if errorlevel 1 (
    echo [ERROR] Python not found! Please install Python 3.8+
    echo Download from: https://www.python.org/downloads/
    pause
    exit /b 1
)

REM Check for .env file
if not exist "..\..\..\.env" (
    echo [WARNING] .env file not found in repo root
    echo.
    echo Please create .env file with:
    echo   TELEGRAM_BOT_TOKEN=your_bot_token_here
    echo   PICOCLAW_BIN=picoclaw
    echo.
    pause
    exit /b 1
)

REM Install dependencies
echo [1/3] Checking dependencies...
pip install python-telegram-bot --quiet
if errorlevel 1 (
    echo [ERROR] Failed to install dependencies
    pause
    exit /b 1
)

echo [2/3] Starting IRC Gateway...
echo.

REM Run the gateway
python gateway.py

REM If gateway stops
echo.
echo ========================================
echo   Gateway stopped
echo ========================================
pause
