@echo off

if "x%1" == "x" echo Usage: install "target directory for installation"
if "x%1" == "x" exit
echo installing to %1 ...

cd ..

del /s /q *.bak

xcopy /q /y /e /k /i стд %1\стд
xcopy /q /y /e /k /i runtime %1\runtime
xcopy /q /y /k /i doc\report\*.pdf %1\doc
xcopy /q /y /k /i doc\*.xml %1\doc
xcopy /q /y /k /i config\*.* %1\config

cd src

go build trivil.go
if exist trivil.exe (
	move trivil.exe %1\tric.exe
) else (
 	echo compilation failed
	exit
)

echo success
