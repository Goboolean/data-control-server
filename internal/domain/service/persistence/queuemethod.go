package cache

import "github.com/Goboolean/data-control-server/internal/domain/value"



func (q *stockQueue) Push(stock value.StockAggregate) {
	q.batch = append(q.batch, stock)
}

func (q *stockQueue) Clear() {
	q.batch = make([]value.StockAggregate, 0)
}