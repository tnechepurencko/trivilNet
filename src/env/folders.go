package env

import (
	"fmt"
	"os"
	"path/filepath"
	//	"strings"
)

var _ = fmt.Printf

// Папка с настройками компилятора
func SettingsFolder() string {
	return filepath.Dir(os.Args[0])
}

func SettingsRelativePath(filename string) string {
	return filepath.Join(SettingsFolder(), filename)
}
