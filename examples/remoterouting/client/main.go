package main

import (
	"fmt"
	"log"
	"runtime"
	"sync"
	"time"

	console "github.com/AsynkronIT/goconsole"
	"gitee.com/lwj8507/light-protoactor-go/actor"
	"gitee.com/lwj8507/light-protoactor-go/examples/remoterouting/messages"
	"gitee.com/lwj8507/light-protoactor-go/mailbox"
	"gitee.com/lwj8507/light-protoactor-go/remote"
	"gitee.com/lwj8507/light-protoactor-go/router"
)

func main() {
	runtime.GOMAXPROCS(runtime.NumCPU())
	runtime.GC()

	remote.Start("127.0.0.1:8100")

	p1 := actor.NewPID("127.0.0.1:8101", "remote")
	p2 := actor.NewPID("127.0.0.1:8102", "remote")
	remotePID := actor.Spawn(router.NewConsistentHashGroup(p1, p2))

	messageCount := 1000000

	var wgStop sync.WaitGroup

	props := actor.
		FromProducer(newLocalActor(&wgStop, messageCount)).
		WithMailbox(mailbox.Bounded(10000))

	pid := actor.Spawn(props)

	log.Println("Starting to send")

	t := time.Now()

	for i := 0; i < messageCount; i++ {
		message := &messages.Ping{User: fmt.Sprintf("User_%d", i)}
		actor.Request(remotePID, message, pid)
	}

	wgStop.Wait()
	actor.StopActor(pid)

	fmt.Printf("elapsed: %v\n", time.Since(t))

	console.ReadLine()
}
