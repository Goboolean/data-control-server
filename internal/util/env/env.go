package env

import (
	"errors"
	"os"
	"path/filepath"

	"github.com/joho/godotenv"
)

// Just import this package to get all the env variables at the root of the project
// Import this package anonymously as shown below:
// import _ "github.com/Goboolean/fetch-server/internal/util/env"

func init() {
	path, err := filepath.Abs(".")
	if err != nil {
		panic(err)
	}

	for base := filepath.Base(path); base != "fetch-server" && base != "app"; {
		path = filepath.Dir(path)
		base = filepath.Base(path)

		if base == "." || base == "/" {
			panic(errRootNotFound)
		}
	}

	if err := os.Chdir(path); err != nil {
		panic(err)
	}

	if err := godotenv.Load(); err != nil {
		panic(err)
	}
}

var errRootNotFound = errors.New("could not find root directory, be sure to set root of the project as fetch-server")