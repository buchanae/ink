package glsl

import (
	"regexp"
)

type Meta struct {
	Vert       []Def
	Frag       []Def
	Uniforms   []Def
	Attributes []Def
}

type Def struct {
	Name     string
	Type     string
	DataType string
}

var rx = regexp.MustCompile(`(?m)^\s*(uniform|in|out) (\w+) (\w+);$`)

func Inspect(vert, frag string) Meta {
	meta := Meta{
		Vert: inspect(vert),
		Frag: inspect(frag),
	}

	uni := map[Def]struct{}{}

	for _, def := range meta.Vert {
		if def.Type == "uniform" {
			if _, seen := uni[def]; !seen {
				meta.Uniforms = append(meta.Uniforms, def)
			}
		}
		if def.Type == "in" {
			meta.Attributes = append(meta.Attributes, def)
		}
	}

	for _, def := range meta.Frag {
		if def.Type == "uniform" {
			if _, seen := uni[def]; !seen {
				meta.Uniforms = append(meta.Uniforms, def)
			}
		}
	}

	return meta
}

func inspect(src string) []Def {
	var defs []Def

	matches := rx.FindAllStringSubmatch(src, -1)
	for _, match := range matches {
		def := Def{match[3], match[1], match[2]}
		defs = append(defs, def)
	}

	return defs
}
