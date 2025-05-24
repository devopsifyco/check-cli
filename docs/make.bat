@ECHO OFF
REM Minimal make.bat for Sphinx documentation

SET SPHINXBUILD=sphinx-build
SET SOURCEDIR=.
SET BUILDDIR=_build

IF "%1"=="html" (
    %SPHINXBUILD% -b html %SOURCEDIR% %BUILDDIR%/html
    GOTO end
)

IF "%1"=="clean" (
    rmdir /S /Q %BUILDDIR%
    GOTO end
)

:end 