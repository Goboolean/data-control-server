package persistence

import (
	"context"
	"errors"
	"log"
	"sync"
	"time"

	"github.com/Goboolean/fetch-server/internal/domain/entity"
)

func (m *PersistenceManager) SubscribeRelayer(stockId string) error {
	received := make([]*entity.StockAggregate, 0)

	if exists := m.s.StockExists(stockId); exists {
		return errors.New("stock already exists")
	}

	ctx := m.s.Map[stockId].Context()

	ch, err := m.relayer.Subscribe(ctx, stockId)
	if err != nil {
		return err
	}

	if err := m.s.StoreStock(stockId); err != nil {
		return err
	}


	var wg sync.WaitGroup

	cacheChan := make(chan struct{})
	defer func(wg *sync.WaitGroup) {
		wg.Wait()
		close(cacheChan)
	}(&wg)

	wg.Add(1)
	go func (ctx context.Context, wg *sync.WaitGroup) {
		defer wg.Done()

		for {
			select {
			case <- ctx.Done():
				return
			case <- time.After(m.o.syncDuration):
				cacheChan <- struct{}{}
			}
		}
	}(ctx, &wg)

	wg.Add(1)
	go func(ctx context.Context, wg *sync.WaitGroup) {
		defer wg.Done()

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
				}

				received = append(received, data)

				if len(received) % m.o.batchSize == 0 {
					ctx, cancel := context.WithCancel(m.s.Ctx.Context())
					defer cancel()

					if err := m.InsertStockOnCache(ctx, stockId, received); err != nil {
						log.Println(err)
	
						if err := m.db.CreateStoringFailedLog(ctx, stockId); err != nil {
							log.Println(err)
						}

						continue
					}

					received = received[:0]
				}

				case <- cacheChan:
					if err := m.SynchronizeCache(ctx, stockId); err != nil {
						log.Println(err)
					}
			}
		}
	}(ctx, &wg)

	return nil
}


func (m *PersistenceManager) UnsubscribeRelayer(stockId string) error {
	return m.s.UnstoreStock(stockId)
}


func (m *PersistenceManager) IsStockStoreable(stockId string) bool {
	return m.s.StockExists(stockId)
}


func (m *PersistenceManager) InsertStockOnDB(ctx context.Context, stockId string, batch []*entity.StockAggregate) error {

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


func (m *PersistenceManager) InsertStockOnCache(ctx context.Context, stockId string, batch []*entity.StockAggregate) error {
	return m.cache.StoreStockBatchOnCache(ctx, stockId, batch)
}


func (m *PersistenceManager) SynchronizeCache(ctx context.Context, stockId string) error {
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