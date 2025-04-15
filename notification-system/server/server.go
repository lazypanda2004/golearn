package server

import (
	"context"
	"encoding/json"
	"log"
	"time"

	"github.com/segmentio/kafka-go"
	pb "github.com/lazypanda2004/notification-system/proto"
)

type NotificationServer struct {
	pb.UnimplementedNotificationServiceServer
	kafkaWriter *kafka.Writer
}

func NewNotificationServer(brokers []string, topic string) *NotificationServer {
	writer := &kafka.Writer{
		Addr:         kafka.TCP(brokers...),
		Topic:        topic,
		Balancer:     &kafka.LeastBytes{},
		RequiredAcks: kafka.RequireAll,
		Async:        false,
	}

	return &NotificationServer{
		kafkaWriter: writer,
	}
}

func (s *NotificationServer) SendNotification(ctx context.Context, req *pb.NotificationRequest) (*pb.NotificationResponse, error) {
	log.Printf("Received notification request for user: %s, type: %s", req.UserId, req.Type)

	// Serialize the request to JSON
	data, err := json.Marshal(req)
	if err != nil {
		log.Printf("Failed to marshal request: %v", err)
		return &pb.NotificationResponse{
			Success: false,
			Message: "Failed to process request",
		}, nil
	}

	// Create Kafka message
	msg := kafka.Message{
		Key:   []byte(req.UserId),
		Value: data,
		Time:  time.Now(),
	}

	// Publish to Kafka
	err = s.kafkaWriter.WriteMessages(ctx, msg)
	if err != nil {
		log.Printf("Failed to write to Kafka: %v", err)
		return &pb.NotificationResponse{
			Success: false,
			Message: "Failed to publish notification",
		}, nil
	}

	return &pb.NotificationResponse{
		Success: true,
		Message: "Notification queued successfully!",
	}, nil
}

func (s *NotificationServer) Close() error {
	return s.kafkaWriter.Close()
}
