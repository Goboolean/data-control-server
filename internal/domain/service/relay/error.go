package relay

import "errors"


var ErrStockNotExists = errors.New("stock does not exist")

var ErrStockNotRelayable = errors.New("stock is not relayable")