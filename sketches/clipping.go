package main

import (
	"github.com/buchanae/ink/clip"
	. "github.com/buchanae/ink/color"
	. "github.com/buchanae/ink/dd"
	"github.com/buchanae/ink/gfx"
	"github.com/buchanae/ink/tess"
)

func Ink(doc gfx.Doc) {
	rect := RectCenter(Center, XY{.5, .5})

	circle := Circle{
		XY:       XY{.25, .35},
		Radius:   .25,
		Segments: 50,
	}

	/*
		gfx.Fill{
			Shape: rect,
			Color: Red,
		}.Draw(doc)

		gfx.Fill{
			Shape: circle,
			Color: Blue,
		}.Draw(doc)
	*/

	//points := Union(circle.Path(), rect.Path())
	//points := Intersection(circle.Path(), rect.Path())
	//points := Difference(circle.Path(), rect.Path())
	points := clip.Difference(rect.Path(), circle.Path())

	circleped := tess.Tesselate(points)
	gfx.Fill{
		Shape: Triangles(circleped),
		Color: Yellow,
	}.Draw(doc)

	//gfx.Dots{points, Black, .005}.Draw(doc)

}

func flipColor(a bool, b, c RGBA) RGBA {
	if a {
		return b
	}
	return c
}
