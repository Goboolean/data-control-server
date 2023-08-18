package config

import (
	"sync"

	"github.com/Goboolean/fetch-server/internal/domain/port"
	"github.com/Goboolean/fetch-server/internal/domain/port/out"
	"github.com/Goboolean/fetch-server/internal/domain/service/persistence"
	"github.com/Goboolean/fetch-server/internal/domain/service/relay"
	"github.com/Goboolean/fetch-server/internal/domain/service/transmission"
)

type Configurator struct {
	relayer     *relay.Manager
	persistence *persistence.Manager
	transmitter *transmission.Manager

	db out.StockMetadataPort
	tx port.TX
}

var (
	instance *Configurator
	once     sync.Once
)

func New(db out.StockMetadataPort, tx port.TX, r *relay.Manager, p *persistence.Manager, t *transmission.Manager) (*Configurator, error) {

	once.Do(func() {
		instance = &Configurator{
			relayer:     r,
			persistence: p,
			transmitter: t,

			db: db,
			tx: tx,
		}
	})

	return instance, nil
}
