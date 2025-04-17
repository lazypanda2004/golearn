package workerpool

import (
	"context"
	"log"
	"net/smtp"
)

type Task struct {
	UserID    string
	Type      string 
	Recipient string 
	Message   string 
}


type WorkerPool struct {
	taskChan chan Task
	workers  int
	ctx      context.Context
	cancel   context.CancelFunc
}

func NewWorkerPool(workerCount int) *WorkerPool {
	ctx, cancel := context.WithCancel(context.Background())

	return &WorkerPool{
		taskChan: make(chan Task, 100), // buffered channel
		workers:  workerCount,
		ctx:      ctx,
		cancel:   cancel,
	}
}

// Start launches the workers
func (wp *WorkerPool) Start() {
	for i := range wp.workers {
		go wp.worker(i)
	}
}

// Stop gracefully shuts down the workers
func (wp *WorkerPool) Stop() {
	wp.cancel()
	close(wp.taskChan)
}

func (wp *WorkerPool) Submit(task Task) {
	wp.taskChan <- task
}

func (wp *WorkerPool) worker(id int) {
	log.Printf("Worker %d started", id)
	for {
		select {
		case <-wp.ctx.Done():
			log.Printf("Worker %d shutting down", id)
			return
		case task := <-wp.taskChan:
			wp.process(task, id)
		}
	}
}


func (wp *WorkerPool) process(task Task, workerID int) {
	log.Printf("Worker %d processing task: %+v", workerID, task)

	switch task.Type {
	case "email":
		wp.sendEmail(task.Recipient, task.Message, workerID)

	case "sms":
		wp.sendSMS(task.Recipient, task.Message, workerID)

	default:
		log.Printf("Worker %d: Unknown task type: %s", workerID, task.Type)
	}
}



func (wp *WorkerPool) sendSMS(to, message string, workerID int) {
	// Replace with actual SMS API logic
	log.Printf("Worker %d: SMS sent to %s - Message: %s", workerID, to, message)
}

func (wp *WorkerPool) sendEmail(to, body string, workerID int) {
	Username := "hrushikeshkareddy2004@gmail.com"
	password := "gigivczkykdbolms"
	smtpHost := "smtp.gmail.com"
	smtpPort := "587"
	subject := "Notification"

	auth := smtp.PlainAuth("", Username, password, smtpHost)

	msg := []byte("To: " + to + "\r\n" +
		"Subject: " + subject + "\r\n" +
		"MIME-Version: 1.0\r\n" +
		"Content-Type: text/html; charset=\"UTF-8\"\r\n" +
		"\r\n" + body + "\r\n")

	err := smtp.SendMail(smtpHost+":"+smtpPort, auth, Username, []string{to}, msg)
	if err != nil {
		log.Printf("Worker %d: Failed to send email to %s: %v", workerID, to, err)
		return
	}

	log.Printf("Worker %d: Email sent to %s", workerID, to)
}



