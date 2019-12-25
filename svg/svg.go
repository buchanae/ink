package svg

import "github.com/buchanae/ink/dd"

func Parse(raw string, width, height float32) (dd.Path, error) {
	s := pathParser{Buffer: raw}
	s.builder.width = width
	s.builder.height = height
	s.Init()
	err := s.Parse()
	if err != nil {
		return dd.Path{}, err
	}
	s.Execute()
	return s.path, nil
}
