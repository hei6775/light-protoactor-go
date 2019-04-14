package actor

import (
	"reflect"
	"testing"

	"time"

	"github.com/stretchr/testify/assert"
)

type ShortLivingActor struct {
}

func (sla *ShortLivingActor) Receive(ctx Context) {

}

func TestStopFuture(t *testing.T) {
	logger.Debug("hello world")

	ID := "UniqueID"
	{
		props := FromInstance(&ShortLivingActor{})
		a, _ := SpawnNamed(props, ID)

		fut := StopFuture(a, 10*time.Second)

		res, errR := fut.Result()
		if errR != nil {
			assert.Fail(t, "Failed to wait stop actor %s", errR)
			return
		}

		_, ok := res.(*Terminated)
		if !ok {
			assert.Fail(t, "Cannot cast %s", reflect.TypeOf(res))
			return
		}

		_, found := ProcessRegistry.Get(a)
		assert.False(t, found)
	}
}
