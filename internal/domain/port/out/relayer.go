package out

import (
	"context"

	"github.com/Goboolean/fetch-server/internal/domain/entity"
)



type RelayerPort interface {
	FetchStock(ctx context.Context, stockId string, stockMeta entity.StockAggsMeta) error
	StopFetchingStock(ctx context.Context, stockId string) error
}