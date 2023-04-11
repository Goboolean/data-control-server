package redis

import (
	"context"
	"fmt"
	"os"

	"github.com/go-redis/redis/v8"
)



var (
	REDIS_HOST     = os.Getenv("REDIS_HOST")
	REDIS_PORT     = os.Getenv("REDIS_PORT")
	REDIS_USER     = os.Getenv("REDIS_USER")
	REDIS_PASS     = os.Getenv("REDIS_PASS")

)

var instance *redis.Client


func NewInstance() *redis.Client {
	if instance == nil {
		rdb := redis.NewClient(&redis.Options{
			Addr:     fmt.Sprintf("%s:%s", REDIS_HOST, REDIS_PORT),
			Password: REDIS_PASS,
			Username: REDIS_USER,
			DB:       0,
		})

		result := rdb.Ping(context.TODO())
		if err := result.Err(); err != nil {
			panic(err)
		}

		instance = rdb

	}

	return instance
}