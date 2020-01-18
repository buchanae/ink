// Package app contains the high-level Ink application.
package app

import (
	"image"
	_ "image/png"

	"github.com/buchanae/ink/color"
	"github.com/buchanae/ink/dd"
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
	Config Config
}

func NewDoc() *Doc {
	doc := &Doc{
		ID: nextID(),
		// TODO need to pull this from app
		Config: DefaultConfig(),
	}
	gfx.Clear(doc, color.Black)
	return doc
}

func (d *Doc) Filter(layerID ...int) *Doc {
	out := &Doc{
		ID:     nextID(),
		Images: d.Images,
		Config: d.Config,
	}
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

	w, h := img.Bounds().Dx(), img.Bounds().Dy()
	win := d.Config.Window
	rw := float32(w) / float32(win.Width)
	rh := float32(h) / float32(win.Height)
	rect := dd.RectCenter(
		dd.XY{0.5, 0.5},
		dd.XY{rw, rh},
	)
	return gfx.Image{id, rect}
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
	// TODO this is pretty confusing in some cases.
	//      the name "Layer" gives the impression of layering
	//      the entire groups in layer order, not interleaved.
	//      consider changing the name, or running shaders
	//      in true layer order.
	//
	// layer writes nodes to the root doc
	// in order to maintain a flat list of nodes.
	l.doc.Ops = append(l.doc.Ops, Op{
		LayerID: l.id,
		Shader:  s,
	})
}
