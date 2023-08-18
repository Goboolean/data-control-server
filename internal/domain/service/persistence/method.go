package persistence

import (
	"context"
	"log"
	"time"

	"github.com/Goboolean/fetch-server/internal/domain/vo"
)



func (m *Manager) SubscribeRelayer(ctx context.Context, stockId string) error {
	received := make([]*vo.StockAggregate, 0)

	if err := m.s.StoreStock(stockId); err != nil {
		return err
	}

	ctx = m.s.Map[stockId].Context()

	ch, err := m.relayer.Subscribe(ctx, stockId)
	if err != nil {
		if err := m.s.UnstoreStock(stockId); err != nil {
			log.Println(err)
		}
		return err
	}


	useCache := m.o.useCache
	batchSize := m.o.BatchSize
	syncCount := m.o.SyncCount
	syncDuration := m.o.SyncDuration
	cacheChan := make(chan struct{}, 10)

	if m.o.SyncDuration != 0 {
		go func (ctx context.Context) {

			for {
				select {
				case <- ctx.Done():
					close(cacheChan)
					return
				case <- time.After(syncDuration):
					cacheChan <- struct{}{}
				}
			}
		}(ctx)
	}

	go func(ctx context.Context) {

		var cacheCount = 0

		if err := m.db.CreateStoringStartedLog(ctx, stockId); err != nil {
			log.Println(err)
		}

		for {
			select {
			case <- ctx.Done():
				if err := m.InsertStockOnDB(ctx, stockId, received); err != nil {
					log.Println(err)
				}

				if err := m.db.CreateStoringStoppedLog(ctx, stockId); err != nil {
					log.Println(err)
				}
				return

			case data, ok := <-ch:

				if !ok {
					m.db.CreateStoringStoppedLog(ctx, stockId)
					return
				}
				received = append(received, data)

				if len(received) % batchSize != 0 {
					continue
				}

				ctx, cancel := context.WithCancel(ctx)
				defer cancel()

				if useCache {
					if err := m.InsertStockOnCache(ctx, stockId, received); err != nil {
						log.Print(err)
						continue
					}

					cacheCount += len(received)
					if syncCount != 0 && cacheCount >= syncCount {
						cacheChan <- struct{}{}
					}

				} else {
					if err := m.InsertStockOnDB(ctx, stockId, received); err != nil {
						log.Print(err)
						continue
					}
				}

				received = received[:0]

			case <- cacheChan:
				if err := m.SynchronizeCache(ctx, stockId); err != nil {
					log.Println(err)
				}

				cacheCount = 0
			}
		}
	}(ctx)

	return nil
}


func (m *Manager) UnsubscribeRelayer(stockId string) error {
	return m.s.UnstoreStock(stockId)
}


func (m *Manager) IsStockStoreable(stockId string) bool {
	return m.s.StockExists(stockId)
}


func (m *Manager) InsertStockOnDB(ctx context.Context, stockId string, batch []*vo.StockAggregate) error {

	tx, err := m.tx.Transaction(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	if err := m.db.StoreStockBatch(tx, stockId, batch); err != nil {
		return err
	}

	return tx.Commit()
}


func (m *Manager) InsertStockOnCache(ctx context.Context, stockId string, batch []*vo.StockAggregate) error {
	return m.cache.StoreStockBatchOnCache(ctx, stockId, batch)
}


func (m *Manager) SynchronizeCache(ctx context.Context, stockId string) error {
	stockBatch, err := m.cache.GetAndEmptyCache(ctx, stockId)
	if err != nil {
		return err
	}

	if err := m.InsertStockOnDB(ctx, stockId, stockBatch); err != nil {
		if err := m.cache.StoreStockBatchOnCache(ctx, stockId, stockBatch); err != nil {
			log.Println(err)
		}
		return err
	}

	return nil
}