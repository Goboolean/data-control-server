package env

import (
	"path/filepath"
)

var Root string

func init() {
	path, err := filepath.Abs(".")
	if err != nil {
		panic(err)
	}

	for base := filepath.Base(path); base != "fetch-server"; {
		path = filepath.Dir(path)
		base = filepath.Base(path)

		if base == "." || base == "/" {
			panic(errRootNotFound)
		}
	}

	Root = path
}