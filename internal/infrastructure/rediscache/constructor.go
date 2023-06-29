package rediscache

import (
	"context"
	"fmt"
	"sync"

	"github.com/Goboolean/shared/pkg/resolver"
	"github.com/go-redis/redis/v8"
)

type Redis struct {
	*redis.Client
}

var (
	instance *Redis
	once     sync.Once
)

func NewInstance(c *resolver.ConfigMap) *Redis {

	user, err := c.GetStringKey("USER")
	if err != nil {
		panic(err)
	}

	password, err := c.GetStringKey("PASSWORD")
	if err != nil {
		panic(err)
	}

	host, err := c.GetStringKey("HOST")
	if err != nil {
		panic(err)
	}

	port, err := c.GetStringKey("PORT")
	if err != nil {
		panic(err)
	}

	database, err := c.GetIntKey("DATABASE")
	if err != nil {
		panic(err)
	}

	address := fmt.Sprintf("%s:%s", host, port)

	if err != nil {
		panic(err)
	}

	once.Do(func() {

		rdb := redis.NewClient(&redis.Options{
			Addr:     address,
			Password: password,
			Username: user,
			DB:       database,
		})

		if err := rdb.Ping(context.TODO()).Err(); err != nil {
			panic(err)
		}
	})

	return instance
}

func (r *Redis) Close() error {
	if err := r.Close(); err != nil {
		return err
	}

	return nil
}

func (r *Redis) Ping() error {
	return r.Ping()
}

func (r *Redis) NewTx(ctx context.Context) (resolver.Transactioner, error) {
	return NewTransaction(r.Pipeline(), ctx), nil
}
