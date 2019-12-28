package color

func NewRGBA(r, g, b, a float32) RGBA {
	return RGBA{r, g, b, a}
}

type RGBA struct {
	R, G, B, A float32
}

var zero RGBA

func (r RGBA) IsZero() bool {
	return r == zero
}

var (
	Transparent          = RGBA{0, 0, 0, 0}
	Cornflowerblue       = RGBA{R: 0.39215687, G: 0.58431375, B: 0.92941177, A: 1}
	Lightgray            = RGBA{R: 0.827451, G: 0.827451, B: 0.827451, A: 1}
	Papayawhip           = RGBA{R: 1, G: 0.9372549, B: 0.8352941, A: 1}
	Turquoise            = RGBA{R: 0.2509804, G: 0.8784314, B: 0.8156863, A: 1}
	Chocolate            = RGBA{R: 0.8235294, G: 0.4117647, B: 0.11764706, A: 1}
	Darkorange           = RGBA{R: 1, G: 0.54901963, B: 0, A: 1}
	Darkorchid           = RGBA{R: 0.6, G: 0.19607843, B: 0.8, A: 1}
	Lemonchiffon         = RGBA{R: 1, G: 0.98039216, B: 0.8039216, A: 1}
	Magenta              = RGBA{R: 1, G: 0, B: 1, A: 1}
	Oldlace              = RGBA{R: 0.99215686, G: 0.9607843, B: 0.9019608, A: 1}
	Saddlebrown          = RGBA{R: 0.54509807, G: 0.27058825, B: 0.07450981, A: 1}
	Slategrey            = RGBA{R: 0.4392157, G: 0.5019608, B: 0.5647059, A: 1}
	Darkslateblue        = RGBA{R: 0.28235295, G: 0.23921569, B: 0.54509807, A: 1}
	Hotpink              = RGBA{R: 1, G: 0.4117647, B: 0.7058824, A: 1}
	Orchid               = RGBA{R: 0.85490197, G: 0.4392157, B: 0.8392157, A: 1}
	Seashell             = RGBA{R: 1, G: 0.9607843, B: 0.93333334, A: 1}
	Blueviolet           = RGBA{R: 0.5411765, G: 0.16862746, B: 0.8862745, A: 1}
	Chartreuse           = RGBA{R: 0.49803922, G: 1, B: 0, A: 1}
	Darkgray             = RGBA{R: 0.6627451, G: 0.6627451, B: 0.6627451, A: 1}
	Palegreen            = RGBA{R: 0.59607846, G: 0.9843137, B: 0.59607846, A: 1}
	Plum                 = RGBA{R: 0.8666667, G: 0.627451, B: 0.8666667, A: 1}
	Wheat                = RGBA{R: 0.9607843, G: 0.87058824, B: 0.7019608, A: 1}
	Brown                = RGBA{R: 0.64705884, G: 0.16470589, B: 0.16470589, A: 1}
	Darkgoldenrod        = RGBA{R: 0.72156864, G: 0.5254902, B: 0.043137256, A: 1}
	Olivedrab            = RGBA{R: 0.41960785, G: 0.5568628, B: 0.13725491, A: 1}
	Peachpuff            = RGBA{R: 1, G: 0.85490197, B: 0.7254902, A: 1}
	Darkmagenta          = RGBA{R: 0.54509807, G: 0, B: 0.54509807, A: 1}
	Goldenrod            = RGBA{R: 0.85490197, G: 0.64705884, B: 0.1254902, A: 1}
	Olive                = RGBA{R: 0.5019608, G: 0.5019608, B: 0, A: 1}
	Palegoldenrod        = RGBA{R: 0.93333334, G: 0.9098039, B: 0.6666667, A: 1}
	Purple               = RGBA{R: 0.5019608, G: 0, B: 0.5019608, A: 1}
	Bisque               = RGBA{R: 1, G: 0.89411765, B: 0.76862746, A: 1}
	Fuchsia              = RGBA{R: 1, G: 0, B: 1, A: 1}
	Lightgreen           = RGBA{R: 0.5647059, G: 0.93333334, B: 0.5647059, A: 1}
	Teal                 = RGBA{R: 0, G: 0.5019608, B: 0.5019608, A: 1}
	Beige                = RGBA{R: 0.9607843, G: 0.9607843, B: 0.8627451, A: 1}
	Dimgrey              = RGBA{R: 0.4117647, G: 0.4117647, B: 0.4117647, A: 1}
	Mediumpurple         = RGBA{R: 0.5764706, G: 0.4392157, B: 0.85882354, A: 1}
	Peru                 = RGBA{R: 0.8039216, G: 0.52156866, B: 0.24705882, A: 1}
	Cyan                 = RGBA{R: 0, G: 1, B: 1, A: 1}
	Indigo               = RGBA{R: 0.29411766, G: 0, B: 0.50980395, A: 1}
	Lightseagreen        = RGBA{R: 0.1254902, G: 0.69803923, B: 0.6666667, A: 1}
	Mistyrose            = RGBA{R: 1, G: 0.89411765, B: 0.88235295, A: 1}
	Steelblue            = RGBA{R: 0.27450982, G: 0.50980395, B: 0.7058824, A: 1}
	Darkgrey             = RGBA{R: 0.6627451, G: 0.6627451, B: 0.6627451, A: 1}
	Lightgrey            = RGBA{R: 0.827451, G: 0.827451, B: 0.827451, A: 1}
	Mediumturquoise      = RGBA{R: 0.28235295, G: 0.81960785, B: 0.8, A: 1}
	Rosybrown            = RGBA{R: 0.7372549, G: 0.56078434, B: 0.56078434, A: 1}
	Black                = RGBA{R: 0, G: 0, B: 0, A: 1}
	Darkcyan             = RGBA{R: 0, G: 0.54509807, B: 0.54509807, A: 1}
	Honeydew             = RGBA{R: 0.9411765, G: 1, B: 0.9411765, A: 1}
	Mediumorchid         = RGBA{R: 0.7294118, G: 0.33333334, B: 0.827451, A: 1}
	Orangered            = RGBA{R: 1, G: 0.27058825, B: 0, A: 1}
	Pink                 = RGBA{R: 1, G: 0.7529412, B: 0.79607844, A: 1}
	Powderblue           = RGBA{R: 0.6901961, G: 0.8784314, B: 0.9019608, A: 1}
	Darkolivegreen       = RGBA{R: 0.33333334, G: 0.41960785, B: 0.18431373, A: 1}
	Darkviolet           = RGBA{R: 0.5803922, G: 0, B: 0.827451, A: 1}
	Lawngreen            = RGBA{R: 0.4862745, G: 0.9882353, B: 0, A: 1}
	Lightsteelblue       = RGBA{R: 0.6901961, G: 0.76862746, B: 0.87058824, A: 1}
	Mintcream            = RGBA{R: 0.9607843, G: 1, B: 0.98039216, A: 1}
	Navy                 = RGBA{R: 0, G: 0, B: 0.5019608, A: 1}
	White                = RGBA{R: 1, G: 1, B: 1, A: 1}
	Violet               = RGBA{R: 0.93333334, G: 0.50980395, B: 0.93333334, A: 1}
	Cornsilk             = RGBA{R: 1, G: 0.972549, B: 0.8627451, A: 1}
	Dimgray              = RGBA{R: 0.4117647, G: 0.4117647, B: 0.4117647, A: 1}
	Gray                 = RGBA{R: 0.5019608, G: 0.5019608, B: 0.5019608, A: 1}
	Greenyellow          = RGBA{R: 0.6784314, G: 1, B: 0.18431373, A: 1}
	Lightblue            = RGBA{R: 0.6784314, G: 0.84705883, B: 0.9019608, A: 1}
	Lightgoldenrodyellow = RGBA{R: 0.98039216, G: 0.98039216, B: 0.8235294, A: 1}
	Gold                 = RGBA{R: 1, G: 0.84313726, B: 0, A: 1}
	Sienna               = RGBA{R: 0.627451, G: 0.32156864, B: 0.1764706, A: 1}
	Deeppink             = RGBA{R: 1, G: 0.078431375, B: 0.5764706, A: 1}
	Floralwhite          = RGBA{R: 1, G: 0.98039216, B: 0.9411765, A: 1}
	Slateblue            = RGBA{R: 0.41568628, G: 0.3529412, B: 0.8039216, A: 1}
	Tan                  = RGBA{R: 0.8235294, G: 0.7058824, B: 0.54901963, A: 1}
	Lightslategrey       = RGBA{R: 0.46666667, G: 0.53333336, B: 0.6, A: 1}
	Mediumslateblue      = RGBA{R: 0.48235294, G: 0.40784314, B: 0.93333334, A: 1}
	Mediumspringgreen    = RGBA{R: 0, G: 0.98039216, B: 0.6039216, A: 1}
	Orange               = RGBA{R: 1, G: 0.64705884, B: 0, A: 1}
	Darkgreen            = RGBA{R: 0, G: 0.39215687, B: 0, A: 1}
	Salmon               = RGBA{R: 0.98039216, G: 0.5019608, B: 0.44705883, A: 1}
	Sandybrown           = RGBA{R: 0.95686275, G: 0.6431373, B: 0.3764706, A: 1}
	Springgreen          = RGBA{R: 0, G: 1, B: 0.49803922, A: 1}
	Darkseagreen         = RGBA{R: 0.56078434, G: 0.7372549, B: 0.56078434, A: 1}
	Darkturquoise        = RGBA{R: 0, G: 0.80784315, B: 0.81960785, A: 1}
	Deepskyblue          = RGBA{R: 0, G: 0.7490196, B: 1, A: 1}
	Forestgreen          = RGBA{R: 0.13333334, G: 0.54509807, B: 0.13333334, A: 1}
	Lightpink            = RGBA{R: 1, G: 0.7137255, B: 0.75686276, A: 1}
	Red                  = RGBA{R: 1, G: 0, B: 0, A: 1}
	Azure                = RGBA{R: 0.9411765, G: 1, B: 1, A: 1}
	Lightskyblue         = RGBA{R: 0.5294118, G: 0.80784315, B: 0.98039216, A: 1}
	Linen                = RGBA{R: 0.98039216, G: 0.9411765, B: 0.9019608, A: 1}
	Mediumvioletred      = RGBA{R: 0.78039217, G: 0.08235294, B: 0.52156866, A: 1}
	Navajowhite          = RGBA{R: 1, G: 0.87058824, B: 0.6784314, A: 1}
	Silver               = RGBA{R: 0.7529412, G: 0.7529412, B: 0.7529412, A: 1}
	Tomato               = RGBA{R: 1, G: 0.3882353, B: 0.2784314, A: 1}
	Aliceblue            = RGBA{R: 0.9411765, G: 0.972549, B: 1, A: 1}
	Aquamarine           = RGBA{R: 0.49803922, G: 1, B: 0.83137256, A: 1}
	Darkblue             = RGBA{R: 0, G: 0, B: 0.54509807, A: 1}
	Darkkhaki            = RGBA{R: 0.7411765, G: 0.7176471, B: 0.41960785, A: 1}
	Lightsalmon          = RGBA{R: 1, G: 0.627451, B: 0.47843137, A: 1}
	Lime                 = RGBA{R: 0, G: 1, B: 0, A: 1}
	Firebrick            = RGBA{R: 0.69803923, G: 0.13333334, B: 0.13333334, A: 1}
	Green                = RGBA{R: 0, G: 0.5019608, B: 0, A: 1}
	Yellow               = RGBA{R: 1, G: 1, B: 0, A: 1}
	Aqua                 = RGBA{R: 0, G: 1, B: 1, A: 1}
	Ivory                = RGBA{R: 1, G: 1, B: 0.9411765, A: 1}
	Maroon               = RGBA{R: 0.5019608, G: 0, B: 0, A: 1}
	Mediumblue           = RGBA{R: 0, G: 0, B: 0.8039216, A: 1}
	Grey                 = RGBA{R: 0.5019608, G: 0.5019608, B: 0.5019608, A: 1}
	Limegreen            = RGBA{R: 0.19607843, G: 0.8039216, B: 0.19607843, A: 1}
	Whitesmoke           = RGBA{R: 0.9607843, G: 0.9607843, B: 0.9607843, A: 1}
	Coral                = RGBA{R: 1, G: 0.49803922, B: 0.3137255, A: 1}
	Darkslategrey        = RGBA{R: 0.18431373, G: 0.30980393, B: 0.30980393, A: 1}
	Indianred            = RGBA{R: 0.8039216, G: 0.36078432, B: 0.36078432, A: 1}
	Lavenderblush        = RGBA{R: 1, G: 0.9411765, B: 0.9607843, A: 1}
	Paleturquoise        = RGBA{R: 0.6862745, G: 0.93333334, B: 0.93333334, A: 1}
	Skyblue              = RGBA{R: 0.5294118, G: 0.80784315, B: 0.92156863, A: 1}
	Crimson              = RGBA{R: 0.8627451, G: 0.078431375, B: 0.23529412, A: 1}
	Darksalmon           = RGBA{R: 0.9137255, G: 0.5882353, B: 0.47843137, A: 1}
	Seagreen             = RGBA{R: 0.18039216, G: 0.54509807, B: 0.34117648, A: 1}
	Yellowgreen          = RGBA{R: 0.6039216, G: 0.8039216, B: 0.19607843, A: 1}
	Darkred              = RGBA{R: 0.54509807, G: 0, B: 0, A: 1}
	Dodgerblue           = RGBA{R: 0.11764706, G: 0.5647059, B: 1, A: 1}
	Khaki                = RGBA{R: 0.9411765, G: 0.9019608, B: 0.54901963, A: 1}
	Lightcyan            = RGBA{R: 0.8784314, G: 1, B: 1, A: 1}
	Slategray            = RGBA{R: 0.4392157, G: 0.5019608, B: 0.5647059, A: 1}
	Thistle              = RGBA{R: 0.84705883, G: 0.7490196, B: 0.84705883, A: 1}
	Antiquewhite         = RGBA{R: 0.98039216, G: 0.92156863, B: 0.84313726, A: 1}
	Gainsboro            = RGBA{R: 0.8627451, G: 0.8627451, B: 0.8627451, A: 1}
	Ghostwhite           = RGBA{R: 0.972549, G: 0.972549, B: 1, A: 1}
	Snow                 = RGBA{R: 1, G: 0.98039216, B: 0.98039216, A: 1}
	Mediumaquamarine     = RGBA{R: 0.4, G: 0.8039216, B: 0.6666667, A: 1}
	Palevioletred        = RGBA{R: 0.85882354, G: 0.4392157, B: 0.5764706, A: 1}
	Blue                 = RGBA{R: 0, G: 0, B: 1, A: 1}
	Cadetblue            = RGBA{R: 0.37254903, G: 0.61960787, B: 0.627451, A: 1}
	Lavender             = RGBA{R: 0.9019608, G: 0.9019608, B: 0.98039216, A: 1}
	Lightcoral           = RGBA{R: 0.9411765, G: 0.5019608, B: 0.5019608, A: 1}
	Lightslategray       = RGBA{R: 0.46666667, G: 0.53333336, B: 0.6, A: 1}
	Lightyellow          = RGBA{R: 1, G: 1, B: 0.8784314, A: 1}
	Darkslategray        = RGBA{R: 0.18431373, G: 0.30980393, B: 0.30980393, A: 1}
	Mediumseagreen       = RGBA{R: 0.23529412, G: 0.7019608, B: 0.44313726, A: 1}
	Moccasin             = RGBA{R: 1, G: 0.89411765, B: 0.70980394, A: 1}
	Royalblue            = RGBA{R: 0.25490198, G: 0.4117647, B: 0.88235295, A: 1}
	Blanchedalmond       = RGBA{R: 1, G: 0.92156863, B: 0.8039216, A: 1}
	Burlywood            = RGBA{R: 0.87058824, G: 0.72156864, B: 0.5294118, A: 1}
	Midnightblue         = RGBA{R: 0.09803922, G: 0.09803922, B: 0.4392157, A: 1}
)

var Map = map[string]RGBA{
	"Cornflowerblue":       Cornflowerblue,
	"Lightgray":            Lightgray,
	"Papayawhip":           Papayawhip,
	"Turquoise":            Turquoise,
	"Chocolate":            Chocolate,
	"Darkorange":           Darkorange,
	"Darkorchid":           Darkorchid,
	"Lemonchiffon":         Lemonchiffon,
	"Magenta":              Magenta,
	"Oldlace":              Oldlace,
	"Saddlebrown":          Saddlebrown,
	"Slategrey":            Slategrey,
	"Darkslateblue":        Darkslateblue,
	"Hotpink":              Hotpink,
	"Orchid":               Orchid,
	"Seashell":             Seashell,
	"Blueviolet":           Blueviolet,
	"Chartreuse":           Chartreuse,
	"Darkgray":             Darkgray,
	"Palegreen":            Palegreen,
	"Plum":                 Plum,
	"Wheat":                Wheat,
	"Brown":                Brown,
	"Darkgoldenrod":        Darkgoldenrod,
	"Olivedrab":            Olivedrab,
	"Peachpuff":            Peachpuff,
	"Darkmagenta":          Darkmagenta,
	"Goldenrod":            Goldenrod,
	"Olive":                Olive,
	"Palegoldenrod":        Palegoldenrod,
	"Purple":               Purple,
	"Bisque":               Bisque,
	"Fuchsia":              Fuchsia,
	"Lightgreen":           Lightgreen,
	"Teal":                 Teal,
	"Beige":                Beige,
	"Dimgrey":              Dimgrey,
	"Mediumpurple":         Mediumpurple,
	"Peru":                 Peru,
	"Cyan":                 Cyan,
	"Indigo":               Indigo,
	"Lightseagreen":        Lightseagreen,
	"Mistyrose":            Mistyrose,
	"Steelblue":            Steelblue,
	"Darkgrey":             Darkgrey,
	"Lightgrey":            Lightgrey,
	"Mediumturquoise":      Mediumturquoise,
	"Rosybrown":            Rosybrown,
	"Black":                Black,
	"Darkcyan":             Darkcyan,
	"Honeydew":             Honeydew,
	"Mediumorchid":         Mediumorchid,
	"Orangered":            Orangered,
	"Pink":                 Pink,
	"Powderblue":           Powderblue,
	"Darkolivegreen":       Darkolivegreen,
	"Darkviolet":           Darkviolet,
	"Lawngreen":            Lawngreen,
	"Lightsteelblue":       Lightsteelblue,
	"Mintcream":            Mintcream,
	"Navy":                 Navy,
	"White":                White,
	"Violet":               Violet,
	"Cornsilk":             Cornsilk,
	"Dimgray":              Dimgray,
	"Gray":                 Gray,
	"Greenyellow":          Greenyellow,
	"Lightblue":            Lightblue,
	"Lightgoldenrodyellow": Lightgoldenrodyellow,
	"Gold":                 Gold,
	"Sienna":               Sienna,
	"Deeppink":             Deeppink,
	"Floralwhite":          Floralwhite,
	"Slateblue":            Slateblue,
	"Tan":                  Tan,
	"Lightslategrey":       Lightslategrey,
	"Mediumslateblue":      Mediumslateblue,
	"Mediumspringgreen":    Mediumspringgreen,
	"Orange":               Orange,
	"Darkgreen":            Darkgreen,
	"Salmon":               Salmon,
	"Sandybrown":           Sandybrown,
	"Springgreen":          Springgreen,
	"Darkseagreen":         Darkseagreen,
	"Darkturquoise":        Darkturquoise,
	"Deepskyblue":          Deepskyblue,
	"Forestgreen":          Forestgreen,
	"Lightpink":            Lightpink,
	"Red":                  Red,
	"Azure":                Azure,
	"Lightskyblue":         Lightskyblue,
	"Linen":                Linen,
	"Mediumvioletred":      Mediumvioletred,
	"Navajowhite":          Navajowhite,
	"Silver":               Silver,
	"Tomato":               Tomato,
	"Aliceblue":            Aliceblue,
	"Aquamarine":           Aquamarine,
	"Darkblue":             Darkblue,
	"Darkkhaki":            Darkkhaki,
	"Lightsalmon":          Lightsalmon,
	"Lime":                 Lime,
	"Firebrick":            Firebrick,
	"Green":                Green,
	"Yellow":               Yellow,
	"Aqua":                 Aqua,
	"Ivory":                Ivory,
	"Maroon":               Maroon,
	"Mediumblue":           Mediumblue,
	"Grey":                 Grey,
	"Limegreen":            Limegreen,
	"Whitesmoke":           Whitesmoke,
	"Coral":                Coral,
	"Darkslategrey":        Darkslategrey,
	"Indianred":            Indianred,
	"Lavenderblush":        Lavenderblush,
	"Paleturquoise":        Paleturquoise,
	"Skyblue":              Skyblue,
	"Crimson":              Crimson,
	"Darksalmon":           Darksalmon,
	"Seagreen":             Seagreen,
	"Yellowgreen":          Yellowgreen,
	"Darkred":              Darkred,
	"Dodgerblue":           Dodgerblue,
	"Khaki":                Khaki,
	"Lightcyan":            Lightcyan,
	"Slategray":            Slategray,
	"Thistle":              Thistle,
	"Antiquewhite":         Antiquewhite,
	"Gainsboro":            Gainsboro,
	"Ghostwhite":           Ghostwhite,
	"Snow":                 Snow,
	"Mediumaquamarine":     Mediumaquamarine,
	"Palevioletred":        Palevioletred,
	"Blue":                 Blue,
	"Cadetblue":            Cadetblue,
	"Lavender":             Lavender,
	"Lightcoral":           Lightcoral,
	"Lightslategray":       Lightslategray,
	"Lightyellow":          Lightyellow,
	"Darkslategray":        Darkslategray,
	"Mediumseagreen":       Mediumseagreen,
	"Moccasin":             Moccasin,
	"Royalblue":            Royalblue,
	"Blanchedalmond":       Blanchedalmond,
	"Burlywood":            Burlywood,
	"Midnightblue":         Midnightblue,
}
