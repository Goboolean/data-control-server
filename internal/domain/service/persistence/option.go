package persistence

import (
	"time"
)


type Option struct {
	// system processes batch of data at a time.
	// if given value is zero(default), it does not use batch.
	batchSize    int
	// system tries to vacate cache and store it to db with given duration.
	// if given value is zero(default), it does not vacate periodically.
	syncDuration time.Duration
	// if data is stored given times at a cache,
	// system tries to vacate cache and store it to db.
	// if given value is zero(default), it does not use cacahe.
	syncCount    int
	// if both syncDuration and syncCount is set,
	// both cache vacate logic is executed.
}


func (m *PersistenceManager) SetBatchSize(batchSize int) {
	m.o.batchSize = batchSize
}

func (m *PersistenceManager) SetSyncDuration(syncDuration time.Duration) {
	m.o.syncDuration = syncDuration
}

func (m *PersistenceManager) SetSyncCount(syncCount int) {
	m.o.syncCount = syncCount
}