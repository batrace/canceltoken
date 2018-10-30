package canceltoken

import (
	"sync/atomic"
	"syscall"
	"testing"
	"time"

	. "github.com/smartystreets/goconvey/convey"
)

func busyWork(ms time.Duration) {
	time.Sleep(ms)
}

func TestCancelToken(t *testing.T) {
	Convey("", t, func() {

		ct := NewCancelToken()

		var exitcounter uint64

		for i := 0; i < 3; i++ {
			ct.Add(1)
			go func(t *CancelToken) {
				defer ct.Done()
				for !t.IsCancelled() {
					busyWork(50 * time.Millisecond)
				}
				atomic.AddUint64(&exitcounter, 1)
			}(ct)
		}

		// Sleep a brief time and then write to the internal
		// signal channel to simulate receiving an interrupt
		busyWork(500 * time.Millisecond)
		ct.c <- syscall.SIGINT

		sig := ct.Wait()
		So(sig, ShouldEqual, syscall.SIGINT)
		So(ct.isCancelled, ShouldBeTrue)
		So(atomic.LoadUint64(&exitcounter), ShouldEqual, 3)
	})
}
