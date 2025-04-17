// loadbalancer.go
package loadbalancer

import (
	"context"
	"encoding/json"
	"log"

	"github.com/segmentio/kafka-go"
	"github.com/lazypanda2004/notification-system/internal/redis"
	"github.com/lazypanda2004/notification-system/internal/workerpool"
)

type NotificationTask struct {
	UserID    string
	Type      string 
	Recipient string 
	Message   string 
}

func Start(kafkaBrokers []string, kafkaTopic string, limiter *redis.Limiter, pools []*workerpool.WorkerPool) error {
	ctx := context.Background()

	reader := kafka.NewReader(kafka.ReaderConfig{
		Brokers:  kafkaBrokers,
		Topic:    kafkaTopic,
		GroupID:  "load-balancer-group",
		MinBytes: 1,
		MaxBytes: 10e6,
	})
	log.Println("Load balancer started")

	workerSelector := 0

	for {
		m, err := reader.ReadMessage(ctx)
		if err != nil {
			log.Printf("Kafka read error: %v", err)
			continue
		}

		var task NotificationTask
		if err := json.Unmarshal(m.Value, &task); err != nil {
			log.Printf("Invalid message format: %v", err)
			continue
		}

		log.Printf("Received task for user %s", task.UserID)

		allowed, err := limiter.AllowOrQueue(ctx, redis.QueuedTask(task))
		if err != nil {
			log.Printf("Rate limit check failed: %v", err)
			continue
		}

		if allowed {
			assignedPool := pools[workerSelector%len(pools)]
			assignedPool.Submit(workerpool.Task(task))
			log.Printf("Task for user %s assigned to pool %d", task.UserID, workerSelector%len(pools))
			workerSelector++
		} else {
			log.Printf("Rate limit exceeded for user %s. Task queued in Redis.", task.UserID)
		}
	}
}
