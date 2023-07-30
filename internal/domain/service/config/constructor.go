package config

import (
	"sync"

	"github.com/Goboolean/fetch-server/internal/domain/port"
	"github.com/Goboolean/fetch-server/internal/domain/port/out"
	"github.com/Goboolean/fetch-server/internal/domain/service/persistence"
	"github.com/Goboolean/fetch-server/internal/domain/service/relayer"
	"github.com/Goboolean/fetch-server/internal/domain/service/transmission"
)

type ConfigurationManager struct {
	relayer     *relayer.RelayerManager
	persistence *persistence.PersistenceManager
	transmitter *transmission.Transmitter

	db out.StockMetadataPort
	tx port.TX
}

var (
	instance *ConfigurationManager
	once     sync.Once
)

func New(db out.StockMetadataPort, tx port.TX, r *relayer.RelayerManager, p *persistence.PersistenceManager, t *transmission.Transmitter) *ConfigurationManager {

	once.Do(func() {
		instance = &ConfigurationManager{
			relayer:     r,
			persistence: p,
			transmitter: t,

			db: db,
			tx: tx,
		}
	})

	return instance
}
