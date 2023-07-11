package env

import "errors"

var errRootNotFound = errors.New("could not find root directory, be sure to set root of the project as fetch-server")