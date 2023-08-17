package redis

import (
	"context"
	"fmt"
	"sync"

	"github.com/Goboolean/shared/pkg/resolver"
	"github.com/go-redis/redis/v8"
)

type Redis struct {
	client *redis.Client
}

var (
	instance *Redis
	once     sync.Once
)

func NewInstance(c *resolver.ConfigMap) (*Redis, error) {

	_, err := c.GetStringKey("USER")
	if err != nil {
		return nil, err
	}

	password, err := c.GetStringKey("PASSWORD")
	if err != nil {
		return nil, err
	}

	host, err := c.GetStringKey("HOST")
	if err != nil {
		return nil, err
	}

	port, err := c.GetStringKey("PORT")
	if err != nil {
		return nil, err
	}

	database, err := c.GetIntKey("DATABASE")
	if err != nil {
		return nil, err
	}

	address := fmt.Sprintf("%s:%s", host, port)

	once.Do(func() {

		rdb := redis.NewClient(&redis.Options{
			Addr:     address,
			Password: password,
			DB:       database,
		})

		if err := rdb.Ping(context.TODO()).Err(); err != nil {
			panic(err)
		}

		instance = &Redis{client: rdb}
	})

	return instance, nil
}

func (r *Redis) Close() error {
	if err := r.client.Close(); err != nil {
		return err
	}

	return nil
}

func (r *Redis) Ping() error {
	_, err := r.client.Ping(context.Background()).Result()
	return err
}