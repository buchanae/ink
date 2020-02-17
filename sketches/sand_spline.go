package main

import (
	. "github.com/buchanae/ink/color"
	. "github.com/buchanae/ink/dd"
	"github.com/buchanae/ink/gfx"
	"github.com/buchanae/ink/rand"
)

const (
	Lines  = 17
	N      = 30
	Passes = 20
	MinY   = 0.005
	MaxY   = 0.04

	MinX = -0.02
	MaxX = 0.02

	LineWidth = 0.003
)

func Ink(doc gfx.Doc) {
	rand.SeedNow()

	redDot := gfx.Dot{Color: Red, Radius: 0.003}

	for k := 0; k < Lines; k++ {
		y := 0.1 + float32(k)*0.05

		for j := 0; j < Passes; j++ {
			var curves Path

			pt := XY{0.05, y}
			inc := XY{0.90 / N, 0}
			var ctrl XY
			var offset float32

			for i := 0; i < N; i++ {

				if i%2 == 0 {
					offset = rand.Range(MinX, MaxX)

					ctrl = pt.Add(XY{
						X: inc.X*0.5 + offset,
						Y: rand.Range(MinY, MaxY) * (float32(Lines-k-1) / Lines),
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

				rd := redDot
				rd.XY = ctrl
				//rd.Draw(doc)
				/*
					TODO want. go proposal
					redDot{
						XY: ctrl,
					}.Draw(doc)
				*/

				pt = pt.Add(inc)
			}

			c := Teal
			c.A = 0.3

			gfx.Stroke{
				Shape: curves,
				Width: LineWidth,
				Color: c,
			}.Draw(doc)
		}
	}
}
