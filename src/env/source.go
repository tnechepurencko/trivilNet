package env

import (
	"fmt"
	"io/ioutil"
	"path"
	"strings"
)

var _ = fmt.Printf

const (
	file_extension = ".tri"
)

type Source struct {
	Original string // исходный путь
	Path     string // путь, используемый для чтения, может быть с расширением файла
	LastName string // последнее имя в пути
	Bytes    []byte
	Lines    []int
	Err      error
	No       int // in sources +1
}

var sources []*Source

func initSources() {
	sources = make([]*Source, 0)
}

//TODO: каноническое имя для сравнения - сейчас используется Path
func AddSource(spath string) *Source {

	var src = &Source{
		Original: spath,
		Path:     spath,
		Lines:    make([]int, 0),
		No:       len(sources) + 1,
	}

	if path.Ext(spath) == "" {
		src.Path += file_extension
	}

	_, filename := path.Split(src.Path)

	var ext = path.Ext(filename)
	if ext != "" {
		filename = strings.TrimSuffix(filename, ext)
	}
	src.LastName = filename

	buf, err := ioutil.ReadFile(src.Path)
	if err != nil {
		src.Err = err
		return src
	}

	src.Bytes = buf

	sources = append(sources, src)

	return src
}

func AddImmSource(text string) *Source {

	var src = &Source{
		Path:  "imm",
		Lines: make([]int, 0),
		No:    len(sources) + 1,
		Bytes: []byte(text),
	}

	sources = append(sources, src)

	return src
}

func (s *Source) AddLine(ofs int) {
	s.Lines = append(s.Lines, ofs)
}

func (s *Source) MakePos(ofs int) int {
	return ofs<<16 + s.No
}

func SourcePos(pos int) (src *Source, line int, col int) {
	no := pos & 0xFFFF
	ofs := pos >> 16

	if no == 0 || no > len(sources) {
		panic("! wrong source index in pos")
	}

	src = sources[no-1]

	line, col = calcTextPos(src, ofs)

	//	line = 0 // TBD: find line number
	//	col = ofs
	return
}

func calcTextPos(src *Source, ofs int) (int, int) {

	if len(src.Lines) == 0 {
		return 0, ofs
	}

	//fmt.Printf("%d in %v\n", ofs, src.Lines)

	var l = 0
	var r = len(src.Lines) - 1

	for {
		if l >= r {
			break
		}
		var x = (l + r) / 2
		var lofs = src.Lines[x]

		//fmt.Printf("%d.%d %d.%d => %d.%d\n", l, src.Lines[l], r, src.Lines[r], x, src.Lines[x])

		if ofs > lofs {
			l = x + 1
		} else if ofs < lofs {
			r = x - 1
		} else {
			l = x
			break
		}
	}

	if ofs < src.Lines[l] && l > 0 {
		l--
	}

	return l + 1, ofs - src.Lines[l]
}
