package runner

import (
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
)

func copyFile(dstPath, srcPath, name string) error {
	dst, err := os.Create(dstPath)
	if err != nil {
		return err
	}
	defer dst.Close()

	src, err := os.Open(srcPath)
	if err != nil {
		return err
	}

	_, err = dst.Write([]byte("//line " + name + ":1\n"))
	if err != nil {
		return err
	}

	_, err = io.Copy(dst, src)
	if err != nil {
		return err
	}
	defer src.Close()
	return nil
}

// search directory tree for a "go.mod" file for the ink module
// in order to decide if ink should build against a local codebase.
func findInkCode() (string, error) {
	if v := os.Getenv("INK_PATH"); v != "" {
		return v, nil
	}

	const inkMod = "module github.com/buchanae/ink"

	wd, err := os.Getwd()
	if err != nil {
		return "", fmt.Errorf("getting current directory: %w", err)
	}

	current := wd
	for {
		mod := filepath.Join(current, "go.mod")

		info, err := os.Stat(mod)
		if err != nil && !os.IsNotExist(err) {
			return "", fmt.Errorf("getting file info: %w", err)
		}

		if info != nil && info.Mode().IsRegular() {
			b, err := ioutil.ReadFile(mod)
			if err != nil {
				return "", fmt.Errorf("reading file %q: %w", mod, err)
			}
			if strings.HasPrefix(string(b), inkMod) {
				return current, nil
			}
		}

		dir := filepath.Dir(current)
		base := filepath.Base(current)
		if dir == base {
			break
		}
		current = dir
	}
	return "", nil
}
