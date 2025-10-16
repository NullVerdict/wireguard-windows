@echo off
setlocal enabledelayedexpansion
set BUILDDIR=%~dp0
set PATH=%BUILDDIR%.deps\llvm-mingw\bin;%BUILDDIR%.deps;%PATH%
set PATHEXT=.exe
cd /d %BUILDDIR% || exit /b 1

if exist .deps\prepared goto :render
:installdeps
	rmdir /s /q .deps 2> NUL
	mkdir .deps || goto :error
	cd .deps || goto :error

	rem --- LLVM-MINGW (последний релиз)
	for /f "tokens=2 delims=:," %%v in ('curl -s https://api.github.com/repos/mstorsjo/llvm-mingw/releases/latest ^| findstr /i "tag_name"') do set LLVM_TAG=%%~v
	set LLVM_TAG=%LLVM_TAG:"=%
	call :download llvm-mingw.zip https://github.com/mstorsjo/llvm-mingw/releases/download/%LLVM_TAG%/llvm-mingw-%LLVM_TAG%-msvcrt-x86_64.zip || goto :error

	rem --- ImageMagick (последний стабильный релиз)
	for /f "tokens=2 delims=:," %%v in ('curl -s https://api.github.com/repos/ImageMagick/ImageMagick/releases/latest ^| findstr /i "tag_name"') do set IM_TAG=%%~v
	set IM_TAG=%IM_TAG:"=%
	call :download imagemagick.zip https://imagemagick.org/archive/binaries/ImageMagick-%IM_TAG%-portable-Q16-x64.zip || goto :error

	rem --- GNU Make (последний без guile для Win32)
	for /f "tokens=2 delims=:," %%v in ('curl -s https://api.github.com/repos/ezwinports/make/releases/latest ^| findstr /i "tag_name"') do set MAKE_TAG=%%~v
	set MAKE_TAG=%MAKE_TAG:"=%
	call :download make.zip https://github.com/ezwinports/make/releases/download/%MAKE_TAG%/make-%MAKE_TAG%-without-guile-w32-bin.zip || goto :error

	rem --- wireguard-tools (последний)
	for /f "tokens=2 delims=:," %%v in ('curl -s https://api.github.com/repos/WireGuard/wireguard-tools/commits ^| findstr /i "sha" ^| findstr /v "parent"') do if not defined WG_SHA set WG_SHA=%%~v
	set WG_SHA=%WG_SHA:"=%
	call :download wireguard-tools.zip https://git.zx2c4.com/wireguard-tools/snapshot/wireguard-tools-%WG_SHA%.zip || goto :error

	rem --- wireguard-nt (фиксированная ссылка, релизы закрыты)
	call :download wireguard-nt.zip https://download.wireguard.com/wireguard-nt/wireguard-nt-0.10.1.zip || goto :error

	copy /y NUL prepared > NUL || goto :error
	cd .. || goto :error

:render
	echo [+] Rendering icons
	for %%a in ("ui\icon\*.svg") do convert -background none "%%~fa" -define icon:auto-resize="256,192,128,96,64,48,40,32,24,20,16" -compress zip "%%~dpna.ico" || goto :error

:build
	for /f "tokens=3" %%a in ('findstr /r "Number.*=.*[0-9.]*" .\version\version.go') do set WIREGUARD_VERSION=%%a
	set WIREGUARD_VERSION=%WIREGUARD_VERSION:"=%
	for /f "tokens=1-4" %%a in ("%WIREGUARD_VERSION:.= % 0 0 0") do set WIREGUARD_VERSION_ARRAY=%%a,%%b,%%c,%%d
	set GOOS=windows
	set GOARM=7
	if "%GoGenerate%"=="yes" (
		go generate ./... || exit /b 1
	)
	call :build_plat x86 i686 386 || goto :error
	call :build_plat amd64 x86_64 amd64 || goto :error
	call :build_plat arm64 aarch64 arm64 || goto :error

:sign
	if exist .\sign.bat call .\sign.bat
	if "%SigningProvider%"=="" goto :success
	if "%TimestampServer%"=="" goto :success
	signtool sign %SigningProvider% /fd sha256 /tr "%TimestampServer%" /td sha256 /d WireGuard x86\wireguard.exe x86\wg.exe amd64\wireguard.exe amd64\wg.exe arm64\wireguard.exe arm64\wg.exe || goto :error

:success
	exit /b 0

:download
	echo [+] Downloading %1
	curl -Lf#o %1 %2 || exit /b 1
	tar -xf %1 %~3 || exit /b 1
	del %1 || exit /b 1
	goto :eof

:build_plat
	set GOARCH=%~3
	mkdir %1 >NUL 2>&1
	%~2-w64-mingw32-windres -I ".deps\wireguard-nt\bin\%~1" -DWIREGUARD_VERSION_ARRAY=%WIREGUARD_VERSION_ARRAY% -DWIREGUARD_VERSION_STR=%WIREGUARD_VERSION% -i resources.rc -o "resources_%~3.syso" -O coff -c 65001 || exit /b %errorlevel%
	set CGO_ENABLED=1
	set CC=%~2-w64-mingw32-gcc
	set CFLAGS=-march=core-avx2
	go build -tags load_wgnt_from_rsrc -ldflags="-H windowsgui -s -w" -trimpath -buildvcs=false -v -o "%~1\wireguard.exe" || exit /b 1
	if not exist "%~1\wg.exe" (
		del .deps\src\*.exe .deps\src\*.o .deps\src\wincompat\*.o .deps\src\wincompat\*.lib 2> NUL
		set LDFLAGS=-s -march=core-avx2
		make --no-print-directory -C .deps\src PLATFORM=windows CC=%~2-w64-mingw32-gcc WINDRES=%~2-w64-mingw32-windres V=1 RUNSTATEDIR= SYSTEMDUNITDIR= -j%NUMBER_OF_PROCESSORS% || exit /b 1
		move /Y .deps\src\wg.exe "%~1\wg.exe" > NUL || exit /b 1
	)
	goto :eof

:error
	echo [-] Failed with error #%errorlevel%.
	cmd /c exit %errorlevel%
