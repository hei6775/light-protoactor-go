package main

import (
	"runtime"

	"github.com/hei6775/light-protoactor-go/actor"
	"github.com/hei6775/light-protoactor-go/remote"

	"github.com/AsynkronIT/goconsole"
)

func main() {
	runtime.GOMAXPROCS(runtime.NumCPU())

	//empty actor just to have something to remote spawn
	props := actor.FromFunc(func(ctx actor.Context) {})
	remote.Register("remote", props)

	remote.Start("127.0.0.1:8080")

	console.ReadLine()
}
