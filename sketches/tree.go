package main

import (
	"log"
	"math"

	"github.com/buchanae/ink/app"
	. "github.com/buchanae/ink/dd"
	"github.com/buchanae/ink/gfx"
	"github.com/buchanae/ink/rand"
)

const (
	Years = 20
)

var Up = XY{0, 1}
var Right = XY{1, 0}
var Left = XY{-1, 0}
var Origin = XY{}

func Ink(doc *app.Doc) {
	log.SetFlags(0)
	rand.SeedNow()
	palette := rand.Palette()

	trunk := &Branch{
		Position:  XY{.5, .05},
		Direction: Up,
	}

	for year := 0; year < Years; year++ {
		Grow(trunk)
	}

	//Dump(trunk, "")

	for _, b := range Flatten(trunk) {
		gfx.Fill{
			Mesh:  b.Stroke(),
			Color: rand.Color(palette),
		}.Draw(doc)
	}
}

func Grow(branch *Branch) {
	branch.Age++

	dir := branch.Direction
	pos := branch.Position

	for i, node := range branch.Nodes {
		node.Width += 0.001
		node.Position = pos
		pos = pos.Add(node.Vec())
		dir = node.Direction

		depth := float32(i) / float32(len(branch.Nodes))

		if depth > 0.5 && rand.Bool(0.20) && len(node.Branches) == 0 {
			b := node.NewBranch()
			b.Depth = branch.Depth + 1
			if rand.Bool(0.5) {
				b.Direction = node.Direction.Rotate(math.Pi/2, Origin)
			} else {
				b.Direction = node.Direction.Rotate(-math.Pi/2, Origin)
			}
		}

		for _, branch := range node.Branches {
			Grow(branch)
		}
	}

	//dir = dir.Rotate(rand.Range(-0.1, 0.1), Origin)

	branch.Nodes = append(branch.Nodes, &Node{
		Position:  pos,
		Direction: dir,
		Length:    0.03,
		Width:     0.003,
	})
}

type Node struct {
	Position  XY
	Direction XY
	Length    float32
	Width     float32
	Branches  []*Branch
}

func (node *Node) NewBranch() *Branch {
	b := &Branch{
		Position: node.Position,
	}
	node.Branches = append(node.Branches, b)
	return b
}

func (node *Node) Vec() XY {
	return node.Direction.MulScalar(node.Length)
}

func (node *Node) Line() Line {
	return Line{
		node.Position, node.Position.Add(node.Vec()),
	}
}

type Branch struct {
	Position  XY
	Direction XY
	Age       int
	Depth     int
	Nodes     []*Node
}

func (branch *Branch) Stroke() Mesh {
	if len(branch.Nodes) == 0 {
		return Mesh{}
	}

	if len(branch.Nodes) == 1 {
		node := branch.Nodes[0]
		line := node.Line()
		normal := line.Normal()
		width := node.Width / 2
		lnorm := normal.SetLength(width)
		rnorm := normal.SetLength(-width)
		return Quad{
			line.A.Add(lnorm),
			line.A.Add(rnorm),
			line.B.Add(lnorm),
			line.B.Add(rnorm),
		}.Mesh()
	}

	var tris []Triangle
	var prev [2]XY

	for i, node := range branch.Nodes {

		line := node.Line()
		normal := line.Normal()
		width := node.Width / 2
		lnorm := normal.SetLength(width)
		rnorm := normal.SetLength(-width)

		var next Line

		// if on the last line, cap the end points
		if i == len(branch.Nodes)-1 {
			tris = append(tris,
				Triangle{
					prev[0],
					prev[1],
					line.B.Add(rnorm),
				},
				Triangle{
					prev[1],
					line.B.Add(lnorm),
					line.B.Add(rnorm),
				},
			)
			break
		} else {
			next = branch.Nodes[i+1].Line()
		}

		// if on the first line, initialize the start points
		if i == 0 {
			prev[0] = line.A.Add(rnorm)
			prev[1] = line.A.Add(lnorm)
		}

		mp := miterPoint(line, next, width)
		mp2 := miterPoint(line, next, -width)

		tris = append(tris,
			Triangle{
				prev[0],
				mp2,
				prev[1],
			},
			Triangle{
				prev[1],
				mp2,
				mp,
			},
		)

		prev[0] = mp2
		prev[1] = mp
	}

	mesh := Mesh{}
	return mesh.AddTriangles(tris)
}

func miterPoint(a, b Line, width float32) XY {
	n := a.Normal()
	miter := n.Add(b.Normal()).Normalize()
	miterWidth := width / miter.Dot(n)
	return a.B.Add(miter.SetLength(miterWidth))
}

func Dump(branch *Branch, indent string) {
	log.Printf("%sBranch %d", indent, len(branch.Nodes))
	// TODO easy to accidentally write to stdout
	//      which is required by Ink
	indent = indent + "    "

	for _, node := range branch.Nodes {
		log.Printf("%sNode %v", indent, node.Position)
		for _, sub := range node.Branches {
			Dump(sub, indent)
		}
	}
}

func Flatten(branch *Branch) []*Branch {
	out := []*Branch{branch}
	for _, node := range branch.Nodes {
		for _, sub := range node.Branches {
			out = append(out, Flatten(sub)...)
		}
	}
	return out
}
