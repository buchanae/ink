package main

import (
	"log"

	"github.com/buchanae/ink/app"
	"github.com/buchanae/ink/color"
	. "github.com/buchanae/ink/color"
	. "github.com/buchanae/ink/dd"
	"github.com/buchanae/ink/gfx"
	"github.com/buchanae/ink/math"
	"github.com/buchanae/ink/rand"
)

const (
	Lines   = 40
	Passes  = 55
	Opacity = 0.02

	JitterY = 0.002

	MinY = 0.005
	MaxY = 0.04

	MinX = -0.03
	MaxX = 0.03

	LineWidth = 0.003
)

func Ink(doc *app.Doc) {
	rand.SeedNow()
	gfx.Clear(doc, White)

	doc.Config.Trace = true

	for k := 0; k < Lines; k++ {
		P := float32(k) / (Lines - 1)
		y := math.Interp(0.05, 0.95, P)

		curves := []Quadratic{}

		pt := XY{0.05, y}
		inc := XY{0.90 / float32(k+1), 0}
		ctrl := XY{}
		var offset float32

		for i := 0; i < k+1; i++ {

			if i%2 == 0 {
				offset = rand.Range(MinX, MaxX)

				ctrl = pt.Add(XY{
					X: inc.X*0.5 + offset,
					Y: rand.Range(MinY, MaxY) * P,
				})
			} else {
				ctrl = pt.Add(XY{
					X: inc.X*0.5 - offset,
					Y: pt.Y - ctrl.Y,
				})
			}

			curves = append(curves, Quadratic{
				A:    pt,
				B:    pt.Add(inc),
				Ctrl: ctrl,
			})

			pt = pt.Add(inc)
		}

		c := Teal
		d := color.RGBA{
			R: c.R*Opacity + (1 - Opacity),
			G: c.G*Opacity + (1 - Opacity),
			B: c.B*Opacity + (1 - Opacity),
			A: 1,
		}
		c.A = 0.3
		c = d
		log.Print(d)
		//c.A = math.Interp(0.05, 0.01, P)
		//c.A = 1

		for j := 0; j < Passes; j++ {

			for i := range curves {
				ctrl := curves[i].Ctrl
				jitter := XY{
					Y: rand.Range(-JitterY, JitterY),
				}
				ctrl = ctrl.Add(jitter)
				curves[i].Ctrl = ctrl
			}

			// TODO super annoying
			path := Path{}
			for _, q := range curves {
				path = append(path, q)
			}

			gfx.Stroke{
				Target: path,
				Color:  c,
				Blend:  gfx.Multiply,
				StrokeOpt: StrokeOpt{
					Width: LineWidth,
				},
			}.Draw(doc)
		}
	}
}
