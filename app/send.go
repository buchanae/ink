// +build sendonly

package app

import (
	"encoding/gob"
	"log"
	"os"

	"github.com/buchanae/ink/gfx"
)

// Render renders the doc with default config and blocks until the window is closed.
// If there is no window open, Render will open one.
// If an error occurs while opening the window, Render panics.
// If an error occurs while rendering the doc, Render logs the error.
func Render(doc *gfx.Doc) {
	err := gob.NewEncoder(os.Stdout).Encode(doc)
	if err != nil {
		log.Println(err)
	}
}
