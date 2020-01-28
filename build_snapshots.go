package main

import (
	"bytes"
	"context"
	"image"
	"image/png"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/alecthomas/chroma"
	"github.com/alecthomas/chroma/formatters/html"
	"github.com/alecthomas/chroma/lexers"
	"github.com/buchanae/ink/app"
)

var formatter chroma.Formatter
var lexer chroma.Lexer
var chromaCSS string
var chromaStyle = chroma.MustNewStyle("ink", chroma.StyleEntries{
	chroma.Comment:               "italic",
	chroma.CommentPreproc:        "noitalic",
	chroma.Keyword:               "bold",
	chroma.KeywordPseudo:         "nobold",
	chroma.KeywordType:           "nobold",
	chroma.OperatorWord:          "bold",
	chroma.NameClass:             "bold",
	chroma.NameNamespace:         "bold",
	chroma.NameException:         "bold",
	chroma.NameEntity:            "bold",
	chroma.NameTag:               "bold",
	chroma.LiteralString:         "italic",
	chroma.LiteralStringInterpol: "bold",
	chroma.LiteralStringEscape:   "bold",
	chroma.GenericHeading:        "bold",
	chroma.GenericSubheading:     "bold",
	chroma.GenericEmph:           "italic",
	chroma.GenericStrong:         "bold",
	chroma.GenericPrompt:         "bold",
	chroma.Error:                 "border:#FF0000",
	chroma.Background:            " bg:#f7f7f7",
})

func init() {
	lexer = lexers.Get("go")
	f := html.New(
		html.WithClasses(true),
		html.TabWidth(4),
	)
	formatter = f

	var buf bytes.Buffer
	err := f.WriteCSS(&buf, chromaStyle)
	if err != nil {
		panic(err)
	}
	chromaCSS = buf.String()
}

func buildExampleData() {

	const dir = "/Users/abuchanan/projects/art-apps/ink/sketches"
	for _, ex := range Examples {
		name := strings.Replace(ex.Sketch, ".go", ".png", 1)
		ex.Snapshot = "assets/snapshots/" + name

		ex.FullPath = filepath.Join(dir, ex.Sketch)
		ex.Code = readCode(ex.FullPath)
	}
}

func buildSnapshots() {

	if os.Getenv("INK_PATH") == "" {
		panic("INK_PATH isn't set")
	}

	conf := app.DefaultConfig()
	conf.Window.Hidden = true
	conf.Snapshot.Width = 400
	conf.Snapshot.Height = 400

	a, err := app.NewApp(conf)
	if err != nil {
		panic(err)
	}

	ctx := context.Background()

	go func() {
		for _, ex := range Examples {

			log.Println()
			log.Print(ex.Sketch)

			err := a.RunSketch(ctx, ex.FullPath)
			if err != nil {
				log.Printf("error: %v", err)
				continue
			}

			img := a.Snapshot()

			err = writeImage(img, ex.Snapshot)
			if err != nil {
				log.Printf("error: %v", err)
				continue
			}
		}
		a.Close()
	}()

	err = a.Run()
	if err != nil {
		log.Printf("error: %v", err)
	}
}

func writeImage(img image.Image, dst string) error {
	f, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer f.Close()
	return png.Encode(f, img)
}

func readCode(path string) string {
	b, err := ioutil.ReadFile(path)
	if err != nil {
		panic(err)
	}
	code := string(b)

	iterator, err := lexer.Tokenise(nil, code)
	if err != nil {
		panic(err)
	}

	var buf bytes.Buffer
	err = formatter.Format(&buf, chromaStyle, iterator)
	if err != nil {
		panic(err)
	}
	code = buf.String()

	//code = html.HTMLEscapeString(code)
	//code = strings.ReplaceAll(code, "\t", "    ")
	return code
}
