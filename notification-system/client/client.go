package main

import (
	"context"
	"fmt"
	"log"
	"sync"
	"time"
	"math/rand"

	pb "github.com/lazypanda2004/notification-system/proto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

const (
	numUsers    = 5
	requestRate = 300 * time.Millisecond // 2 request per second
	testDuration = 3 * time.Second // run for 1 minute
)

func simulateUser(userID string, client pb.NotificationServiceClient, wg *sync.WaitGroup) {
	defer wg.Done()

	ticker := time.NewTicker(requestRate)
	defer ticker.Stop()

	end := time.After(testDuration)

	for {
		select {
		case <-end:
			return
		case <-ticker.C:
			// Send request
			typei := rand.Intn(2)
			var reqType string
			if typei == 0 {
				reqType = "email"
			} else {
				reqType = "sms"
			}
			req := &pb.NotificationRequest{
				UserId:    userID,
				Type:      reqType,
				Recipient: fmt.Sprintf("%s@example.com", userID),
				Message:   fmt.Sprintf("Hello from %s at %s", userID, time.Now().Format(time.RFC3339)),
			}

			ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
			res, err := client.SendNotification(ctx, req)
			defer cancel()

			if err != nil {
				log.Printf("%s error: %v", userID, err)
				continue
			}

			log.Printf("%s => %v | %s", userID, res.Success, res.Message)
		}
	}
}

func main() {
	conn, err := grpc.NewClient("localhost:50051", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("Failed to connect to server: %v", err)
	}
	defer conn.Close()

	client := pb.NewNotificationServiceClient(conn)

	var wg sync.WaitGroup

	// Simulate multiple users
	for i := 1; i <= numUsers; i++ {
		wg.Add(1)
		go simulateUser(fmt.Sprintf("user_%d", i), client, &wg)
	}

	wg.Wait()
	log.Println("Load simulation completed.")
}
