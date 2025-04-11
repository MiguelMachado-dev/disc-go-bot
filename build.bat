@echo off
setlocal enabledelayedexpansion

REM Define the output binary name
set OUTPUT_NAME=disc-go-bot

REM Define the output directory
set DIST_DIR=dist

REM Create dist directory if it doesn't exist
if not exist %DIST_DIR% mkdir %DIST_DIR%

REM Define the platforms and architectures you want to build for
set PLATFORMS=linux/amd64 windows/amd64 darwin/amd64

REM Check if we should build Windows as GUI application
set BUILD_GUI=0
if "%1"=="--gui" set BUILD_GUI=1

REM Check if building only for current platform
set BUILD_CURRENT_ONLY=0
if "%1"=="--current" set BUILD_CURRENT_ONLY=1
if "%2"=="--current" set BUILD_CURRENT_ONLY=1

REM Get current platform
for /f "tokens=*" %%a in ('go env GOOS') do set CURRENT_OS=%%a
for /f "tokens=*" %%a in ('go env GOARCH') do set CURRENT_ARCH=%%a
echo Current platform: %CURRENT_OS%/%CURRENT_ARCH%

if %BUILD_CURRENT_ONLY%==1 (
    echo Building only for current platform: %CURRENT_OS%/%CURRENT_ARCH%
    set PLATFORMS=%CURRENT_OS%/%CURRENT_ARCH%
)

REM Iterate over the platforms and build the binary for each
for %%p in (%PLATFORMS%) do (
    for /f "tokens=1,2 delims=/" %%a in ("%%p") do (
        set GOOS=%%a
        set GOARCH=%%b

        echo Building for %%a/%%b...
        set "output_name=%DIST_DIR%\%OUTPUT_NAME%_%%a_%%b"

        if "%%a"=="windows" (
            set "output_name=!output_name!.exe"

            if %BUILD_GUI%==1 (
                echo Building with GUI flag...
                set CGO_ENABLED=1
                go env -w GOOS=%%a GOARCH=%%b
                go build -ldflags="-H=windowsgui" -o !output_name!
            ) else (
                set CGO_ENABLED=0
                go env -w GOOS=%%a GOARCH=%%b
                go build -o !output_name!
            )

            if !errorlevel! neq 0 (
                echo Error building for %%a/%%b
            ) else (
                echo Built binary: !output_name!
            )
        ) else if "%%a"=="!CURRENT_OS!" (
            REM Building for current OS, CGO should work
            set CGO_ENABLED=1
            go env -w GOOS=%%a GOARCH=%%b
            go build -o !output_name!

            if !errorlevel! neq 0 (
                echo Error building for %%a/%%b
            ) else (
                echo Built binary: !output_name!
            )
        ) else (
            REM Cross-compilation, disable CGO
            set CGO_ENABLED=0
            go env -w GOOS=%%a GOARCH=%%b
            go build -o !output_name!

            if !errorlevel! neq 0 (
                echo Error building for %%a/%%b
            ) else (
                echo Built binary: !output_name!
            )
        )
    )
)

REM Also create a copy of the current platform binary in the dist root with the basic name
if %BUILD_CURRENT_ONLY%==0 (
    echo Creating shortcut binary for current platform...
    if "%CURRENT_OS%"=="windows" (
        copy "%DIST_DIR%\%OUTPUT_NAME%_%CURRENT_OS%_%CURRENT_ARCH%.exe" "%DIST_DIR%\%OUTPUT_NAME%.exe"
    ) else (
        copy "%DIST_DIR%\%OUTPUT_NAME%_%CURRENT_OS%_%CURRENT_ARCH%" "%DIST_DIR%\%OUTPUT_NAME%"
    )
)

REM Reset environment variables back to current OS
go env -u GOOS
go env -u GOARCH

endlocal