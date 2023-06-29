package stock

import (
	"fmt"

	"github.com/Goboolean/fetch-server/internal/domain/port"
	"github.com/Goboolean/fetch-server/internal/domain/value"
)

func (a *StockAdapter) GetStockType(tx port.Transactioner, stock string) (value.StockType, error) {

	types, err := a.postgres.GetStockOrigin(tx.Context(), stock)
	if err != nil {
		return 0, err
	}

	switch types {
	case "kor":
		return value.Domestic, nil
	case "usa":
		return value.International, nil
	default:
		return 0, fmt.Errorf("invalid stock origin")
	}
}

func (a *StockAdapter) StockExists(tx port.Transactioner, stock string) (bool, error) {
	return a.postgres.CheckStockExist(tx.Context(), stock)
}
