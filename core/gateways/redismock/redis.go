package redismock

import "context"

type RedisGateway struct {
	WriteFunc func(ctx context.Context, key string, value string) error
	ReadFunc  func(ctx context.Context, key string) (string, error)
}

func NewGateway() *RedisGateway {
	return &RedisGateway{}
}

func (g *RedisGateway) Write(ctx context.Context, key string, value string) error {
	return g.WriteFunc(ctx, key, value)
}

func (g *RedisGateway) Read(ctx context.Context, key string) (string, error) {
	return g.ReadFunc(ctx, key)
}
