package remote

import "github.com/hei6775/light-protoactor-go/actor"

func remoteHandler(pid *actor.PID) (actor.Process, bool) {
	ref := newRemoteProcess(pid)
	return ref, true
}
