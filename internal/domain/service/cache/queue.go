package cache

import (
	"sync"

	"github.com/Goboolean/data-control-server/internal/domain/value"
)


type stockQueue struct {
	batch []value.StockAggregate

	mu sync.Mutex
}

func newQueue() *stockQueue {
	return &stockQueue{
		batch: make([]value.StockAggregate, 0),
	}
}

func (q *stockQueue) Lock() {
	q.mu.Lock()
}

func (q *stockQueue) Unlock() {
	q.mu.Unlock()
}