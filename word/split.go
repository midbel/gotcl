package word

import (
	"strings"
)

func Split(str string) ([]Word, error) {
	s, err := Scan(strings.NewReader(str))
	if err != nil {
		return nil, err
	}
	var list []Word
	for {
		w := s.Split()
		if w.Type == EOF {
			break
		}
		list = append(list, w)
	}
	return list, nil
}
