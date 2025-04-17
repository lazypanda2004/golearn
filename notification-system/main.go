package main

import (
	"log"
	"net"
	"time"

	"github.com/lazypanda2004/notification-system/cmd/loadbalancer"
	"github.com/lazypanda2004/notification-system/internal/workerpool"
	"github.com/lazypanda2004/notification-system/internal/redis"
	"github.com/lazypanda2004/notification-system/server"
	pb "github.com/lazypanda2004/notification-system/proto"
	"google.golang.org/grpc"
)

const (
	grpcPort    = ":50051"
	kafkaBroker = "localhost:9092"
	kafkaTopic  = "notifications"
	redisAddr   = "localhost:6379"
	rateLimit   = 100
	timeWindow  = 1 * time.Minute
)

func main() {
	listener, err := net.Listen("tcp", grpcPort)
	if err != nil {
		log.Fatalf("Failed to listen grpc on port %s: %v", grpcPort, err)
	}
	grpcServer := grpc.NewServer()

	notificationServer := server.NewNotificationServer([]string{kafkaBroker}, kafkaTopic)
	defer func() {
		if err := notificationServer.Close(); err != nil {
			log.Printf("Error closing Kafka producer: %v", err)
		}
	}()
	pb.RegisterNotificationServiceServer(grpcServer, notificationServer)

	limiter := redis.NewLimiter(redisAddr, rateLimit, timeWindow)

	pool1 := workerpool.NewWorkerPool(8)
	pool2 := workerpool.NewWorkerPool(8)
	pool1.Start()
	pool2.Start()

	// --- Load balancer ---
	go func() {
		err := loadbalancer.Start([]string{kafkaBroker}, kafkaTopic, limiter, []*workerpool.WorkerPool{pool1, pool2})
		if err != nil {
			log.Fatalf("Load balancer error: %v", err)
		}
	}()

	log.Printf("gRPC server is running on port %s...", grpcPort)

	if err := grpcServer.Serve(listener); err != nil {
		log.Fatalf("Failed to serve gRPC server: %v", err)
	}
}
