package main

import (
	"fmt"

	"github.com/AsynkronIT/goconsole"
	"github.com/hei6775/light-protoactor-go/actor"
)

type Hello struct{ Who string }
type SetBehaviorActor struct{}

func (state *SetBehaviorActor) Receive(context actor.Context) {
	switch msg := context.Message().(type) {
	case Hello:
		fmt.Printf("Hello %v\n", msg.Who)
		context.SetBehavior(state.Other)
	}
}

func (state *SetBehaviorActor) Other(context actor.Context) {
	switch msg := context.Message().(type) {
	case Hello:
		fmt.Printf("%v, ey we are now handling messages in another behavior", msg.Who)
	}
}

func NewSetBehaviorActor() actor.Actor {
	return &SetBehaviorActor{}
}

func main() {
	props := actor.FromProducer(NewSetBehaviorActor)
	pid := actor.Spawn(props)
	actor.Tell(pid, Hello{Who: "Roger"})
	actor.Tell(pid, Hello{Who: "Roger"})
	console.ReadLine()
}
