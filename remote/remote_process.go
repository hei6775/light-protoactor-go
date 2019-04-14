package remote

import (
	"gitee.com/lwj8507/light-protoactor-go/actor"

	"github.com/gogo/protobuf/proto"
)

type remoteProcess struct {
	pid *actor.PID
}

func newRemoteProcess(pid *actor.PID) actor.Process {
	return &remoteProcess{
		pid: pid,
	}
}

func (ref *remoteProcess) SendUserMessage(pid *actor.PID, message interface{}, sender *actor.PID) {
	sendRemoteMessage(pid, message, sender)
}

func sendRemoteMessage(pid *actor.PID, message interface{}, sender *actor.PID) {
	switch msg := message.(type) {
	case proto.Message:

		rd := &remoteDeliver{
			message: msg,
			sender:  sender,
			target:  pid,
		}
		actor.Tell(endpointManagerPID, rd)
	default:
		logger.Error("failed, trying to send non Proto message[%#v], pid=[%v]", msg, pid)
	}
}

func (ref *remoteProcess) SendSystemMessage(pid *actor.PID, message interface{}) {

	//intercept any Watch messages and direct them to the endpoint manager
	switch msg := message.(type) {
	case *actor.Watch:
		rw := &remoteWatch{
			Watcher: msg.Watcher,
			Watchee: pid,
		}
		actor.Tell(endpointManagerPID, rw)
	case *actor.Unwatch:
		ruw := &remoteUnwatch{
			Watcher: msg.Watcher,
			Watchee: pid,
		}
		actor.Tell(endpointManagerPID, ruw)
	default:
		sendRemoteMessage(pid, message, nil)
	}
}

func (ref *remoteProcess) Stop(pid *actor.PID) {
	ref.SendSystemMessage(pid, stopMessage)
}

type remoteDeliver struct {
	message proto.Message
	target  *actor.PID
	sender  *actor.PID
}
