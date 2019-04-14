package main

import (
	"log"
	"runtime"

	"github.com/AsynkronIT/goconsole"
	"github.com/hei6775/light-protoactor-go/actor"
	"github.com/hei6775/light-protoactor-go/examples/remotebenchmark/messages"
	"github.com/hei6775/light-protoactor-go/mailbox"
	"github.com/hei6775/light-protoactor-go/remote"
)

func main() {
	runtime.GOMAXPROCS(runtime.NumCPU() * 1)
	runtime.GC()

	remote.Start("127.0.0.1:8080")
	var sender *actor.PID
	props := actor.
		FromFunc(
			func(context actor.Context) {
				switch msg := context.Message().(type) {
				case *messages.StartRemote:
					log.Println("Starting")
					sender = msg.Sender
					context.Respond(&messages.Start{})
				case *messages.Ping:
					context.Tell(sender, &messages.Pong{})
				}
			}).
		WithMailbox(mailbox.Bounded(1000000))

	actor.SpawnNamed(props, "remote")

	console.ReadLine()
}
