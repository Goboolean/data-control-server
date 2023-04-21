package rediscache

import "github.com/go-redis/redis/v8"

type Queries struct {
	rds *redis.Client
}

func New() *Queries {
	return &Queries{
		rds: NewInstance(),
	}
}


