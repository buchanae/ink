package gfx

import (
	"image"
	_ "image/png"
	"os"

	"github.com/buchanae/ink/color"
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

type Config struct {
	Snapshot SnapshotConfig
	Window   WindowConfig
	Trace    bool
}

type WindowConfig struct {
	Title         string
	Width, Height int
	X, Y          int
	Hidden        bool
}

type SnapshotConfig struct {
	Width, Height int
	Dir           string
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
	}
	// TODO this shouldn't be here
	Clear(doc, color.White)
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

func (d *Doc) NewLayer() Layer {
	return &layer{
		id:  nextID(),
		doc: d,
	}
}

func (d *Doc) NewImage(img image.Image) Image {
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
	return Image{id, rect}
}

func (d *Doc) LoadImage(path string) Image {
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

func (d *Doc) AddShader(s *Shader) {
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
