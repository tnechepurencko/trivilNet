package env

import (
	er "errors"
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
		return er.New("это файл")
	}
	return nil
}
