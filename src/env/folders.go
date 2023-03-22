package env

import (
	goerr "errors"
	"fmt"
	"os"

	"path"
	"path/filepath"
	"strings"
)

var _ = fmt.Printf

const (
	separator = "::"
	stdName   = "стд"
)

// Обработка входного пути для импорта:
// Поиск по кодовой базе или нормализация
type ImportLookup struct {
	Root       string // "", если не задана кодовая база
	Err        error  // ошибка при обработке кодовой базы или нормализации
	ImportPath string // Найденный или нормализованный путь
}

// Объект для поиска
var Lookup ImportLookup

// Словарь кодовых баз
var sourceRoots = make(map[string]string)

//====

func (il *ImportLookup) Process(fpath string) {

	il.Root = ""
	il.Err = nil

	var parts = strings.SplitN(fpath, separator, 2)
	if len(parts) != 2 || len(parts) >= 1 && parts[0] == "" {

		il.ImportPath, il.Err = filepath.Abs(fpath)

		return
	}

	il.Root = parts[0]
	rootPath, ok := sourceRoots[il.Root]
	if !ok {
		il.Err = goerr.New(fmt.Sprintf("путь для кодовой базы '%s' не задан", il.Root))
		return
	}
	il.ImportPath = path.Join(rootPath, parts[1])
}

//====

var baseFolder = ""

// Возвращает ошибку, если путь не указывает на папку
func EnsureFolder(fpath string) error {
	fi, err := os.Stat(fpath)
	if err != nil {
		return err
	}
	if !fi.IsDir() {
		return goerr.New("это файл")
	}
	return nil
}

// Папка с настройками компилятора
func SettingsFolder() string {
	return baseFolder
}

func SettingsRelativePath(filename string) string {
	return path.Join(baseFolder, filename)
}

func initFolders() {
	baseFolder = filepath.ToSlash(filepath.Dir(os.Args[0]))

	var stdPath = path.Join(baseFolder, stdName)

	var err = EnsureFolder(stdPath)
	if err == nil {
		sourceRoots[stdName] = stdPath
	}

	// TODO: прочитать файл со списком кодовых баз
}
