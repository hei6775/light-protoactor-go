package main

import (
	"fmt"

	"github.com/AsynkronIT/goconsole"
	"github.com/hei6775/light-protoactor-go/actor"
	"github.com/hei6775/light-protoactor-go/actor/middleware"
)

type hello struct{ Who string }

func receive(context actor.Context) {
	switch msg := context.Message().(type) {
	case *hello:
		fmt.Printf("Hello %v\n", msg.Who)
	}
}

func main() {
	props := actor.FromFunc(receive).WithMiddleware(middleware.Logger)
	pid := actor.Spawn(props)
	actor.Tell(pid, &hello{Who: "Roger"})
	console.ReadLine()
}
