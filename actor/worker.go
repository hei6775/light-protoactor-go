package actor

import (
	"fmt"
	"sync"
	"sync/atomic"
	"time"
)

////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
type IWork interface {
	Do(inputs ...interface{}) ([]interface{}, error)
}

type WorkFunc func(inputs ...interface{}) ([]interface{}, error)

func (f WorkFunc) Do(inputs ...interface{}) ([]interface{}, error) {
	return f(inputs...)
}

type WorkCallback func(err error, outputs ...interface{})

type WorkRequest struct {
	pid    *PID
	iWork  IWork
	inputs []interface{}
	cb     WorkCallback
}

type WorkResult struct {
	outputs []interface{}
	err     error
}

func newWorkRequest(pid *PID, iWork IWork, inputs ...interface{}) *WorkRequest {
	wr := &WorkRequest{
		pid:    pid,
		iWork:  iWork,
		inputs: inputs,
	}
	return wr
}

func (wr *WorkRequest) WaitForResult(timeout time.Duration) ([]interface{}, error) {
	future := RequestFuture(wr.pid, wr, timeout)
	iResult, err := future.Result()
	if err != nil {
		return nil, err
	}
	result := iResult.(*WorkResult)
	return result.outputs, result.err
}

func (wr *WorkRequest) Post(cb WorkCallback) {
	wr.cb = cb
	Tell(wr.pid, wr)
}

////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
type Worker struct {
	*Super

	workNum  int64
	stopFlag int32
}

func NewWorker(
	vfOnStarted VFOnStarted,
	vfOnStopping VFOnStopping,
	vfOnStopped VFOnStopped,
	vfOnReceiveMessage VFOnReceiveMessage,
	vfOnActorTerminated VFOnActorTerminated,
	vfOnRestarting VFOnRestarting,
	vfOnRestarted VFOnRestarted,
	startedWg *sync.WaitGroup,
	stoppedWg *sync.WaitGroup,
) *Worker {
	worker := &Worker{
		Super: NewSuper(vfOnStarted, vfOnStopping, vfOnStopped, vfOnReceiveMessage,
			vfOnActorTerminated, vfOnRestarting, vfOnRestarted, startedWg, stoppedWg),
	}
	return worker
}

func (w *Worker) NewRequest(iWork IWork, inputs ...interface{}) (*WorkRequest, error) {
	if iWork == nil {
		return nil, fmt.Errorf("can not do a nil work")
	}

	if atomic.LoadInt32(&w.stopFlag) == 1 {
		// 退出中, 不再接受新操作
		return nil, fmt.Errorf("in stopping state, can not do a new work")
	}

	atomic.AddInt64(&w.workNum, 1)

	return newWorkRequest(w.PID, iWork, inputs...), nil
}

func (w *Worker) Receive(ctx Context) {
	w.Super.Receive(ctx)

	switch msg := ctx.Message().(type) {
	case *WorkRequest:
		defer atomic.AddInt64(&w.workNum, -1)

		if msg.iWork == nil {
			logger.Error("[%s] can not do a nil iWork", w.PID.Id)
			return
		}
		result := &WorkResult{}
		result.outputs, result.err = msg.iWork.Do(msg.inputs...)
		if msg.cb != nil {
			msg.cb(result.err, result.outputs...)
		}
		sender := ctx.Sender()
		if sender != nil {
			Tell(sender, result)
		}
	}
}

const (
	checkToStopInterval = 3 * time.Second // 检测是否可退出的间隔
)

func (w *Worker) WorkNum() int64 {
	return atomic.LoadInt64(&w.workNum)
}

func (w *Worker) IsStopping() bool {
	return atomic.LoadInt32(&w.stopFlag) == 1
}

func (w *Worker) Stop() {
	if !atomic.CompareAndSwapInt32(&w.stopFlag, 0, 1) {
		logger.Error("[%s] has in stopping state", w.PID.Id)
		return
	}

	logger.Info("[%s] enter stopping state", w.PID.Id)

	if atomic.LoadInt64(&w.workNum) == 0 {
		logger.Info("[%s] all work has done, stop", w.PID.Id)
		StopActor(w.PID)
		return
	}

	// 每隔一段时间检测是否可以退出
	w.NewTimer(checkToStopInterval, w.checkToStop)
}

func (w *Worker) checkToStop(id TimerID) {
	logger.Info("[%s] check left work num to stop", w.PID.Id)

	if atomic.LoadInt64(&w.workNum) != 0 {
		w.NewTimer(checkToStopInterval, w.checkToStop)
		return
	}

	logger.Info("[%s] all work has done, stop", w.PID.Id)

	StopActor(w.PID)
}
