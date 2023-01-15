package env

import (
	"fmt"
	"os"
	"strings"
)

type Error struct {
	id     string
	source *Source
	line   int
	col    int
	text   string
}

var (
	errors   []*Error
	messages map[string]string
)

func initErrors() {
	errors = make([]*Error, 0)
	messages = make(map[string]string)

	buf, err := os.ReadFile("errors.txt")
	if err != nil {
		fmt.Println("! error reading errors.txt file ", err.Error())
		return
	}

	var lines = strings.Split(string(buf[:]), "\n")

	for _, s := range lines {
		pair := strings.SplitN(s, ": ", 2)
		if len(pair) == 2 {
			messages[pair[0]] = pair[1]
		}
	}
}

func PosString(pos int) string {
	source, line, col := SourcePos(pos)
	return fmt.Sprintf("%s:%d:%d", source.Path, line, col)
}

func AddError(pos int, id string, args ...interface{}) string {

	source, line, col := SourcePos(pos)

	var err = &Error{
		id:     id,
		source: source,
		line:   line,
		col:    col,
	}

	template, ok := messages[id]
	msg := ""

	if ok {
		msg = fmt.Sprintf(template, args...)
	} else {
		msg = fmt.Sprintf("сообщение для ошибки '%s' не задано!", id)
	}

	err.text = fmt.Sprintf("%s:%d:%d:%s: %s", source.Path, line, col, id, msg)

	errors = append(errors, err)

	return err.text
}

func AddProgramError(id string, args ...interface{}) string {

	var err = &Error{
		id:     id,
		source: nil,
		line:   0,
		col:    0,
	}

	template, ok := messages[id]
	msg := ""

	if ok {
		msg = fmt.Sprintf(template, args...)
	} else {
		msg = fmt.Sprintf("сообщение для ошибки '%s' не задано!", id)
	}

	err.text = fmt.Sprintf("%s: %s", id, msg)

	errors = append(errors, err)

	return err.text
}

func FatalError(id string, args ...interface{}) {
	AddProgramError(id, args...)
	ShowErrors()
	os.Exit(1)
}

func ShowErrors() {
	for _, e := range errors {
		fmt.Println(e.text)
	}
}

func ErrorCount() int {
	return len(errors)
}

func GetError(i int) string {
	return errors[i].text
}

func GetErrorId(i int) string {
	return errors[i].id
}

func ClearErrors() {
	errors = make([]*Error, 0)
}
