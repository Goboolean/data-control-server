package relayer

import (
	"context"
	"time"

	"github.com/Goboolean/fetch-server/internal/domain/entity"
	"github.com/Goboolean/fetch-server/internal/domain/service/store"
)


const duration = time.Second / 100


type Mock struct {
	s *store.Store
}



func (*Mock) Subscribe(ctx context.Context, stockId string) (<-chan *entity.StockAggregate, error) {
	
	ch := make(chan *entity.StockAggregate)

	go func(ctx context.Context) {
		for {
			select {
			case <- ctx.Done():
				return
			case <- time.After(duration):
				ch <- &entity.StockAggregate{}
			}
		}
	}(ctx)

	return ch, nil
}


