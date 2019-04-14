package actor

import "time"

//Stop the given PID
func StopActor(pid *PID) {
	pid.ref().Stop(pid)
}

// Tell sends a messages asynchronously to the PID
func Tell(pid *PID, message interface{}) {
	pid.ref().SendUserMessage(pid, message, nil)
}

// Request sends a messages asynchronously to the PID. The actor may send a response back via respondTo, which is
// available to the receiving actor via Context.Sender
func Request(pid *PID, message interface{}, respondTo *PID) {
	pid.ref().SendUserMessage(pid, message, respondTo)
}

// RequestFuture sends a message to a given PID and returns a Future
func RequestFuture(pid *PID, message interface{}, timeout time.Duration) *Future {
	future := NewFuture(timeout)
	pid.ref().SendUserMessage(pid, message, future.pid)
	return future
}

func StopFuture(pid *PID, timeout time.Duration) *Future {
	future := NewFuture(timeout)

	pid.sendSystemMessage(&Watch{Watcher: future.pid})
	StopActor(pid)

	return future
}

func StopGraceful(pid *PID, timeout time.Duration) {
	StopFuture(pid, timeout).Wait()
}
