package main

import (
	"log"
	"runtime"

	"gitee.com/lwj8507/light-protoactor-go/actor"
	"gitee.com/lwj8507/light-protoactor-go/remote"

	"github.com/AsynkronIT/goconsole"

	"gitee.com/lwj8507/light-protoactor-go/examples/chat/messages"
)

func main() {
	//SetLogger(mylog.New("./log", "client"))

	runtime.GOMAXPROCS(runtime.NumCPU())

	if err := remote.Start("127.0.0.1:8081"); err != nil {
		panic(err)
	}

	server := actor.NewPID("127.0.0.1:8080", "chatserver")
	//spawn our chat client inline
	props := actor.FromFunc(func(context actor.Context) {
		switch msg := context.Message().(type) {
		case *messages.Connected:
			log.Println(msg.Message)
		case *messages.SayResponse:
			log.Printf("%v: %v", msg.UserName, msg.Message)
		case *messages.NickResponse:
			log.Printf("%v is now known as %v", msg.OldUserName, msg.NewUserName)
		}
	})

	client := actor.Spawn(props)

	actor.Tell(server, &messages.Connect{
		Sender: client,
	})

	nick := "Roger"
	cons := console.NewConsole(func(text string) {
		actor.Tell(server, &messages.SayRequest{
			UserName: nick,
			Message:  text,
		})
	})
	//write /nick NAME to change your chat username
	cons.Command("/nick", func(newNick string) {
		actor.Tell(server, &messages.NickRequest{
			OldUserName: nick,
			NewUserName: newNick,
		})
	})
	cons.Run()
	console.ReadLine()
}
