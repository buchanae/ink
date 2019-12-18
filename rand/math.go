package rand

import "math"

func sqrt(x float32) float32 {
	y := float64(x)
	z := math.Sqrt(y)
	return float32(z)
}

func atan2(a, b float32) float32 {
	c, d := float64(a), float64(b)
	e := math.Atan2(c, d)
	return float32(e)
}

func acos(x float32) float32 {
	y := float64(x)
	z := math.Acos(y)
	return float32(z)
}

func sin(x float32) float32 {
	y := float64(x)
	z := math.Sin(y)
	return float32(z)
}

func cos(x float32) float32 {
	y := float64(x)
	z := math.Cos(y)
	return float32(z)
}
