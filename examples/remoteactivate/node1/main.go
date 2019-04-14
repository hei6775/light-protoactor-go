package main

import (
	"fmt"
	"time"

	console "github.com/AsynkronIT/goconsole"
	"gitee.com/lwj8507/light-protoactor-go/actor"
	"gitee.com/lwj8507/light-protoactor-go/examples/remoteactivate/messages"
	"gitee.com/lwj8507/light-protoactor-go/remote"
)

func main() {
	timeout := 5 * time.Second
	remote.Start("127.0.0.1:8081")
	pid, _ := remote.SpawnNamed("127.0.0.1:8080", "remote", "hello", timeout)
	res, _ := actor.RequestFuture(pid, &messages.HelloRequest{}, timeout).Result()
	response := res.(*messages.HelloResponse)
	fmt.Printf("Response from remote %v", response.Message)

	console.ReadLine()
}
