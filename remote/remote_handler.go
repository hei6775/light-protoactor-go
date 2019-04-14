package remote

import "gitee.com/lwj8507/light-protoactor-go/actor"

func remoteHandler(pid *actor.PID) (actor.Process, bool) {
	ref := newRemoteProcess(pid)
	return ref, true
}
