package relayer

import "github.com/Goboolean/stock-fetch-server/internal/domain/value"

type MiddleRelayer struct {
	queue map[string][]value.StockAggregate
}

func NewMiddleRelayer() MiddleRelayer {
	return MiddleRelayer{
		queue: make(map[string][]value.StockAggregate),
	}
}

func (r *MiddleRelayer) isExist(stock string) bool {
	_, ok := r.queue[stock]
	return ok
}

func (r *MiddleRelayer) push(stock string, data *value.StockAggregate) {

	if r.isExist(stock) == false {
		r.queue[stock] = make([]value.StockAggregate, 0)
	}

	r.queue[stock] = append(r.queue[stock], *data)
}

func (r *MiddleRelayer) emptyQueue(stock string) ([]value.StockAggregate, bool) {
	if len(r.queue[stock]) > 100 {
		batch := r.queue[stock]

		delete(r.queue, stock)

		return batch, true
	}
	return nil, false
}
