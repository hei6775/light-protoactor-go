package actor

import (
	"log"
	"sync"
	"testing"
	"time"

	//mylog "gitee.com/lwj8507/nggs/log"
)

func startLog() {
	//SetLogger(mylog.New("./log", "super_test"))
}

func stopLog() {
	logger.Close()
}

type example struct {
	*Super

	startedWg sync.WaitGroup
	stoppedWg sync.WaitGroup
}

func newExample() *example {
	e := &example{}

	e.Super = NewSuper(
		e.OnStarted,
		e.OnStopping,
		e.OnStopped,
		e.OnReceiveMessage,
		nil,
		e.OnRestarting,
		e.OnRestarted,
		&e.startedWg,
		&e.stoppedWg,
	)

	return e
}

func (e *example) OnStarted(ctx Context) {
	log.Println("OnStarted")

	//e.NewLoopTimer(1*time.Second, func(timer.TimerID) {
	//	actor.Tell(e.PID, "receive message")
	//})

	e.NewTimer(3*time.Second, func(TimerID) {
		panic("test")
	})
}

func (e *example) OnStopping(ctx Context) {
	log.Println("OnStopping")
}

func (e *example) OnStopped(ctx Context) {
	log.Println("OnStopped")
}

func (e *example) OnReceiveMessage(ctx Context) {
	log.Println("OnReceiveMessage")
}

func (e *example) OnRestarting(ctx Context) {
	log.Println("OnRestarting")
}

func (e *example) OnRestarted(ctx Context) {
	log.Println("OnRestarted")
	panic("test")
}

func (e *example) start() {
	Spawn(FromInstance(e))
}

func (e *example) waitForStarted() {
	e.startedWg.Wait()
}

func (e *example) stop() {
	StopActor(e.PID)
}

func (e *example) waitForStopped() {
	e.stoppedWg.Wait()
}

func TestSuper_StartStop(t *testing.T) {
	startLog()

	e := newExample()

	e.start()

	e.waitForStarted()

	//util.WaitExitSignal()
	//time.Sleep(10 * time.Second)

	//e.stop()

	e.waitForStopped()

	stopLog()
}
