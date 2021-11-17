package dd

var Center = XY{.5, .5}

type Fillable interface {
	Fill() Mesh
}

type Strokeable interface {
	Stroke(StrokeOpt) Mesh
}

// TODO http://www.realtimerendering.com/intersections.html
