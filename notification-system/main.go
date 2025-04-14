package main

import (
	"log"
	"net"

	"github.com/lazypanda2004/notification-system/server"
	pb "github.com/lazypanda2004/notification-system/proto"
	"google.golang.org/grpc"
)

func main() {
	listener, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}

	grpcServer := grpc.NewServer()
	pb.RegisterNotificationServiceServer(grpcServer, &server.NotificationServer{})

	log.Println("gRPC server started on port 50051")
	if err := grpcServer.Serve(listener); err != nil {
		log.Fatalf("Failed to serve: %v", err)
	}
}
