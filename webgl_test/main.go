package main

import (
	"github.com/buchanae/ink/app"
	"github.com/buchanae/ink/color"
	"github.com/buchanae/ink/dd"
	"github.com/buchanae/ink/gfx"
	"github.com/buchanae/ink/render"
)

func main() {
	rend := render.NewRenderer(800, 800)
	doc := app.NewDoc()
	rend.StartTrace()
	doc.Ops = nil

	// The "hello, world" of graphics:
	// a triangle with different colored vertices.
	t := dd.Triangle{
		dd.XY{0.2, 0.2},
		dd.XY{0.8, 0.2},
		dd.XY{0.5, 0.8},
	}

	s := gfx.NewShader(t)
	s.Set("a_color", []color.RGBA{
		color.Red, color.Green, color.Blue,
	})
	s.Draw(doc)

	plan := app.BuildPlan(doc)
	rend.Render(plan)
	rend.ToScreen(doc.LayerID())

	done := make(chan struct{})
	<-done
}
