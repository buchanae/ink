package main

import (
	"log"

	"github.com/buchanae/ink/app"
	. "github.com/buchanae/ink/color"
	. "github.com/buchanae/ink/dd"
	"github.com/buchanae/ink/gfx"
	"github.com/buchanae/ink/rand"
)

func Ink(doc *app.Doc) {
	log.SetFlags(0)
	rand.SeedNow()

	var shape Shape = Shapes{
		RectShape{
			RectCenter(Center, XY{.2, .2}),
		},
	}

	for i := 0; i < 150; i++ {
		shape = shape.Next()
	}

	shape.Draw(doc)
}

type Shape interface {
	Next() Shape
	Draw(*app.Doc)
}

type RectShape struct {
	rect Rect
}

func (rs RectShape) Draw(doc *app.Doc) {
	gfx.Fill{
		Mesh:  rs.rect,
		Color: randGray(),
	}.Draw(doc)
}

func (rs RectShape) Next() Shape {

	const (
		Shrink int = iota
		Grow
		Rotate
		Translate
		Duplicate
	)

	actions := []Weighted{
		{150, Shrink},
		{150, Grow},
		{2, Rotate},
		{5, Translate},
		{10, Duplicate},
	}

	switch randWeighted(actions) {
	case Shrink:
		return RectShape{
			rs.rect.ShrinkXY(XY{
				X: rand.Range(0.01, 0.02),
				Y: rand.Range(0.01, 0.02),
			}),
		}
	case Grow:
		return RectShape{
			rs.rect.GrowXY(XY{
				X: rand.Range(0.01, 0.02),
				Y: rand.Range(0.01, 0.02),
			}),
		}
	case Translate:
		return RectShape{
			rs.rect.Translate(XY{
				X: rand.Range(-0.02, 0.02),
				Y: rand.Range(-0.02, 0.02),
			}),
		}

	case Rotate:
		return QuadShape{
			rs.rect.Rotate(rand.Angle()),
		}
	case Duplicate:
		return Shapes{rs, rs}
	}

	return rs
}

type QuadShape struct {
	quad Quad
}

func (qs QuadShape) Draw(doc *app.Doc) {
	gfx.Fill{
		Mesh:  qs.quad,
		Color: randGray(),
	}.Draw(doc)
}

func (qs QuadShape) Next() Shape {
	const (
		Shrink int = iota
		Grow
		Rotate
		Translate
		Duplicate
	)

	actions := []Weighted{
		{2, Rotate},
		{5, Translate},
		{10, Duplicate},
	}

	switch randWeighted(actions) {
	case Translate:
		return QuadShape{
			qs.quad.Translate(XY{
				X: rand.Range(-0.02, 0.02),
				Y: rand.Range(-0.02, 0.02),
			}),
		}
	case Rotate:
		return QuadShape{
			qs.quad.Rotate(rand.Angle()),
		}
	case Duplicate:
		return Shapes{qs, qs}
	}

	return qs
}

type Shapes []Shape

func (shapes Shapes) Draw(doc *app.Doc) {
	for _, shape := range shapes {
		shape.Draw(doc)
	}
}

func (shapes Shapes) Next() Shape {
	next := Shapes{}
	for _, shape := range shapes {
		next = append(next, shape.Next())
	}
	return next
}

type Weighted struct {
	Weight int
	Action int
}

func randWeighted(items []Weighted) int {
	if len(items) == 0 {
		return 0
	}
	if len(items) == 1 {
		return items[0].Action
	}

	sum := 0
	for _, x := range items {
		sum += x.Weight
	}

	v := rand.Intn(sum)
	var item Weighted

	for _, x := range items {
		item = x
		if v < item.Weight {
			break
		}
		v -= item.Weight
	}
	return item.Action
}

func randGray() RGBA {
	g := rand.Range(0, 1)
	return RGBA{g, g, g, .5}
}
