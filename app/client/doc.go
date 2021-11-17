package client

import (
	"image"
	_ "image/png"
	"os"

	"github.com/buchanae/ink/dd"
	"github.com/buchanae/ink/gfx"
	"github.com/buchanae/ink/render"
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

	Conf *gfx.Config
}

func NewDoc() *Doc {
	doc := &Doc{
		ID:   nextID(),
		Conf: &gfx.Config{},
	}
	return doc
}

func (d *Doc) Plan() render.Plan {
	return buildPlan(d)
}

func (d *Doc) Clear() {
	d.Ops = nil
}

func (d *Doc) Config() *gfx.Config {
	return d.Conf
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
	win := d.Conf
	rw := float32(w) / float32(win.Width)
	rh := float32(h) / float32(win.Height)
	rect := dd.RectCenter(
		dd.XY{0.5, 0.5},
		dd.XY{rw, rh},
	)
	return gfx.Image{id, rect}
}

func (d *Doc) LoadImage(path string) gfx.Image {
	fh, err := os.Open(path)
	if err != nil {
		panic(err)
	}
	defer fh.Close()

	img, _, err := image.Decode(fh)
	if err != nil {
		panic(err)
	}

	return d.NewImage(img)
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

func (l *layer) Clear() {
	// TODO wrong. needs fixes to layer ops. see AddShader.
	l.doc.Ops = nil
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
