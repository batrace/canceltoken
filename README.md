CancelToken is a simple helper module that allows a main go process to launch multiple go routines and then upon receiving a signal, cancell those go routines and wait for their exit.

Example usage:

```go
package main

import (
	"github.com/batrace/canceltoken"
)

func main() {
	ct := canceltoken.NewCancelToken()
	ct.Add(1)
	go func(t *canceltoken.CancelToken) {
		defer ct.Done()
		for !t.IsCancelled() {
			// Do some work
		}
	}(ct)

	ct.Wait()
}
```