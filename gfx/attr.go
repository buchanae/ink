package gfx

import "github.com/buchanae/ink/color"

type Attrs struct {
	Data map[string]interface{}
}

func (a *Attrs) SetColor(c color.RGBA) {
	a.Set("a_color", c)
}

func (a *Attrs) Set(name string, val interface{}) {
	if a.Data == nil {
		a.Data = map[string]interface{}{}
	}
	a.Data[name] = val
}
