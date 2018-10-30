package canceltoken

import (
	"os"
	"os/signal"
	"sync"
	"syscall"
)

// CancelToken recieves signals, cancels launched go routines and waits for them to exit
type CancelToken struct {
	c           chan os.Signal
	isCancelled bool
	wg          sync.WaitGroup
}

// NewCancelToken makes a cancel token
func NewCancelToken() *CancelToken {
	t := CancelToken{}
	t.c = make(chan os.Signal, 1)
	signal.Notify(t.c, syscall.SIGINT, syscall.SIGTERM)
	t.isCancelled = false
	return &t
}

// Wait for the token to be cancelled and all WaitGroup members to exit
func (t *CancelToken) Wait() os.Signal {
	s := <-t.c
	t.isCancelled = true
	t.wg.Wait()
	return s
}

// IsCancelled is true when signal is recieved
func (t *CancelToken) IsCancelled() bool {
	return t.isCancelled
}

// Add is called to increment the WaitGroup
func (t *CancelToken) Add(delta int) {
	t.wg.Add(delta)
}

// Done is called to decrement the WaitGroup
func (t *CancelToken) Done() {
	t.wg.Done()
}
