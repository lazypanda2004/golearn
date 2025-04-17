package main

import (
	"context"
	"fmt"
	"log"
	"sync"
	"time"

	pb "github.com/lazypanda2004/notification-system/proto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

const (
	numUsers     = 5
	requestRate  = 200 * time.Millisecond // ~3 requests per second
	testDuration = 10 * time.Second        // run for 5 seconds
)

const html = `
<table align="center" width="100%" cellpadding="0" cellspacing="0" style="padding: 20px;">
  <tr>
    <td>
      <table align="center" width="600" cellpadding="0" cellspacing="0" style="background-color: #ffffff; border-radius: 10px; overflow: hidden; box-shadow: 0 0 10px rgba(0,0,0,0.1);">
        <tr>
          <td style="background-color: #fcc70; padding: 30px; text-align: center;">
            <h1 style="margin: 0; color: #333;">🎉 Happy Birthday ra PUKA! 🎉</h1>
          </td>
        </tr>
        <tr>
          <td style="padding: 30px; text-align: center;">
            <img src="https://media1.tenor.com/m/WOqAMSYnCGoAAAAd/telugu-brahmanandam.gif" alt="Birthday Cake" width="300" height = "300" style="border-radius: 10px; margin-bottom: 20px;" />
            <p style="font-size: 18px; color: #555;">
              Wishing you a day filled with happiness and a year filled with joy.
            </p>
            <p style="font-size: 16px; color: #777;">
              Bokka le kani PARTY ekkada!
            </p>
          </td>
        </tr>
        <tr>
          <td style="background-color: #fcc70; padding: 20px; text-align: center;">
            <p style="margin: 0; font-size: 14px; color: #444;">
              With love,<br>Your BABE 💌
            </p>
          </td>
        </tr>
      </table>
    </td>
  </tr>
</table>
`

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
			reqType := "email"
			req := &pb.NotificationRequest{
				UserId:    userID,
				Type:      reqType,
				Recipient: "g.prashams@iitg.ac.in",
				
				Message:   html,
			}

			ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
			res, err := client.SendNotification(ctx, req)
			cancel()

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

	for i := 1; i <= numUsers; i++ {
		wg.Add(1)
		go simulateUser(fmt.Sprintf("user_%d", i), client, &wg)
	}

	wg.Wait()
	log.Println("Load simulation completed.")
}
