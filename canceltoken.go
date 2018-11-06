package canceltoken

import (
	"os"
	"os/signal"
	"sync"
	"sync/atomic"
	"syscall"
)

// CancelToken recieves signals, cancels launched go routines and waits for them to exit
type CancelToken struct {
	c           chan os.Signal
	isCancelled uint32
	wg          sync.WaitGroup
}

// NewCancelToken makes a cancel token
func NewCancelToken() *CancelToken {
	t := CancelToken{}
	t.c = make(chan os.Signal, 1)
	signal.Notify(t.c, syscall.SIGINT, syscall.SIGTERM)
	atomic.StoreUint32(&t.isCancelled, 0)
	return &t
}

// Wait for the token to be cancelled and all WaitGroup members to exit
func (t *CancelToken) Wait() os.Signal {
	s := <-t.c
	atomic.StoreUint32(&t.isCancelled, 1)
	t.wg.Wait()
	return s
}

// IsCancelled is true when signal is recieved
func (t *CancelToken) IsCancelled() bool {
	cancelled := atomic.LoadUint32(&t.isCancelled)
	return (cancelled == 1)
}

// Add is called to increment the WaitGroup
func (t *CancelToken) Add(delta int) {
	t.wg.Add(delta)
}

// Done is called to decrement the WaitGroup
func (t *CancelToken) Done() {
	t.wg.Done()
}
