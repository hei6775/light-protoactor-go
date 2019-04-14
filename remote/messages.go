package remote

import "gitee.com/lwj8507/light-protoactor-go/actor"

type EndpointTerminatedEvent struct {
	Address string
}

type EndpointReaderFailedToRead struct {
	Err error
}

type remoteWatch struct {
	Watcher *actor.PID
	Watchee *actor.PID
}

type remoteUnwatch struct {
	Watcher *actor.PID
	Watchee *actor.PID
}

type remoteTerminate struct {
	Watcher *actor.PID
	Watchee *actor.PID
}

var (
	stopMessage interface{} = &actor.Stop{}
)
