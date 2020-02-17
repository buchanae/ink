// Package gfx contains types for drawing graphics to an Ink application.
package gfx

import (
	"image"
	_ "image/png"
	"os"

	"github.com/buchanae/ink/dd"
)

type Op struct {
	LayerID int
	Shader  *Shader
}

var currentID int

func nextID() int {
	currentID++
	return currentID
}

type Doc struct {
	ID            int
	Images        map[int]image.Image
	Ops           []Op
	Width, Height int
}

func NewDoc() *Doc {
	doc := &Doc{
		ID: nextID(),
	}
	// TODO this shouldn't be here
	//Clear(doc, color.White)
	return doc
}

func (doc *Doc) LayerID() int {
	return doc.ID
}

func (doc *Doc) NewLayer() Layer {
	return &layer{
		id:  nextID(),
		doc: doc,
	}
}

func (doc *Doc) NewImage(img image.Image) Image {
	if doc.Images == nil {
		doc.Images = map[int]image.Image{}
	}
	id := nextID()
	doc.Images[id] = img

	w, h := img.Bounds().Dx(), img.Bounds().Dy()
	rw := float32(w) / float32(doc.Width)
	rh := float32(h) / float32(doc.Height)
	rect := dd.RectCenter(
		dd.XY{0.5, 0.5},
		dd.XY{rw, rh},
	)
	return Image{id, rect}
}

func (doc *Doc) LoadImage(path string) Image {
	fh, err := os.Open(path)
	if err != nil {
		panic(err)
	}
	defer fh.Close()

	img, _, err := image.Decode(fh)
	if err != nil {
		panic(err)
	}

	return doc.NewImage(img)
}

func (doc *Doc) AddShader(s *Shader) {
	doc.Ops = append(doc.Ops, Op{
		LayerID: doc.ID,
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

func (l *layer) NewLayer() Layer {
	return l.doc.NewLayer()
}

func (l *layer) NewImage(img image.Image) Image {
	return l.doc.NewImage(img)
}

func (l *layer) AddShader(s *Shader) {
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
