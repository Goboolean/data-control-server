package relayer

import outport "github.com/Goboolean/stock-fetch-server/internal/domain/port/out"

type RelayerManager struct {
	LinedUpRelayer
	RawRelayer
	MiddleRelayer

	ws outport.RelayerPort
}

var instance *RelayerManager

func New() *RelayerManager {
	return instance
}

func init() {
	instance = &RelayerManager{
		LinedUpRelayer: NewLinedUpRelayer(),
		RawRelayer:     NewRawRelayer(),
		MiddleRelayer:  NewMiddleRelayer(),
	}

	go instance.TransferRawToLinedUp()
}
