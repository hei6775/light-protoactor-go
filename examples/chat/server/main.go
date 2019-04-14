package main

import (
	"log"
	"runtime"

	//mylog "gitee.com/lwj8507/nggs/log"

	"github.com/hei6775/light-protoactor-go/actor"
	"github.com/hei6775/light-protoactor-go/remote"

	"github.com/AsynkronIT/goconsole"
	"github.com/emirpasic/gods/sets/hashset"

	"github.com/hei6775/light-protoactor-go/examples/chat/messages"
)

func notifyAll(context actor.Context, clients *hashset.Set, message interface{}) {
	for _, tmp := range clients.Values() {
		client := tmp.(*actor.PID)
		context.Tell(client, message)
	}
}

func main() {
	//SetLogger(mylog.New("./log", "server"))

	runtime.GOMAXPROCS(runtime.NumCPU())

	if err := remote.Start("127.0.0.1:8080"); err != nil {
		panic(err)
	}

	clients := hashset.New()
	props := actor.FromFunc(func(context actor.Context) {
		switch msg := context.Message().(type) {
		case *messages.Connect:
			log.Printf("client %v connected\n", msg.Sender)
			clients.Add(msg.Sender)
			context.Tell(msg.Sender, &messages.Connected{Message: "Welcome!"})
			context.Watch(msg.Sender)

		case *messages.SayRequest:
			log.Printf("recv SayResponse=[%v] from user=[%v]\n", msg.Message, msg.UserName)
			notifyAll(context, clients, &messages.SayResponse{
				UserName: msg.UserName,
				Message:  msg.Message,
			})

		case *messages.NickRequest:
			notifyAll(context, clients, &messages.NickResponse{
				OldUserName: msg.OldUserName,
				NewUserName: msg.NewUserName,
			})

		case *actor.Terminated:
			log.Printf("client[%v] terminated\n", msg.Who)
			clients.Remove(msg.Who)
			context.Unwatch(msg.Who)
		}
	})
	actor.SpawnNamed(props, "chatserver")
	console.ReadLine()
}
