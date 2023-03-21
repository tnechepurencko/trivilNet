package env

import (
	er "errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

var _ = fmt.Printf

const (
	file_extension = ".tri"
)

type Source struct {
	FolderPath string
	FolderName string
	FileName   string

	Path string // путь, используемый для чтения

	Bytes []byte
	Lines []int
	Err   error
	No    int // in sources +1
}

var sources []*Source

func initSources() {
	sources = make([]*Source, 0)
}

// Если параметр - папка, делаеть GetFolderSource
// Иначе, читает файл, заданный параметром, и добавляет остальные файлы из папки с нужным расширением
func GetSources(spath string) []*Source {

	if CheckFolder(spath) == nil {
		return GetFolderSources(spath)
	}

	folder, filename := filepath.Split(spath)
	if folder == "" {
		folder = "."
	}
	//fmt.Printf("GetSources: '%v' '%v'\n", folder, filename)

	var src = &Source{
		FolderPath: folder,
		FolderName: filepath.Base(folder),
		FileName:   filename,
		Path:       spath,

		Lines: make([]int, 0),
	}

	var list = make([]*Source, 1)
	list[0] = src

	if !strings.HasSuffix(filename, file_extension) {
		src.Err = er.New("неверное расширение файла " + filepath.Ext(filename))
		return list
	}

	buf, err := os.ReadFile(src.Path)
	if err != nil {
		src.Err = err
		return list
	}
	src.Bytes = buf

	f, err := os.Open(folder)
	if err != nil {
		panic(fmt.Sprintf("panic GetSources(%s): %s", folder, err.Error()))
	}

	names, err := f.Readdirnames(0)
	for _, name := range names {
		if name != filename && strings.HasSuffix(name, file_extension) {
			var src = readSource(folder, name)
			if src.Err != nil {
				list = make([]*Source, 1)
				list[0] = src
				return list
			}
			list = append(list, src)
		}
	}

	for _, s := range list {
		sources = append(sources, s)
		s.No = len(sources)
	}

	return list
}

func NormalizeFolderPath(fpath string) string {
	// обработать репу - возможно, добавить имя реры, если оно не указано
	return fpath
}

// Выдает список прочитанных исходников из папки
// Если исходник не удается прочитать, то на этом составление списка заканчивается
// и возвращается список из одного непрочитанного исходника.
// В случае успеха, все исходники добавляются в список всех исходников
func GetFolderSources(folder string) []*Source {

	var list = make([]*Source, 0)

	f, err := os.Open(folder)
	if err != nil {
		panic(fmt.Sprintf("panic GetFolderSources(%s): %s", folder, err.Error()))
	}

	names, err := f.Readdirnames(0)
	for _, name := range names {
		if strings.HasSuffix(name, file_extension) {
			var src = readSource(folder, name)
			if src.Err != nil {
				list = make([]*Source, 1)
				list[0] = src
				return list
			}
			list = append(list, src)
		}
	}

	for _, s := range list {
		sources = append(sources, s)
		s.No = len(sources)
	}

	return list
}

func readSource(folder, filename string) *Source {

	var src = &Source{
		FolderPath: folder,
		FolderName: filepath.Base(folder),
		FileName:   filename,
		Path:       filepath.Join(folder, filename),

		Lines: make([]int, 0),
	}

	buf, err := os.ReadFile(src.Path)
	if err != nil {
		src.Err = err
		return src
	}

	src.Bytes = buf

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

//====

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
