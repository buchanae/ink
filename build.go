package main

import (
	"log"
	"os"
	"text/template"
)

const AssetRoot = "./assets"
const CacheBuster = "?2"

var tpls *template.Template

func init() {

	tpls = template.New("ink")

	_, err := tpls.ParseGlob("templates/*.html")
	if err != nil {
		panic(err)
	}
}

func main() {
	log.SetFlags(0)

	buildExampleData()
	buildSnapshots()

	out, err := os.Create("index.html")
	if err != nil {
		panic(err)
	}
	defer out.Close()

	err = tpls.ExecuteTemplate(out, "base.html", Data{
		"Style":     asset("style.css"),
		"ChromaCSS": chromaCSS,
		"Examples":  Examples,
	})
	if err != nil {
		log.Printf("error: %v", err)
	}
}

func asset(name string) string {
	return AssetRoot + "/" + name + CacheBuster
}

type Data map[string]interface{}
