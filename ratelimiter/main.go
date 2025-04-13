package main

import (
	"fmt"
	"sync"
	"time"
)

// Limiter interface
type Limiter interface {
	Allow() bool
}

// TokenBucketLimiter implements Limiter
type TokenBucketLimiter struct {
	capacity   int           // Maximum number of tokens
	tokens     int           // Current number of tokens
	rate       int           // Tokens added per second
	mutex      sync.Mutex    // To protect shared state
	refillDone chan struct{} // For stopping the refill goroutine
}

// NewTokenBucketLimiter creates a new TokenBucketLimiter
func NewTokenBucketLimiter(capacity, rate int) *TokenBucketLimiter {
	limiter := &TokenBucketLimiter{
		capacity:   capacity,
		tokens:     capacity,
		rate:       rate,
		refillDone: make(chan struct{}),
	}

	go limiter.refillTokens()
	return limiter
}

// refillTokens refills tokens at the given rate using a ticker
func (l *TokenBucketLimiter) refillTokens() {
	ticker := time.NewTicker(time.Second / time.Duration(l.rate))
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			l.mutex.Lock()
			if l.tokens < l.capacity {
				l.tokens++
			}
			l.mutex.Unlock()
		case <-l.refillDone:
			return
		}
	}
}

// Allow returns true if a token is available, false otherwise
func (l *TokenBucketLimiter) Allow() bool {
	l.mutex.Lock()
	defer l.mutex.Unlock()

	if l.tokens > 0 {
		l.tokens--
		return true
	}
	return false
}

// Stop stops the refill goroutine
func (l *TokenBucketLimiter) Stop() {
	close(l.refillDone)
}

// SimulateRequests simulates incoming requests and checks if allowed
func SimulateRequests(limiter Limiter, requestRate int, duration time.Duration) {
	ticker := time.NewTicker(time.Second / time.Duration(requestRate))
	defer ticker.Stop()

	end := time.After(duration)
	for {
		select {
		case <-ticker.C:
			if limiter.Allow() {
				fmt.Println("Request allowed at", time.Now().Format("15:04:05.000"))
			} else {
				fmt.Println("Request rejected at", time.Now().Format("15:04:05.000"))
			}
		case <-end:
			return
		}
	}
}

func main() {
	capacity := 5       // Max tokens
	rate := 2           // Refill rate per second
	requestRate := 4    // Incoming requests per second
	duration := 5 * time.Second

	limiter := NewTokenBucketLimiter(capacity, rate)
	defer limiter.Stop()

	SimulateRequests(limiter, requestRate, duration)
}


// If a type (TokenBucketLimiter) implements all the methods of an interface (Limiter), 
// it is considered to implement that interface — no explicit declaration needed.
// That’s it. Since *TokenBucketLimiter (the pointer receiver) has a method Allow() bool, it automatically implements the Limiter interface.

// You can swap TokenBucketLimiter with another implementation (e.g., FixedWindowLimiter, SlidingWindowLimiter) without changing SimulateRequests.

// so we can have multiple structs which implements the allow function then all of them can have same simulaterequests function right