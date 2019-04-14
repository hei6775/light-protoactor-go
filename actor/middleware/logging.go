package middleware

import (
	"log"
	"reflect"

	"github.com/hei6775/light-protoactor-go/actor"
)

// logger is message middleware which logs messages before continuing to the next middleware
func Logger(next actor.ActorFunc) actor.ActorFunc {
	fn := func(c actor.Context) {
		message := c.Message()
		log.Printf("%v got %v %+v", c.Self(), reflect.TypeOf(message), message)
		next(c)
	}

	return fn
}
