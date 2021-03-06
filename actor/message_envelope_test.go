package actor

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestNormalMessageGivesEmptyMessageHeaders(t *testing.T) {
	props := FromFunc(func(ctx Context) {
		switch ctx.Message().(type) {
		case string:
			l := len(ctx.MessageHeader().Keys())
			ctx.Respond(l)
		}
	})
	a := Spawn(props)
	defer StopGraceful(a, 10*time.Second)
	f := RequestFuture(a, "hello", testTimeout)
	res := assertFutureSuccess(f, t).(int)
	assert.Equal(t, 0, res)
}
