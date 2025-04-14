package server

import (
	"context"
	"log"

	pb "github.com/lazypanda2004/notification-system/proto"
)

type NotificationServer struct {
	pb.UnimplementedNotificationServiceServer
}

func (s *NotificationServer) SendNotification(ctx context.Context, req *pb.NotificationRequest) (*pb.NotificationResponse, error) {
	log.Printf("Received notification request for user: %s, type: %s", req.UserId, req.Type)

	// Dummy logic: always succeed
	return &pb.NotificationResponse{
		Success: true,
		Message: "Notification sent successfully!",
	}, nil
}
