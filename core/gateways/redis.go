package gateways

import (
	"context"
	"github.com/go-redis/redis/v8"
)

type RedisGateway interface {
	Write(ctx context.Context, key string, value string) error
	Read(ctx context.Context, key string) (string, error)
}

type redisGateway struct {
	client *redis.Client
}

func NewGateway() *redisGateway {
	client := redis.NewClient(&redis.Options{
		Addr:     "redis:6379",
		Password: "", // no password set
		DB:       0,  // use default DB
	})

	err := client.Ping(context.Background()).Err()
	if err != nil {
		panic(err)
	}

	return &redisGateway{
		client: client,
	}
}

func (g *redisGateway) Write(ctx context.Context, key string, value string) error {
	err := g.client.Set(ctx, key, value, 0).Err()
	if err != nil {
		return err
	}
	return nil
}

func (g *redisGateway) Read(ctx context.Context, key string) (string, error) {
	value, err := g.client.Get(ctx, key).Result()
	if err != nil {
		return "", err
	}
	return value, nil
}
