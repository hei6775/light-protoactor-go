package router_test

import (
	"strconv"
	"sync"
	"testing"
	"time"

	actor "github.com/hei6775/light-protoactor-go/actor"
	"github.com/hei6775/light-protoactor-go/router"
)

type myMessage struct {
	i   int
	pid *actor.PID
}
type getRoutees struct {
	pid *actor.PID
}

func (m *myMessage) Hash() string {
	return strconv.Itoa(m.i)
}

var wait sync.WaitGroup

type routerActor struct{}
type tellerActor struct{}
type managerActor struct {
	set  []*actor.PID
	rpid *actor.PID
}

func (state *routerActor) Receive(context actor.Context) {
	switch msg := context.Message().(type) {
	case *myMessage:
		//log.Printf("%v got message %d", context.Self(), msg.i)
		msg.i++
		wait.Done()
	}
}
func (state *tellerActor) Receive(context actor.Context) {
	switch msg := context.Message().(type) {
	case *myMessage:
		for i := 0; i < 100; i++ {
			context.Tell(msg.pid, msg)
			time.Sleep(10 * time.Millisecond)
		}

	}
}

func (state *managerActor) Receive(context actor.Context) {
	switch msg := context.Message().(type) {
	case *router.Routees:
		state.set = msg.PIDs
		for i, v := range state.set {
			if i%2 == 0 {
				context.Tell(state.rpid, &router.RemoveRoutee{PID: v})
				//log.Println(v)

			} else {

				props := actor.FromInstance(&routerActor{})
				pid := actor.Spawn(props)
				context.Tell(state.rpid, &router.AddRoutee{PID: pid})
				//log.Println(v)
			}
		}
		context.Tell(context.Self(), &getRoutees{state.rpid})
	case *getRoutees:
		state.rpid = msg.pid
		context.Request(msg.pid, &router.GetRoutees{})
	}
}

func TestConcurrency(t *testing.T) {
	if testing.Short() {
		t.SkipNow()
	}

	wait.Add(100 * 10000)
	rpid := actor.Spawn(router.NewConsistentHashPool(100).WithInstance(&routerActor{}))

	props := actor.FromInstance(&tellerActor{})
	for i := 0; i < 10000; i++ {
		pid := actor.Spawn(props)
		actor.Tell(pid, &myMessage{i, rpid})
	}

	props = actor.FromInstance(&managerActor{})
	pid := actor.Spawn(props)
	actor.Tell(pid, &getRoutees{rpid})
	wait.Wait()
}
