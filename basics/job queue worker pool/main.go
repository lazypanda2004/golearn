package main

import (
	"fmt"
	"log"
	"sync"
	"time"
)

// Runnable interface for all jobs
type Runnable interface {
	Run()
}

// Job struct implementing Runnable
type Job struct {
	ID int
}

func (j *Job) Run() {
	log.Printf("Running Job ID: %d", j.ID)
	time.Sleep(time.Second) // Simulate work
	log.Printf("Completed Job ID: %d", j.ID)
}

// Worker struct
type Worker struct {
	ID        int
	JobQueue  <-chan Runnable
	WaitGroup *sync.WaitGroup
	Counter   *JobCounter
}

// JobCounter to track completed jobs safely
type JobCounter struct {
	sync.Mutex
	Count int
}

func (w *Worker) Start() {
	go func() {
		for job := range w.JobQueue {
			job.Run()
			w.Counter.Increment()
			w.WaitGroup.Done()
		}
	}()
}

// Increment safely updates job counter
func (jc *JobCounter) Increment() {
	jc.Lock()
	defer jc.Unlock()
	jc.Count++
}

// Dispatcher struct
type Dispatcher struct {
	WorkerPoolSize int
	JobQueue       chan Runnable
	WaitGroup      sync.WaitGroup
	Counter        JobCounter
}

// NewDispatcher creates and starts workers
func NewDispatcher(poolSize int) *Dispatcher {
	jobQueue := make(chan Runnable)

	dispatcher := &Dispatcher{
		WorkerPoolSize: poolSize,
		JobQueue:       jobQueue,
		Counter:        JobCounter{},
	}

	for i := 1; i <= poolSize; i++ {
		worker := Worker{
			ID:        i,
			JobQueue:  jobQueue,
			WaitGroup: &dispatcher.WaitGroup,
			Counter:   &dispatcher.Counter,
		}
		worker.Start()
	}

	return dispatcher
}

// Submit adds jobs to the queue
func (d *Dispatcher) Submit(job Runnable) {
	d.WaitGroup.Add(1)
	d.JobQueue <- job
}

// Wait blocks until all jobs are done
func (d *Dispatcher) Wait() {
	d.WaitGroup.Wait()
}

// Stop gracefully closes the queue
func (d *Dispatcher) Stop() {
	close(d.JobQueue)
}

func main() {
	dispatcher := NewDispatcher(3) // 3 worker goroutines

	// Submit 10 jobs
	for i := 1; i <= 10; i++ {
		dispatcher.Submit(&Job{ID: i})
	}

	dispatcher.Wait()  // Wait for all jobs to finish
	dispatcher.Stop()  // Close job queue

	fmt.Printf("\n📊 Total Jobs Completed: %d\n", dispatcher.Counter.Count)
}
