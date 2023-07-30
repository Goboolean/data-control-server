package out

import (
	"context"
)



type RelayerPort interface {
	FetchStock(ctx context.Context, stockId string, platform string, symbol string) error
	StopFetchingStock(ctx context.Context, stockId string) error
}