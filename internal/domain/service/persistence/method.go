package persistence

import (
	"context"
	"log"

	"github.com/Goboolean/fetch-server/internal/domain/entity"
)

func (m *PersistenceManager) SubscribeRelayer(ctx context.Context, stockId string) error {

	ch, err := m.relayer.Subscribe(stockId)
	if err != nil {
		return err
	}

	if err := m.s.StoreStock(stockId); err != nil {
		return err
	}

	go func() {

		if err := m.db.CreateStoringStartedLog(ctx, stockId); err != nil {
			log.Println(err)
		}

		for {
			select {
			case <- m.s.Map[stockId].Done():
				if err := m.db.CreateStoringStoppedLog(ctx, stockId); err != nil {
					log.Println(err)
				}
				return

			case data := <-ch:

				ctx := context.Background()

				if err := m.InsertStockOnDB(ctx, stockId, data); err != nil {
					log.Println(err)

					if err := m.db.CreateStoringFailedLog(ctx, stockId); err != nil {
						log.Println(err)
					}
					return
				}
			}
		}
	}()

	return nil
}


func (m *PersistenceManager) UnsubscribeRelayer(ctx context.Context, stockId string) error {
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

