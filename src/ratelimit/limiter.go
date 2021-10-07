package ratelimit

import (
	"context"
	"fmt"
	"github.com/go-redis/redis/v8"
	"time"
)

const namespace = "rl"

func NewLimiter(client *redis.Client, action string, period time.Duration, limit int64) *Limiter {

	return &Limiter{client, action, period, limit}
}

type Limiter struct {
	client *redis.Client

	action string
	period time.Duration
	limit  int64
}

func (l *Limiter) CanDoAt(ctx context.Context, ts time.Time) (bool, error) {
	key := l.key(ts)

	var incr *redis.IntCmd
	_, err := l.client.Pipelined(ctx, func(p redis.Pipeliner) error {
		incr = p.Incr(ctx, key)
		p.Do(ctx, "PEXPIRE", key, 2*l.period.Milliseconds(), "NX")
		return nil
	})
	if err != nil {
		return false, err
	}
	count, err := incr.Result()
	if err != nil {
		return false, err
	}

	return count <= l.limit, nil
}

func (l *Limiter) key(ts time.Time) string {
	interval := ts.UTC().UnixNano() / l.period.Nanoseconds()
	return fmt.Sprintf("%s:%s:%x", namespace, l.action, interval)
}
