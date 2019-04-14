package actor

import (
	"gitee.com/lwj8507/light-protoactor-go/eventstream"
)

//SupervisorEvent is sent on the EventStream when a supervisor have applied a directive to a failing child actor
type SupervisorEvent struct {
	Child     *PID
	Reason    interface{}
	Directive Directive
}

var (
	supervisionSubscriber *eventstream.Subscription
)

func init() {
	supervisionSubscriber = eventstream.Subscribe(func(evt interface{}) {
		if supervisorEvent, ok := evt.(*SupervisorEvent); ok {
			logger.Debug("[SUPERVISION] actor=[%v], directive=[%v], reason=[%v]", supervisorEvent.Child, supervisorEvent.Directive, supervisorEvent.Reason)
		}
	})
}
