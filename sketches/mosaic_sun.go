package main

import (
	"github.com/buchanae/ink/app"
	"github.com/buchanae/ink/color"
	. "github.com/buchanae/ink/dd"
	"github.com/buchanae/ink/gfx"
	"github.com/buchanae/ink/math"
	"github.com/buchanae/ink/rand"
)

const (
	Gap   = 0.001
	Space = 0.0045
	Count = 25

	Start = 0.01
	Width = (0.5 - Start) / Count

	MinChord = 0.005
	MaxChord = 0.008

	JumpChance    = 0.2
	LightenChance = 0.4
	LightenAmt    = 0.2
	TweakChance   = 0.3
	TweakAmt      = -0.001
)

func Ink(doc *app.Doc) {
	rand.SeedNow()

	bg := color.HexString("#260d05")
	bg = color.White
	gfx.Clear(doc, bg)

	A := color.HexString("#ebc334")
	B := color.HexString("#0c79e8")

	//palette := rand.Palette()
	//A = rand.Color(palette)
	//B = rand.Color(palette)

	gfx.Fill{
		Mesh: Circle{
			XY:       XY{.5, .5},
			Radius:   Start - Gap,
			Segments: 5,
		},
		Color: A,
	}.Draw(doc)

	for i := float32(0); i < Count; i++ {

		rings := Rings{
			Offset: rand.Range(0, 3),
			Inner:  Start + i*Width,
			Outer:  Start + (i+1)*Width - Space,
			Gap:    Gap,
		}

		var min float32 = MinChord
		var max float32 = MaxChord

		min += i * 0.001
		max += i * 0.003
		max = math.Min(max, 0.05)

		chords := GenChords(rings.Inner, min, max)
		// TODO interpcolor isn't based on visual interpolation
		//      going form orange to blue goes through green
		col := InterpColor(A, B, i/Count)

		for _, in := range chords {
			rx := rings
			rx.From = in.From
			rx.To = in.To
			rx.Gap += rand.Range(0, 0.002)
			rx.Color = col

			if rand.Bool(JumpChance) {
				rx.Color = InterpColor(A, B,
					rand.Range(0, Count)/Count,
				)
			}

			if rand.Bool(LightenChance) {
				rx.Color = Lighten(
					rx.Color, rand.Range(-LightenAmt, LightenAmt),
				)
			}
			rx.Draw(doc)
		}
	}
}

func Lighten(c color.RGBA, amt float32) color.RGBA {
	h, s, v := RGBToHSV(c.R, c.G, c.B)
	// TODO broken?
	v += amt
	r, g, b := HSVToRGB(h, s, v)
	return color.RGBA{r, g, b, c.A}
}

func HSVToRGB(h, s, v float32) (r, g, b float32) {
	f := func(n float32) float32 {
		k := math.Mod(n+h/60, 6)
		m := math.Min(math.Min(k, 4-k), 1)
		if m < 0 {
			m = 0
		}
		return v - (v * s * m)
	}

	return f(5), f(3), f(1)
}

func RGBToHSV(r, g, b float32) (h, s, v float32) {
	// https://en.wikipedia.org/wiki/HSL_and_HSV#From_RGB

	min := math.Min(math.Min(r, g), b)
	max := math.Max(math.Max(r, g), b)

	switch {
	case max == min:
	case max == r:
		h = 60 * (0 + ((g - b) / (max - min)))
	case max == g:
		h = 60 * (2 + ((g - b) / (max - min)))
	case max == b:
		h = 60 * (4 + ((g - b) / (max - min)))
	}

	if h < 0 {
		h += 360
	}

	if max != 0 {
		s = (max - min) / max
	}

	v = max
	return
}

func InterpColor(from, to color.RGBA, p float32) color.RGBA {
	return color.RGBA{
		R: math.Interp(from.R, to.R, p),
		G: math.Interp(from.G, to.G, p),
		B: math.Interp(from.B, to.B, p),
		A: math.Interp(from.A, to.A, p),
	}
}

type Rings struct {
	Inner, Outer float32
	From, To     float32
	Offset       float32
	Gap          float32
	Color        color.RGBA
}

func (r Rings) Draw(doc gfx.Layer) {
	center := XY{.5, .5}
	inner := Circle{
		XY:     center,
		Radius: r.Inner,
	}
	outer := Circle{
		XY:     center,
		Radius: r.Outer,
	}
	from := r.Offset + r.From
	to := r.Offset + r.To
	innerGap := ChordAngle(r.Inner, r.Gap)
	outerGap := ChordAngle(r.Outer, r.Gap)
	quad := Quad{
		inner.XYFromAngle(from + innerGap),
		inner.XYFromAngle(to - innerGap),
		outer.XYFromAngle(to - outerGap),
		outer.XYFromAngle(from + outerGap),
	}
	if rand.Bool(TweakChance) {
		quad = rand.TweakQuad(quad, TweakAmt)
	}
	fill := gfx.Fill{quad, r.Color}
	fill.Draw(doc)
}

type Chord struct {
	From, To float32
}

func GenChords(radius, min, max float32) []Chord {

	var out []Chord

	p := float32(0)
	for {
		length := rand.Range(min, max)
		// Protect against NaN
		// from Asin in ChordAngle
		if length >= radius {
			length = radius
		}
		ang := ChordAngle(radius, length)
		next := p + ang

		if next >= math.Pi*2 {
			out = append(out, Chord{
				From: p,
			})
			break
		}

		out = append(out, Chord{
			From: p,
			To:   next,
		})
		p = next
	}

	return out
}

func ChordAngle(radius, length float32) float32 {
	return 2 * math.Asin(length/(2*radius))
}

func CircleCircumference(radius float32) float32 {
	return 2 * math.Pi * radius
}
