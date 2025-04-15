package main

import (
	"log"
	"net"

	"github.com/lazypanda2004/notification-system/server"
	pb "github.com/lazypanda2004/notification-system/proto"
	"google.golang.org/grpc"
)

const (
	grpcPort     = ":50051"
	kafkaBroker  = "localhost:9092"
	kafkaTopic   = "notifications"
)

func main() {
	listener, err := net.Listen("tcp", grpcPort)
	if err != nil {
		log.Fatalf("Failed to listen on port %s: %v", grpcPort, err)
	}

	grpcServer := grpc.NewServer()

	notificationServer := server.NewNotificationServer([]string{kafkaBroker}, kafkaTopic)
	defer func() {
		if err := notificationServer.Close(); err != nil {
			log.Printf("Error closing Kafka producer: %v", err)
		}
	}()

	pb.RegisterNotificationServiceServer(grpcServer, notificationServer)

	log.Printf("gRPC server is running on port %s...", grpcPort)

	if err := grpcServer.Serve(listener); err != nil {
		log.Fatalf("Failed to serve gRPC server: %v", err)
	}
}
