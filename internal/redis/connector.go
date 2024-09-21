// Package redis provides a Redis connector.
package redis

import (
	"context"

	"github.com/VinukaThejana/go-utils/logger"
	env "github.com/VinukaThejana/todoapp/internal/config"
	"github.com/redis/go-redis/v9"
)

// Init initializes a Redis connection.
func Init(e *env.Env) *redis.Client {
	opt, err := redis.ParseURL(e.RedisURL)
	if err != nil {
		logger.Errorf(err)
	}

	rdb := redis.NewClient(opt)
	if err := rdb.Ping(context.Background()).Err(); err != nil {
		logger.Errorf(err)
	}

	return rdb
}
