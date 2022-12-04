package env

import (
	"io/ioutil"
)

type Source struct {
	Path  string
	Bytes []byte
	Lines []int
	Err   error
	No    int // in sources +1
}

var sources []*Source

func initSources() {
	sources = make([]*Source, 0)
}

func AddSource(path string) *Source {

	var src = &Source{
		Path:  path,
		Lines: make([]int, 0),
		No:    len(sources) + 1,
	}

	buf, err := ioutil.ReadFile(path)
	if err != nil {
		src.Err = err
		return src
	}

	src.Bytes = buf

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
	line = 0 // TBD: find line number
	col = ofs
	return
}
