package main

import (
	"log"
	"time"

	"github.com/hei6775/light-protoactor-go/actor"
	"github.com/hei6775/light-protoactor-go/remote"

	"github.com/AsynkronIT/goconsole"
)

type watch struct {
}

func main() {
	timeout := 5 * time.Second
	remote.Start("127.0.0.1:8081")

	props := actor.FromFunc(func(ctx actor.Context) {
		switch msg := ctx.Message().(type) {
		//case *actor.Started:
		//	log.Println("Local actor started")
		//	remotePID, err := remote.SpawnNamed("127.0.0.1:8080", "myRemote", "remote", timeout)
		//	if err != nil {
		//		log.Print("Local failed to spawn remote actor")
		//		return
		//	}
		//	log.Println("Local spawned remote actor")
		//	ctx.Watch(remotePID)
		//	log.Println("Local is watching remote actor")

		case *watch:
			log.Println("Local actor started")
			remotePID, err := remote.SpawnNamed("127.0.0.1:8080", "myRemote", "remote", timeout)
			if err != nil {
				log.Print("Local failed to spawn remote actor")
				return
			}
			log.Println("Local spawned remote actor")
			ctx.Watch(remotePID)
			log.Println("Local is watching remote actor")

		case *actor.Terminated:
			log.Printf("Local got terminated message %+v", msg)
		}
	})
	pid := actor.Spawn(props)

	var input string
	var err error
	var sendMsg = &watch{}

loop:
	for {
		input, err = console.ReadLine()
		if err != nil {
			break loop
		}
		switch input {
		case "exit", "quit", "bye":
			break loop
		default:
			actor.Tell(pid, sendMsg)
		}
	}
}
