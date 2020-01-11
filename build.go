package main

import (
	"bytes"
	html "html/template"
	"io/ioutil"
	"log"
	"os"
	"strings"
	"text/template"
)

const AssetRoot = "./assets"

var tpls *template.Template

func init() {

	tpls = template.New("ink")
	tpls.Funcs(map[string]interface{}{
		"Example": Example,
	})

	_, err := tpls.ParseGlob("templates/*.html")
	if err != nil {
		panic(err)
	}
}

func main() {
	log.SetFlags(0)
	out := os.Stdout

	data := Data{
		"AssetRoot": AssetRoot,
		"Content":   toString("index.html", nil),
	}
	err := tpls.ExecuteTemplate(out, "base.html", data)
	if err != nil {
		log.Printf("error: %v", err)
	}
}

func asset(name string) string {
	return AssetRoot + "/" + name
}

func Example(name string) (string, error) {

	path := "examples/" + name + ".go"
	b, err := ioutil.ReadFile(path)
	if err != nil {
		return "", err
	}
	code := string(b)
	code = html.HTMLEscapeString(code)
	code = strings.ReplaceAll(code, "\t", "    ")

	return toString("example.html", Data{
		"Asset": asset(name + ".png"),
		"Code":  code,
	}), nil
}

type Data map[string]interface{}

func toString(name string, data Data) string {
	buf := bytes.NewBuffer(nil)

	err := tpls.ExecuteTemplate(buf, name, data)
	if err != nil {
		log.Printf("error: %v", err)
	}
	return buf.String()
}
