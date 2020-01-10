package math

import "math"

const (
	Pi = math.Pi
)

func Sqrt(x float32) float32 {
	y := float64(x)
	z := math.Sqrt(y)
	return float32(z)
}

func Sin(x float32) float32 {
	y := float64(x)
	z := math.Sin(y)
	return float32(z)
}

func Cos(x float32) float32 {
	y := float64(x)
	z := math.Cos(y)
	return float32(z)
}

func Acos(x float32) float32 {
	y := float64(x)
	z := math.Acos(y)
	return float32(z)
}

func Atan2(a, b float32) float32 {
	c, d := float64(a), float64(b)
	e := math.Atan2(c, d)
	return float32(e)
}

func Clamp(val, min, max float32) float32 {
	if val > max {
		return max
	}
	if val < min {
		return min
	}
	return val
}

func Interp(from, to, percent float32) float32 {
	return from + ((to - from) * percent)
}
