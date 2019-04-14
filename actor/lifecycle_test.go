package actor

import (
	"testing"
	"time"
)

func TestActorCanReplyOnStarting(t *testing.T) {
	future := NewFuture(testTimeout)
	a := Spawn(FromFunc(func(context Context) {
		switch context.Message().(type) {
		case *Started:
			context.Tell(future.PID(), EchoResponse{})
		}
	}))
	StopGraceful(a, 10*time.Second)
	assertFutureSuccess(future, t)
}

func TestActorCanReplyOnStopping(t *testing.T) {
	future := NewFuture(testTimeout)
	a := Spawn(FromFunc(func(context Context) {
		switch context.Message().(type) {
		case *Stopping:
			context.Tell(future.PID(), EchoResponse{})
		}
	}))
	StopGraceful(a, 10*time.Second)
	assertFutureSuccess(future, t)
}
