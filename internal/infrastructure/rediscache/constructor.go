package rediscache

import (
	"context"
	"fmt"
	"strconv"
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

func NewInstance(c *resolver.Config) *Redis {

	if err := c.ShouldHostExist(); err != nil {
		panic(err)
	}

	if err := c.ShouldPortExist(); err != nil {
		panic(err)
	}

	if err := c.ShouldUserExist(); err != nil {
		panic(err)
	}

	if err := c.ShouldPWExist(); err != nil {
		panic(err)
	}

	if err := c.ShouldDBExist(); err != nil {
		panic(err)
	}

	c.Address = fmt.Sprintf("%s:%s", c.Host, c.Port)

	database, err := strconv.Atoi(c.Database)
	if err != nil {
		panic(err)
	}

	once.Do(func() {

		rdb := redis.NewClient(&redis.Options{
			Addr:     c.Address,
			Password: c.Password,
			Username: c.User,
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
