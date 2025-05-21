@echo off
echo Cleaning up temporary directories...

if exist build rmdir /s /q build
if exist dist rmdir /s /q dist
if exist venv rmdir /s /q venv
if exist __pycache__ rmdir /s /q __pycache__

echo Cleanup completed!
echo Remaining files:
dir 