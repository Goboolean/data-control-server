package websocket

import "errors"

var ErrFetcherNotRegistered = errors.New("fetcher is not registered")

var ErrFetcherAlreadyRegistered = errors.New("fetcher is already registered")

var ErrFetcherNotFoundByPlatformName = errors.New("fetcher is not found by platform name")

var ErrPlatformNotFoundByStockId = errors.New("platform is not found by stock id")

var ErrSymbolUnrecognized = errors.New("symbol is unrecognized")

var ErrSymbolNotFoundByStockId = errors.New("symbol is not found by stock id")