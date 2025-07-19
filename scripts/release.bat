@echo off
setlocal enabledelayedexpansion

REM Goca Release Script for Windows
REM Usage: scripts\release.bat [version]

echo [INFO] Goca Release Script for Windows
echo.

REM Get version from argument or prompt
if "%1"=="" (
    set /p NEW_VERSION="Enter new version (e.g., 1.0.1): "
) else (
    set NEW_VERSION=%1
)

set TAG_VERSION=v!NEW_VERSION!

echo [INFO] Preparing release for version !TAG_VERSION!

REM Check if we're on the main branch
for /f %%i in ('git branch --show-current') do set CURRENT_BRANCH=%%i
if not "!CURRENT_BRANCH!"=="master" if not "!CURRENT_BRANCH!"=="main" (
    echo [WARNING] You're not on the master/main branch. Current branch: !CURRENT_BRANCH!
    set /p CONTINUE="Do you want to continue? (y/N): "
    if not "!CONTINUE!"=="y" if not "!CONTINUE!"=="Y" (
        echo [ERROR] Release cancelled
        exit /b 1
    )
)

REM Check if tag already exists
git tag -l | findstr /r "^!TAG_VERSION!$" >nul
if !errorlevel! equ 0 (
    echo [ERROR] Tag !TAG_VERSION! already exists
    exit /b 1
)

REM Update version in version.go
echo [INFO] Updating version in cmd\version.go
powershell -Command "(Get-Content cmd\version.go) -replace 'Version.*=.*', 'Version   = \"!NEW_VERSION!\"' | Set-Content cmd\version.go"
for /f "tokens=1-6 delims=/:. " %%a in ("%date% %time%") do (
    set TIMESTAMP=%%c-%%a-%%b
    set TIME_PART=%%d:%%e:%%f
)
powershell -Command "(Get-Content cmd\version.go) -replace 'BuildTime.*=.*', 'BuildTime = \"!TIMESTAMP!T!TIME_PART!Z\"' | Set-Content cmd\version.go"

REM Build and test
echo [INFO] Running tests
go test -v ./...
if !errorlevel! neq 0 (
    echo [ERROR] Tests failed
    exit /b 1
)

echo [INFO] Building application
go build -o goca.exe .
if !errorlevel! neq 0 (
    echo [ERROR] Build failed
    exit /b 1
)

REM Test CLI
echo [INFO] Testing CLI functionality
goca.exe version

REM Git operations
echo [INFO] Committing changes
git add cmd\version.go CHANGELOG.md
git commit -m "chore: bump version to !NEW_VERSION!"
if !errorlevel! neq 0 (
    echo [ERROR] Commit failed
    exit /b 1
)

echo [INFO] Creating tag !TAG_VERSION!
git tag -a "!TAG_VERSION!" -m "Release !TAG_VERSION!"
if !errorlevel! neq 0 (
    echo [ERROR] Tag creation failed
    exit /b 1
)

echo [INFO] Pushing changes and tag
git push origin !CURRENT_BRANCH!
git push origin !TAG_VERSION!
if !errorlevel! neq 0 (
    echo [ERROR] Push failed
    exit /b 1
)

REM Clean up
del goca.exe 2>nul

echo.
echo [SUCCESS] Release !TAG_VERSION! created successfully!
echo [INFO] GitHub Actions will now build and publish the release automatically.
echo [INFO] Check the progress at: https://github.com/sazardev/goca/actions
echo.
echo 🎉 Your release is ready! Here's what happens next:
echo 1. GitHub Actions will build binaries for all platforms
echo 2. A new release will be published automatically  
echo 3. Users can install via: go install github.com/sazardev/goca@!TAG_VERSION!

pause
