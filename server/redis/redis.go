package rd

import (
	"fmt"

	"github.com/go-redis/redis"
)

type RedisStorage interface {
	Set(key string, value string) (string, error)
	Get(key string) (string, error)
}

type redisStorage struct {
	client redis.Cmdable
}

func New() RedisStorage {
	return &redisStorage{
		client: redis.NewClient(&redis.Options{
			Addr: "redis:6379",
			DB:   0,
		}),
	}
}

func (r *redisStorage) Set(key string, value string) (string, error) {
	command := r.client.Set(key, value, 0)
	if err := command.Err(); err != nil {
		return "", fmt.Errorf("failed to set : %w", err)
	}
	return command.Val(), nil
}

func (r *redisStorage) Get(key string) (string, error) {
	command := r.client.Get(key)
	if err := command.Err(); err != nil {
		return "", fmt.Errorf("failed to get : %w", err)
	}
	return command.Val(), nil
}
