package config

import (
	"github.com/Goboolean/stock-fetch-server/internal/domain/service/persistence"
	"github.com/Goboolean/stock-fetch-server/internal/domain/service/relayer"
)

type ConfigurationManager struct {
	relayer     *relayer.RelayerManager
	persistence *persistence.PersistenceManager
}

var instance *ConfigurationManager

func init() {
	instance = &ConfigurationManager{
		relayer: relayer.New(),
	}
}

func NewConfigurationManager() *ConfigurationManager {
	return instance
}
