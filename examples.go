package main

type Example struct {
	Title       string
	Description string
	Sketch      string
	FullPath    string
	Code        string
	Snapshot    string
}

var Examples = []*Example{
	{
		Title:       "Hello, Triangle",
		Description: `The "hello, world" of graphics: a triangle with different colored vertices.`,
		Sketch:      "001_triangle.go",
	},
	{
		Title:       "Basic Shapes",
		Description: `The "dd" package holds 2D geometry types. The "gfx" package holds operations like "Fill", "Stroke", etc.`,
		Sketch:      "002_shapes.go",
	},
	{
		Title:       "Rotate",
		Description: "Rotate a shape using OpenGL.",
		Sketch:      "004_rotate.go",
	},
	{
		Title:       "Paths",
		Description: "Draw paths using a Pen.",
		Sketch:      "005_path.go",
	},
	{
		Title:       "Blue Noise",
		Description: "Generate evenly spaced random points.",
		Sketch:      "006_blue_noise.go",
	},
	{
		Title:       "Grid",
		Description: "Grids are useful for laying out shapes in a grid pattern.",
		Sketch:      "007_grid.go",
	},
	{
		Title:       "Triangulation",
		Description: "Turn a set of points into triangles.",
		Sketch:      "009_tesselation.go",
	},
	{
		Title:       "Gradient",
		Description: "",
		Sketch:      "010_gradient.go",
	},
	{
		Title:       "Blur",
		Description: "Gaussian blur.",
		Sketch:      "013_blur.go",
	},
	{
		Title:       "Noise",
		Description: "Generate perlin (simplex?) noise using OpenGL.",
		Sketch:      "014_noise.go",
	},
	{
		Title:       "Tweak",
		Description: "Tweak the vertices of a mesh, to give it some character.",
		Sketch:      "015_tweak.go",
	},
	{
		Title:       "Image",
		Description: "Display an image (currently only PNG?)",
		Sketch:      "016_image.go",
	},
	{
		Title:       "Opacity",
		Description: "Testing opacity and blending (a tricky thing to get right, so it's probably wrong...)",
		Sketch:      "017_opacity.go",
	},
	{
		Title:       "Instancing",
		Description: "Generate thousands of instances of the same shape efficiently using OpenGL.",
		Sketch:      "018_instancing.go",
	},
	/*
		{
			Title:       "Window Size",
			Description: "Change the window size from the sketch code.",
			Sketch:      "019_window_size.go",
		},
	*/
	{
		Title:       "Hex Color",
		Description: "Convert hex to a color.",
		Sketch:      "020_color.go",
	},
	{
		Title:       "Stateful Context",
		Description: `Experimenting with adding a stateful "context" drawing API, since people are very familiar with this pattern and it can save some verbose lines of code.`,
		Sketch:      "021_context.go",
	},
	{
		Title:       "Circle Pack",
		Description: "Currently, a failed experiement with finding an efficient circle packing algorithm. Maybe some day I'll crack it.",
		Sketch:      "circle_pack.go",
	},
	{
		Title:       "Combos",
		Description: "Generate all 3x3 combos.",
		Sketch:      "combos.go",
	},
	{
		Title:       "Coqart Grid",
		Description: "Replica of a piece by Roger Coqart.",
		Sketch:      "coqart_grid.go",
	},
	{
		Title:       "Grayblocks",
		Description: "Just for fun.",
		Sketch:      "grayblock.go",
	},
	{
		Title:       "Hex Grid",
		Description: "Generating hexagons in a grid.",
		Sketch:      "hexgrid.go",
	},
	{
		Title:       "Hobbs Split",
		Description: "Inspired by an essay from Tyler Hobbs.",
		Sketch:      "hobbs_split.go",
	},
	{
		Title:       "Molnar2",
		Description: "Inspired by Vera Molnar.",
		Sketch:      "molnar2.go",
	},
	{
		Title:       "Mosaic Sun",
		Description: "Trying to figure out mosaic styling.",
		Sketch:      "mosaic_sun.go",
	},
	{
		Title:       "Nees",
		Description: "Inspired by George Nees.",
		Sketch:      "nees.go",
	},
	{
		Title:       "Noise 1D",
		Description: "1D noise line.",
		Sketch:      "noise_1d.go",
	},
	{
		Title:       "Rotline",
		Description: "",
		Sketch:      "rotline.go",
	},
	{
		Title:       "Sand Spline",
		Description: "Inspired by Anders Hoff and Jared Tarbell. Not working the way I want yet. Splines need love.",
		Sketch:      "sand_spline.go",
	},
	{
		Title:       "Sand Stroke",
		Description: "Inspired by Jared Tarbell.",
		Sketch:      "sand_stroke.go",
	},
	{
		Title:       "Fast Sand Stroke",
		Description: "Fast version of Sand Stroke, because it relies more on OpenGL/GPU.",
		Sketch:      "sand_stroke_fast.go",
	},
	{
		Title:       "Squiggle Grid",
		Description: "Inspired by George Nees.",
		Sketch:      "squiggle_grid.go",
	},
	{
		Title:       "Tribox2",
		Description: "Just for fun.",
		Sketch:      "tribox_2.go",
	},
	{
		Title:       "Tribox3",
		Description: "Just for fun.",
		Sketch:      "tribox_3.go",
	},
	{
		Title:       "Voronoi",
		Description: "Generating a voronoi mesh.",
		Sketch:      "voronoi.go",
	},
}
