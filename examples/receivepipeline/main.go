package main

import (
	"fmt"

	"github.com/AsynkronIT/goconsole"
	"gitee.com/lwj8507/light-protoactor-go/actor"
	"gitee.com/lwj8507/light-protoactor-go/actor/middleware"
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
