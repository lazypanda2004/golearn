package redis

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
)

type Limiter struct {
	rdb        *redis.Client
	rateLimit  int
	timeWindow time.Duration
}

type QueuedTask struct {
	UserID    string
	Type      string 
	Recipient string 
	Message   string 
}

func NewLimiter(addr string, rateLimit int, timeWindow time.Duration) *Limiter {
	rdb := redis.NewClient(&redis.Options{
		Addr: addr,
	})
	return &Limiter{
		rdb:        rdb,
		rateLimit:  rateLimit,
		timeWindow: timeWindow,
	}
}

func (l *Limiter) AllowOrQueue(ctx context.Context, task QueuedTask) (bool, error) {
	key := fmt.Sprintf("rate_limit:%s", task.UserID)
	count, err := l.rdb.Incr(ctx, key).Result()
	if err != nil {
		return false, err
	}

	if count == 1 {
		l.rdb.Expire(ctx, key, l.timeWindow)
	}

	if int(count) > l.rateLimit {
		// Exceeded limit - queue the task
		queueKey := fmt.Sprintf("queue:%s", task.UserID)
		data, err := json.Marshal(task)
		if err != nil {
			return false, err
		}
		err = l.rdb.RPush(ctx, queueKey, data).Err()
		if err != nil {
			return false, err
		}
		return false, nil
	}
	return true, nil
}

func (l *Limiter) PopQueuedTask(ctx context.Context, userID string) (*QueuedTask, error) {
	queueKey := fmt.Sprintf("queue:%s", userID)
	data, err := l.rdb.LPop(ctx, queueKey).Result()
	if err == redis.Nil {
		return nil, nil // empty
	} else if err != nil {
		return nil, err
	}

	var task QueuedTask
	err = json.Unmarshal([]byte(data), &task)
	if err != nil {
		return nil, err
	}
	return &task, nil
}
