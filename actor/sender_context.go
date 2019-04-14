package actor

import "time"

type SenderContext interface {
	// Tell sends a messages asynchronously to the PID
	Tell(pid *PID, message interface{})

	// Request sends a messages asynchronously to the PID. The actor may send a response back via respondTo, which is
	// available to the receiving actor via Context.Sender
	Request(pid *PID, message interface{}, respondTo *PID)

	// RequestFuture sends a message to a given PID and returns a Future
	RequestFuture(pid *PID, message interface{}, timeout time.Duration) *Future
}

type rootContext struct {
}

func EmptyContext() SenderContext {
	return &rootContext{}
}

// Tell sends a messages asynchronously to the PID
func (*rootContext) Tell(pid *PID, message interface{}) {
	Tell(pid, message)
}

// Request sends a messages asynchronously to the PID. The actor may send a response back via respondTo, which is
// available to the receiving actor via Context.Sender
func (*rootContext) Request(pid *PID, message interface{}, respondTo *PID) {
	Request(pid, message, respondTo)
}

// RequestFuture sends a message to a given PID and returns a Future
func (*rootContext) RequestFuture(pid *PID, message interface{}, timeout time.Duration) *Future {
	return RequestFuture(pid, message, timeout)
}
