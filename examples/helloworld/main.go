package main

import (
	"fmt"

	"github.com/AsynkronIT/goconsole"
	"github.com/hei6775/light-protoactor-go/actor"
)

type hello struct{ Who string }
type helloActor struct{}

func (state *helloActor) Receive(context actor.Context) {
	switch msg := context.Message().(type) {
	case *hello:
		fmt.Printf("Hello %v\n", msg.Who)
	}
}

func main() {
	props := actor.FromInstance(&helloActor{})
	pid := actor.Spawn(props)
	actor.Tell(pid, &hello{Who: "Roger"})
	console.ReadLine()
}
