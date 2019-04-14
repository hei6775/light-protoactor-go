package main

import (
	"runtime"

	"github.com/AsynkronIT/goconsole"
	"gitee.com/lwj8507/light-protoactor-go/actor"
	"gitee.com/lwj8507/light-protoactor-go/examples/remoteactivate/messages"
	"gitee.com/lwj8507/light-protoactor-go/remote"
)

type helloActor struct{}

func (*helloActor) Receive(ctx actor.Context) {
	switch ctx.Message().(type) {
	case *messages.HelloRequest:
		ctx.Respond(&messages.HelloResponse{
			Message: "Hello from remote node",
		})
	}
}

func newHelloActor() actor.Actor {
	return &helloActor{}
}

func init() {
	remote.Register("hello", actor.FromProducer(newHelloActor))
}

func main() {
	runtime.GOMAXPROCS(runtime.NumCPU())

	remote.Start("127.0.0.1:8080")

	console.ReadLine()
}
