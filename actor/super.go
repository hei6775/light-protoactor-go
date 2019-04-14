package actor

import (
	"sync"
	"sync/atomic"
	"time"
)

type VFOnStarted func(ctx Context)
type VFOnStopping func(ctx Context)
type VFOnStopped func(ctx Context)
type VFOnReceiveMessage func(ctx Context)
type VFOnActorTerminated func(who *PID, ctx Context)
type VFOnRestarting func(ctx Context)
type VFOnRestarted func(ctx Context)

type Super struct {
	vfOnStarted         VFOnStarted
	vfOnStopping        VFOnStopping
	vfOnStopped         VFOnStopped
	vfOnReceiveMessage  VFOnReceiveMessage
	vfOnActorTerminated VFOnActorTerminated
	vfOnRestarting      VFOnRestarting
	vfOnRestarted       VFOnRestarted

	PID       *PID
	ParentPID *PID

	restarting bool

	startedWg     *sync.WaitGroup
	stoppedWg     *sync.WaitGroup
	stoppedWgOnce sync.Once

	timerMgr *TimerManager

	stopFlag int32
}

////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
func NewSuper(
	vfOnStarted VFOnStarted,
	vfOnStopping VFOnStopping,
	vfOnStopped VFOnStopped,
	vfOnReceiveMessage VFOnReceiveMessage,
	vfOnActorTerminated VFOnActorTerminated,
	vfOnRestarting VFOnRestarting,
	vfOnRestarted VFOnRestarted,
	startedWg *sync.WaitGroup,
	stoppedWg *sync.WaitGroup,
) *Super {

	s := &Super{
		vfOnStarted:         vfOnStarted,
		vfOnStopping:        vfOnStopping,
		vfOnStopped:         vfOnStopped,
		vfOnReceiveMessage:  vfOnReceiveMessage,
		vfOnActorTerminated: vfOnActorTerminated,
		vfOnRestarting:      vfOnRestarting,
		vfOnRestarted:       vfOnRestarted,
		startedWg:           startedWg,
		stoppedWg:           stoppedWg,
		stopFlag:            0,
	}

	if s.vfOnStarted == nil {
		s.vfOnStarted = func(ctx Context) {}
	}
	if s.vfOnStopping == nil {
		s.vfOnStopping = func(ctx Context) {}
	}
	if s.vfOnStopped == nil {
		s.vfOnStopped = func(ctx Context) {}
	}
	if s.vfOnReceiveMessage == nil {
		s.vfOnReceiveMessage = func(ctx Context) {}
	}
	if s.vfOnActorTerminated == nil {
		s.vfOnActorTerminated = func(who *PID, ctx Context) {}
	}
	if s.vfOnRestarting == nil {
		s.vfOnRestarting = func(ctx Context) {}
	}
	if s.vfOnRestarted == nil {
		s.vfOnRestarted = func(ctx Context) {}
	}

	if s.startedWg != nil {
		s.startedWg.Add(1)
	}

	return s
}

////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
func (s *Super) Receive(ctx Context) {
	switch msg := ctx.Message().(type) {
	case *Started:
		if !s.restarting {
			s.PID = ctx.Self()
			s.ParentPID = ctx.Parent()
			s.timerMgr = NewTimerManager(s.PID)
			s.vfOnStarted(ctx)
			if s.startedWg != nil {
				s.startedWg.Done()
			}
			if s.stoppedWg != nil {
				s.stoppedWg.Add(1)
			}
			logger.Debug("[%s] started", s.PID.Id)
		} else {
			logger.Error("[%s] restarted", s.PID.Id)
			s.restarting = false
			s.vfOnRestarted(ctx)
		}

	case *Stopping:
		logger.Debug("[%s] stopping", s.PID.Id)
		s.vfOnStopping(ctx)

	case *Stopped:
		logger.Debug("[%s] stopped", s.PID.Id)
		s.vfOnStopped(ctx)
		s.timerMgr.StopAll()
		if s.stoppedWg != nil {
			s.stoppedWgOnce.Do(func() {
				// 避免重复Done, 导致panic
				s.stoppedWg.Done()
			})
		}

	case *Terminated:
		s.vfOnActorTerminated(msg.Who, ctx)

	case *Restarting:
		logger.Error("[%s] restarting", s.PID.Id)
		s.restarting = true
		s.vfOnRestarting(ctx)

	case *timeout:
		if err := s.timerMgr.Trigger(msg.ID); err != nil {
			logger.Error("[%s] trigger time fail, %s", s.PID.Id, err)
		}

	default:
		s.vfOnReceiveMessage(ctx)
	}
}

////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
func (s Super) Sign() string {
	return s.PID.Id
}

////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
type ErrActorStopped struct {
}

func (ErrActorStopped) Error() string {
	return "actor stopped"
}

func (s *Super) IsStopped() bool {
	return atomic.LoadInt32(&s.stopFlag) != 0
}

func (s *Super) Stop() error {
	if !atomic.CompareAndSwapInt32(&s.stopFlag, 0, 1) {
		return &ErrActorStopped{}
	}
	if s.PID != nil {
		StopActor(s.PID)
	}
	return nil
}

////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
func (s *Super) NewTimer(dur time.Duration, cb TimerCallback) TimerID {
	return s.timerMgr.NewTimer(dur, cb)
}

func (s *Super) NewLoopTimer(interval time.Duration, cb TimerCallback) TimerID {
	return s.timerMgr.NewLoopTimer(interval, cb)
}

func (s *Super) StopTimer(id TimerID) error {
	return s.timerMgr.Stop(id)
}
