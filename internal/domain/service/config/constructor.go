package config

import (
	"sync"

	"github.com/Goboolean/stock-fetch-server/internal/domain/service/persistence"
	"github.com/Goboolean/stock-fetch-server/internal/domain/service/relayer"
	"github.com/Goboolean/stock-fetch-server/internal/domain/service/transmission"
)

type ConfigurationManager struct {
	relayer     *relayer.RelayerManager
	persistence *persistence.PersistenceManager
	transmitter *transmission.Transmitter
}

var (
	instance *ConfigurationManager
	once sync.Once
)



func New(r *relayer.RelayerManager, p *persistence.PersistenceManager, t *transmission.Transmitter) *ConfigurationManager {

	once.Do(func() {
		instance = &ConfigurationManager{
			relayer: r,
			persistence: p,
			transmitter: t,
		}
	})

	return instance
}
