package config

import (
	"context"

	"github.com/Goboolean/fetch-server/internal/domain/entity"
	"github.com/Goboolean/fetch-server/internal/domain/errs"
)

func (m *ConfigurationManager) SetStockRelayableTrue(ctx context.Context, stockId string) error {
	return m.relayer.FetchStock(ctx, stockId)
}

func (m *ConfigurationManager) SetStockRelayableFalse(ctx context.Context, stockId string) error {
	return m.relayer.StopFetchingStock(ctx, stockId)
}

func (m *ConfigurationManager) SetStockStoreableTrue(ctx context.Context, stockId string) error {
	return m.persistence.SubscribeRelayer(ctx, stockId)
}

func (m *ConfigurationManager) SetStockStoreableFalse(ctx context.Context, stockId string) error {
	return m.persistence.UnsubscribeRelayer(ctx, stockId)
}

func (m *ConfigurationManager) SetStockTransmittableTrue(ctx context.Context, stockId string) error {
	return m.transmitter.SubscribeRelayer(ctx, stockId)
}

func (m *ConfigurationManager) SetStockTransmittableFalse(ctx context.Context, stockId string) error {
	return m.transmitter.UnsubscribeRelayer(ctx, stockId)
}


func (m *ConfigurationManager) GetStockConfiguration(ctx context.Context, stock string) (entity.StockConfiguration, error) {

	tx, err := m.tx.Transaction(context.Background())
	defer tx.Rollback()
	if err != nil {
		return entity.StockConfiguration{}, err
	}

	exists, err := m.db.CheckStockExists(tx, stock)
	if err != nil {
		return entity.StockConfiguration{}, err
	}

	if !exists {
		return entity.StockConfiguration{}, errs.StockNotFound
	}
	// check stock exist
	// reflect stock info to entity


	if err := tx.Commit(); err != nil {
		return entity.StockConfiguration{}, err
	}

	return entity.StockConfiguration{}, nil
}


func (m *ConfigurationManager) GetAllStockConfiguration(ctx context.Context) ([]entity.StockConfiguration, error) {

	tx, err := m.tx.Transaction(context.Background())
	defer tx.Rollback()
	if err != nil {
		return nil, err
	}

	metaList, err := m.db.GetAllStockMetadata(tx)
	if err != nil {
		return nil, err
	}

	confList := make([]entity.StockConfiguration, 0)

	for _, meta := range metaList {
		stockId := meta.StockID
		
		if !m.relayer.IsStockRelayable(stockId) {
			continue
		}

		isStoreable := m.persistence.IsStockStoreable(stockId)
		isTransmittable := m.transmitter.IsStockTransmittable(stockId)

		conf := entity.StockConfiguration{
			StockId: stockId,
			Relayable: true,
			Storeable: isStoreable,
			Transmittable: isTransmittable,
		}

		confList = append(confList, conf)
	}

	return confList, nil
}