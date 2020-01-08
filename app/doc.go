// Package app contains the high-level Ink application.
package app

import (
	"image"

	"github.com/buchanae/ink/gfx"
)

type Op struct {
	LayerID int
	Shader  *gfx.Shader
}

var currentID int

func nextID() int {
	currentID++
	return currentID
}

type Doc struct {
	ID     int
	Images map[int]image.Image
	Ops    []Op
}

func NewDoc() *Doc {
	return &Doc{ID: nextID()}
}

func (d *Doc) Filter(layerID ...int) *Doc {
	out := NewDoc()
	out.Images = d.Images
	for _, op := range d.Ops {
		for _, id := range layerID {
			if op.LayerID == id {
				out.Ops = append(out.Ops, op)
				break
			}
		}
	}
	return out
}

func (d *Doc) LayerID() int {
	return d.ID
}

func (d *Doc) NewLayer() gfx.Layer {
	return &layer{
		id:  nextID(),
		doc: d,
	}
}

func (d *Doc) NewImage(img image.Image) gfx.Image {
	if d.Images == nil {
		d.Images = map[int]image.Image{}
	}
	id := nextID()
	d.Images[id] = img
	return gfx.Image{id}
}

func (d *Doc) AddShader(s *gfx.Shader) {
	d.Ops = append(d.Ops, Op{
		LayerID: d.ID,
		Shader:  s,
	})
}

type layer struct {
	id  int
	doc *Doc
}

func (l *layer) LayerID() int {
	return l.id
}

func (l *layer) NewLayer() gfx.Layer {
	return l.doc.NewLayer()
}

func (l *layer) NewImage(img image.Image) gfx.Image {
	return l.doc.NewImage(img)
}

func (l *layer) AddShader(s *gfx.Shader) {
	// layer writes nodes to the root doc
	// in order to maintain a flat list of nodes.
	l.doc.Ops = append(l.doc.Ops, Op{
		LayerID: l.id,
		Shader:  s,
	})
}
