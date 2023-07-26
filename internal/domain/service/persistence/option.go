package persistence

import (
	"time"
)


type Option struct {
	// system processes batch of data at a time.
	// if given value is zero(default), it does not use batch.
	BatchSize    int
	// system tries to vacate cache and store it to db with given duration.
	// if given value is zero(default), it does not vacate periodically.
	SyncDuration time.Duration
	// if data is stored given times at a cache,
	// system tries to vacate cache and store it to db.
	// if given value is zero(default), it does not use cacahe.
	SyncCount    int
	// if both syncDuration and syncCount is set,
	// both cache vacate logic is executed.
	useCache     bool
}


func (m *PersistenceManager) SetBatchSize(batchSize int) {
	m.o.BatchSize = batchSize
}

func (m *PersistenceManager) SetSyncDuration(syncDuration time.Duration) {
	m.o.SyncDuration = syncDuration
	m.setUseCache()
}

func (m *PersistenceManager) SetSyncCount(syncCount int) {
	m.o.SyncCount = syncCount
	m.setUseCache()
}

func (m *PersistenceManager) setUseCache() {
	m.o.useCache = m.o.SyncDuration != 0 || m.o.SyncCount != 0 
}