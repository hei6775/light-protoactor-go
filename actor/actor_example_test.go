package actor_test

import (
	"fmt"
	"sync"
	"testing"
	"time"

	"gitee.com/lwj8507/light-protoactor-go/actor"
)

// Demonstrates how to create an actor using a function literal and how to send a message asynchronously
func TestExample(t *testing.T) {
	var props *actor.Props = actor.FromFunc(func(c actor.Context) {
		if msg, ok := c.Message().(string); ok {
			fmt.Println(msg) // outputs "Hello World"
		}
	})

	pid := actor.Spawn(props)

	actor.Tell(pid, "Hello World")
	actor.StopGraceful(pid, 10*time.Second) // wait for the actor to stop

	// Output: Hello World
}

// Demonstrates how to send a message from one actor to another and for the caller to wait for a response before
// proceeding
func TestExampleSynchronous(t *testing.T) {
	var wg sync.WaitGroup
	wg.Add(1)

	// callee will wait for the PING message
	callee := actor.Spawn(actor.FromFunc(func(c actor.Context) {
		if msg, ok := c.Message().(string); ok {
			fmt.Println(msg) // outputs PING
			c.Respond("PONG")
		}
	}))

	// caller will send a PING message and wait for the PONG
	caller := actor.Spawn(actor.FromFunc(func(c actor.Context) {
		switch msg := c.Message().(type) {
		// the first message an actor receives after it has started
		case *actor.Started:
			// send a PING to the callee, and specify the response
			// is sent to Self, which is this actor's PID
			c.Request(callee, "PING")

		case string:
			fmt.Println(msg) // PONG
			wg.Done()
		}
	}))

	wg.Wait()
	actor.StopGraceful(callee, 10*time.Second)
	actor.StopGraceful(caller, 10*time.Second)

	// Output:
	// PING
	// PONG
}
