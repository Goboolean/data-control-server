package broker

import "errors"


var errTimeoutRequired = errors.New("timeout setting on ctx required")

var errFatalWhileDeletingTopic = errors.New("fatal error while deleting topic")

var errTrivalWhileDeletingTopic = errors.New("trival error while deleting topic")