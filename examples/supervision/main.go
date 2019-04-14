package main

import (
	"fmt"

	"sync"

	"github.com/hei6775/light-protoactor-go/actor"
	"github.com/AsynkronIT/goconsole"
)

type hello struct{ Who string }
type parentActor struct {
	*actor.Super

	childrenStoppedWg sync.WaitGroup
}

func (a *parentActor) onStarted(ctx actor.Context) {

}

func (a *parentActor) onStopping(ctx actor.Context) {

}

func (a *parentActor) onStopped(ctx actor.Context) {

}

func (a *parentActor) onReceive(ctx actor.Context) {
	switch msg := ctx.Message().(type) {
	case *hello:
		props := actor.FromInstance(newChildActor(&a.childrenStoppedWg))
		child := ctx.Spawn(props)
		actor.StopActor(child)
		ctx.Tell(child, msg)
	}
}

func newParentActor() actor.Actor {
	a := &parentActor{}
	a.Super = actor.NewSuper(
		a.onStarted,
		a.onStopping,
		a.onStopped,
		a.onReceive,
		nil,
		nil,
		nil,
		nil,
		nil,
	)
	return a
}

type childActor struct {
	*actor.Super
}

func (a *childActor) onStarted(ctx actor.Context) {

}

func (a *childActor) onStopping(ctx actor.Context) {
	panic("test")
}

func (a *childActor) onStopped(ctx actor.Context) {

}

func (a *childActor) onReceive(ctx actor.Context) {
	switch /*msg := */ ctx.Message().(type) {
	//case *actor.Started:
	//	fmt.Println("Starting, initialize actor here")
	//case *actor.Stopping:
	//	fmt.Println("Stopping, actor is about to shut down")
	//case *actor.Stopped:
	//	fmt.Println("Stopped, actor and its children are stopped")
	//case *actor.Restarting:
	//	fmt.Println("Restarting, actor is about to restart")
	case *hello:
		//fmt.Printf("Hello %v\n", msg.Who)
		panic("Ouch")
	}
}

func newChildActor(stoppedWg *sync.WaitGroup) actor.Actor {
	a := &childActor{}
	a.Super = actor.NewSuper(
		a.onStarted,
		a.onStopping,
		a.onStopped,
		a.onReceive,
		nil,
		nil,
		nil,
		nil,
		stoppedWg,
	)
	return a
}

func main() {
	decider := func(reason interface{}) actor.Directive {
		fmt.Println("handling failure for child")
		return actor.RestartDirective
	}
	supervisor := actor.NewOneForOneStrategy(0, 0, decider)
	props := actor.
		FromProducer(newParentActor).
		WithSupervisor(supervisor)

	pid := actor.Spawn(props)
	actor.Tell(pid, &hello{Who: "Roger"})
	console.ReadLine()
}
