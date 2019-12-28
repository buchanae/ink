package gfx

var currentID int

func nextID() int {
	currentID++
	return currentID
}

type Doc struct {
	OnFrame func(Frame)

	nodes []Node
}

type Node struct {
	LayerID int
	Op      interface{}
}

func (d *Doc) NewLayer() Layer {
	return newLayer(d)
}

func (d *Doc) LayerID() int {
	return 0
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

func (l *layer) AddShader(s *Shader) {
	// layer writes nodes to the root doc
	// in order to maintain a flat list of nodes.
	l.doc.nodes = append(l.doc.nodes, Node{
		LayerID: l.id,
		Op:      s,
	})
}
