package eventstream_test

import (
	"fmt"

	"gitee.com/lwj8507/light-protoactor-go/eventstream"
)

// Subscribe subscribes to events
func Example_Subscribe() {
	sub := eventstream.Subscribe(func(event interface{}) {
		fmt.Println(event)
	})

	// only allow strings
	sub.WithPredicate(func(evt interface{}) bool {
		_, ok := evt.(string)
		return ok
	})

	eventstream.Publish("Hello World")
	eventstream.Publish(1)

	eventstream.Unsubscribe(sub)

	// Output: Hello World
}
