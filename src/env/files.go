package env

import (
	e "errors"
	"fmt"
	"os"
	//"path"
	//"strings"
)

var _ = fmt.Printf

func CheckFolder(fpath string) error {
	fi, err := os.Stat(fpath)
	if err != nil {
		return err
	}
	if !fi.IsDir() {
		return e.New("это файл")
	}
	return nil
}
