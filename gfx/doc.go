package gfx

import "image"

var currentID int

func nextID() int {
	currentID++
	return currentID
}

type Doc struct {
	OnFrame func(Frame)

	id    int
	nodes []Node
}

func NewDoc() *Doc {
	return &Doc{id: nextID()}
}

type Node struct {
	LayerID int
	Op      interface{}
}

func (d *Doc) Nodes() []Node {
	return d.nodes
}

func (d *Doc) NewLayer() Layer {
	return newLayer(d)
}

func (d *Doc) LayerID() int {
	return d.id
}

func (d *Doc) Clear() {
	d.nodes = nil
}

func (d *Doc) NewImage(img image.Image) Layer {
	l := newLayer(d)
	d.nodes = append(d.nodes, Node{
		LayerID: l.id,
		Op:      img,
	})
	return l
}

func (d *Doc) AddShader(s *Shader) {
	d.nodes = append(d.nodes, Node{
		Op: s,
	})
}

func newLayer(doc *Doc) *layer {
	return &layer{
		id:  nextID(),
		doc: doc,
	}
}

type layer struct {
	id  int
	doc *Doc
}

func (l *layer) NewLayer() Layer {
	return newLayer(l.doc)
}

func (l *layer) LayerID() int {
	return l.id
}

func (l *layer) NewImage(img image.Image) Layer {
	return l.doc.NewImage(img)
}

func (l *layer) AddShader(s *Shader) {
	// layer writes nodes to the root doc
	// in order to maintain a flat list of nodes.
	l.doc.nodes = append(l.doc.nodes, Node{
		LayerID: l.id,
		Op:      s,
	})
}
