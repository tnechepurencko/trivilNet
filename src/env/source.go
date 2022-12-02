package env

import (
	"io/ioutil"
)

type Source struct {
	Path  string
	Bytes []byte
	Lines []int
	Err   error
}

func ReadSource(path string) *Source {

	var src = &Source{
		Path:  path,
		Lines: make([]int, 0),
	}

	buf, err := ioutil.ReadFile(path)
	if err != nil {
		src.Err = err
		return src
	}

	src.Bytes = buf

	return src
}

func (s *Source) AddLine(ofs int) {
	s.Lines = append(s.Lines, ofs)
}
