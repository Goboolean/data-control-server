package rediscache

import (
	"context"
	"fmt"
	"os"
	"strconv"

	"github.com/go-redis/redis/v8"
)



var (
	REDIS_HOST     = os.Getenv("REDIS_HOST")
	REDIS_PORT     = os.Getenv("REDIS_PORT")
	REDIS_USER     = os.Getenv("REDIS_USER")
	REDIS_PASS     = os.Getenv("REDIS_PASS")
	REDIS_DATABASE, err = strconv.Atoi(os.Getenv("REDIS_DATABASE"))
)

func init() {

	if err != nil {
		panic(err)
	}

	rdb := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%s", REDIS_HOST, REDIS_PORT),
		Password: REDIS_PASS,
		Username: REDIS_USER,
		DB:       REDIS_DATABASE,
	})

	result := rdb.Ping(context.TODO())
	if err := result.Err(); err != nil {
		panic(err)
	}

	instance = rdb
}



var instance *redis.Client

func NewInstance() *redis.Client {
	return instance
}

func Close() error {
	if err := instance.Close(); err != nil {
		return err
	}
	
	return nil
}

