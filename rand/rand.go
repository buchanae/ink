package rand

import (
	"math/rand"

	"github.com/buchanae/ink/color"
	"github.com/buchanae/ink/dd"
	"github.com/buchanae/ink/math"
)

func New(seed int64) *Rand {
	return &Rand{rand.New(rand.NewSource(seed))}
}

type Rand struct {
	src *rand.Rand
}

func (r *Rand) Intn(max int) int {
	return r.src.Intn(max)
}

func (r *Rand) IntRange(min, max int) int {
	return r.src.Intn(max-min) + min
}
func (r *Rand) Float() float32 {
	return r.src.Float32()
}

func (r *Rand) Bool(chance float32) bool {
	return r.Float() < chance
}

func (r *Rand) Angle() float32 {
	return r.Float() * math.Pi * 2
}

func (r *Rand) Range(min, max float32) float32 {
	return (max-min)*r.Float() + min
}

func (r *Rand) XY() dd.XY {
	return dd.XY{r.Float(), r.Float()}
}

func (r *Rand) XYRange(min, max float32) dd.XY {
	return dd.XY{
		X: r.Range(min, max),
		Y: r.Range(min, max),
	}
}

func (r *Rand) XYInTriangle(t dd.Triangle) dd.XY {
	// https://stackoverflow.com/questions/19654251/random-point-inside-triangle-inside-java#
	r1 := math.Sqrt(r.Float())
	r2 := r.Float()
	t1 := 1 - r1
	t2 := r1 * (1 - r2)
	t3 := r1 * r2
	x := t1*t.A.X + t2*t.B.X + t3*t.C.X
	y := t1*t.A.Y + t2*t.B.Y + t3*t.C.Y
	return dd.XY{x, y}
}

var src = New(1)

func Intn(max int) int {
	return src.Intn(max)
}

func IntRange(min, max int) int {
	return src.IntRange(min, max)
}

func Float() float32 {
	return src.Float()
}

func Bool(chance float32) bool {
	return src.Bool(chance)
}

func Angle() float32 {
	return src.Angle()
}

func Range(min, max float32) float32 {
	return src.Range(min, max)
}

func XY() dd.XY {
	return src.XY()
}

func XYRange(min, max float32) dd.XY {
	return src.XYRange(min, max)
}

func XYInTriangle(t dd.Triangle) dd.XY {
	return src.XYInTriangle(t)
}

func BlueNoise(n int, w, h, d float32) []dd.XY {
	return src.BlueNoise(n, w, h, d)
}

func Color(c []color.RGBA) color.RGBA {
	return src.Color(c)
}

func Palette() []color.RGBA {
	return src.Palette()
}
