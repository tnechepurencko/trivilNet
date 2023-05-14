## Компилятор Trivil-0

### Что читать

* Как внести свой вклад в проект: doc/CONTRIBUTING.md
* Описание языка: doc/report/report.pdf
* Тривиль проект является частью проекта "Языки выходного дня" (ЯВД или WEL), введение в проект: http://алексейнедоря.рф/?p=419
* ЯВД является частью проекта "Интенсивное программирование", введение в проект: http://digital-economy.ru/stati/интенсивное-программирование

### Установка и использование компилятора Тривиль на Windows

1) Перед установкой Тривиля должны быть установлены:
* Go: https://go.dev/doc/install
* Clang

2) Запустить установочный скрипт

```
cd scripts
install <path to install directory>
```
3) Добавить install directory в PATH

4) Написать модуль на Тривиле. Для создания/ редактирования исходных текстов рекомендуем использовать Notepad++ к которому есть подсветка синтаксиса: doc/tri-notepad.xml

5) Перейти в рабочую папку (в ней будет создан exe), запустить:

```
tric <путь к папке головного модуля>
```
Например:
```
tric проекты/мой 
tric . // если головной модуль лежит в рабочей папке
```

### Тестирование пакетов компилятора

```
cd src
go test ./...
```

### Размер компилятора:

Строчек кода на Go:

* 2022.12.16   3135 in 31 files
* 2022.12.22   4450 in 37 files
* 2022.12.31   5611 in 41 files
* 2023.01.08   6583 in 44 files, runtime:   509 lines (C)
* 2023.01.13   7199 in 45 files, runtime:   568 lines (C)
* 2023.01.22   8061 in 50 files, runtime:   685 lines (C)
* 2023.02.03   9001 in 51 files, runtime:   767 lines (C)
* 2023.03.11   9417 in 54 files, runtime:   852 lines (C)
* 2023.04.08   9836 in 54 files, runtime:   930 lines (C), libraries:   5,   561 lines (Тривиль)
* 2023.05.14 10369 in 56 files, runtime: 1312 lines (C), libraries: 10, 1230 lines (Тривиль), трик 987 lines (Тривиль)

#### Как посчитать в Windows (PowerShell)

```
#Count lines in Powershell:
(dir -Include *.go -Recurse | select-string "$").Count
#Count files:
(dir -Include *.go -Recurse ).Count
```

#### Как посчитать на Linux

```
cd src
find . -name '*.go' | xargs wc -l
find . -name '*.go' | wc -l

cd ../runtime
find . -name '*.?' | xarg wc -l
find . -name '*.go' | wc -l
```


### Использование Unicode
Концептуально логика работы с Unicode взята из библиотеки языка Julia (минимум).

UTF-8: https://github.com/JuliaStrings/utf8proc

### Русификация консоли

#### Windows

См: https://remontka.pro/fix-cyrillic-windows-10/

Панель управления -> Региональные стандарты -> Дополнительное -> Изменить язык системы:

Включить галочку: Использовать UTF-8

#### Linux

В подавляющем большинстве случаев на Linux установлена локаль `ru_RU.UTF-8`, этого достаточно для работы
с Trivil.

#### Русский язык в именах файлов в Git

Однако, при работе с Git имена файлов отображаются в виде escape-последовательностей, что не очень удобно.

Для изменения поведения Git нужно выполнить команду

```bash
git config --global core.quotePath false
```
