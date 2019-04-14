package router

import "gitee.com/lwj8507/light-protoactor-go/actor"

// A type that satisfies router.Interface can be used as a router
type Interface interface {
	RouteMessage(message interface{}, sender *actor.PID)
	SetRoutees(routees *actor.PIDSet)
	GetRoutees() *actor.PIDSet
}
