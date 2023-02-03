Trivil-0 compiler

test all compiler packages:
cd src; go test ./...

Compiler lines (Go):
2022.12.16 3135 in 31 files
2022.12.22 4450 in 37 files
2022.12.31 5611 in 41 files
2023.01.08 6583 in 44 files, runtime: 509 lines (C)
2023.01.13 7199 in 45 files, runtime: 568 lines (C)
2023.01.22 8061 in 50 files, runtime: 685 lines (C)
2023.02.03 9001 in 51 files, runtime: 767 lines (C)


#Count lines in Powershell:
(dir -Include *.go -Recurse | select-string "$").Count
#Count files:
(dir -Include *.go -Recurse ).Count


UTF-8: https://github.com/JuliaStrings/utf8proc


Руссификация консоли, см: https://remontka.pro/fix-cyrillic-windows-10/
Панель управления -> Региональные стандарты -> Дополнительное -> Изменить язык системы:
  включить галочку: Использовать UTF-8

