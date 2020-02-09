<a href="https://godoc.org/github.com/buchanae/ink"><img src="https://godoc.org/github.com/buchanae/ink?status.svg" alt="GoDoc"></a>

Ink is a framework for creative 2D graphics in [Go](https://golang.org), based on OpenGL. Visit [buchanae.github.io/ink](https://buchanae.github.io/ink/) for more.

### Example: a simple triangle
Install:
```
go get github.com/buchanae/ink
```

(Building Ink is a little tricky, because it depends on GLFW. You might need to install these packages:
```
build-essential
xorg-dev 
libgflw3-dev
libxcursor-dev 
libxinerama-dev 
libxi-dev
pkg-config
```

Write `example.go`:
```go
package main

import (
	"github.com/buchanae/ink/app"
	. "github.com/buchanae/ink/color"
	. "github.com/buchanae/ink/dd"
	"github.com/buchanae/ink/gfx"
)

func Ink(doc *app.Doc) {
	t := Triangle{
		XY{0.2, 0.2},
		XY{0.8, 0.2},
		XY{0.5, 0.8},
	}
	s := gfx.Fill{Shape: t}.Shader()
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

![Triangle example](https://buchanae.github.io/ink/assets/snapshots/001_triangle.png)

There are more examples in the [sketches](./sketches) directory.

### Implementation Notes

- Ink is based on OpenGL only, although other backends are desired. There is an experimental WebGL branch.
- Ink is focused on 2D graphics, although 3D is desired some day.
- Most numbers are float32, because OpenGL APIs are mostly based on float32s.
  - This is also why I haven't used the Go stdlib `image/color`. I admit this feels messy.
- Angles are in radians, unless otherwise noted.
