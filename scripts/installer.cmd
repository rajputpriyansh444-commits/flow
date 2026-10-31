@echo off
setlocal

set "PS1_URL=https://raw.githubusercontent.com/rajputpriyansh444-commits/flow/main/scripts/installer.ps1"
set "TEMP_PS1=%TEMP%\flow-installer-%RANDOM%.ps1"

REM Internal source override for deterministic tests.
if defined _FLOW_INSTALLER_PS1_PATH (
  copy /y "%_FLOW_INSTALLER_PS1_PATH%" "%TEMP_PS1%" >nul
  if errorlevel 1 (
    echo Failed to stage the PowerShell installer. >&2
    exit /b 1
  )
) else (
  powershell -NoProfile -ExecutionPolicy Bypass -Command ^
    "$ProgressPreference='SilentlyContinue'; Invoke-WebRequest -UseBasicParsing -Uri '%PS1_URL%' -OutFile '%TEMP_PS1%';"
  if errorlevel 1 (
    echo Failed to download the PowerShell installer. >&2
    exit /b 1
  )
)

powershell -NoProfile -ExecutionPolicy Bypass -Command ^
  "$content=[IO.File]::ReadAllText('%TEMP_PS1%', [Text.Encoding]::UTF8); [IO.File]::WriteAllText('%TEMP_PS1%', $content, [Text.Encoding]::UTF8);"
if errorlevel 1 (
  echo Failed to prepare the PowerShell installer. >&2
  exit /b 1
)

powershell -NoProfile -ExecutionPolicy Bypass -File "%TEMP_PS1%" %*
set "EXIT_CODE=%ERRORLEVEL%"

del /q "%TEMP_PS1%" >nul 2>&1
exit /b %EXIT_CODE%
