<a href="https://godoc.org/github.com/buchanae/ink"><img src="https://godoc.org/github.com/buchanae/ink?status.svg" alt="GoDoc"></a>

Ink is a framework for creative 2D graphics in [Go](https://golang.org), based on OpenGL.

# Example: a simple triangle

Write:
```go
// example.go
package main

import (
	. "github.com/buchanae/ink/color"
	. "github.com/buchanae/ink/dd"
	. "github.com/buchanae/ink/gfx"
)

func Ink(doc *Doc) {

	t := Triangle{
		XY{0.2, 0.2},
		XY{0.8, 0.2},
		XY{0.5, 0.8},
	}

	s := NewShader(t)
	s.Set("a_color", []RGBA{
		Red, Green, Blue,
	})
	s.Draw(doc)
}
```

Run:
```
ink example.go
```

Ink opens a window and renders your triangle:


# Status

This is a young project. It's tested only on MacOS 10.14.

# Implementation Notes

Currently, Ink is based on OpenGL only, although other backends are desired. Ink is also focused primarily on 2D graphics so far, because everything is simpler in two dimensions, although I'd like to extend it to 3D some day.
