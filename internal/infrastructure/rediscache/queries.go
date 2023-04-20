package rediscache

import "github.com/go-redis/redis/v8"

type Queries struct {
	rds *redis.Client
	localLock bool
}

func New(client *redis.Client) *Queries {
	return &Queries{
		rds: client,
		localLock: false,
	}
}


