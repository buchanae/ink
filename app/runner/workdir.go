package runner

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"time"
)

type workdir struct {
	path string
}

func newWorkdir() (wd workdir, err error) {

	wd.path, err = ioutil.TempDir("", "ink-run-")
	if err != nil {
		return
	}

	mainPath := filepath.Join(wd.path, "ink_main_wrapper.go")
	err = ioutil.WriteFile(mainPath, []byte(head), 0644)
	if err != nil {
		return
	}

	// TODO should look for an existing go.mod file in the sketch directory
	modPath := filepath.Join(wd.path, "go.mod")
	modName := "ink_sketch_" + fmt.Sprint(time.Now().Unix())
	modContent := "module " + modName + "\n"

	inkCode, err := findInkCode()
	if err != nil {
		log.Printf("error: finding ink code: %v", err)
	}
	if inkCode != "" {
		modContent += "\n\nreplace github.com/buchanae/ink => " + inkCode
	}

	err = ioutil.WriteFile(modPath, []byte(modContent), 0644)
	if err != nil {
		return
	}

	return
}

func (w workdir) cleanup() {
	if w.path != "" {
		os.RemoveAll(w.path)
	}
}
