package actor

import (
	"fmt"
	"sync/atomic"
	"time"
)

type TimerID = uint64

type timeout struct {
	ID TimerID
}

type TimerCallback func(TimerID)

type timer struct {
	id     TimerID
	cb     TimerCallback
	looped bool
	stopCh chan struct{}
}

func (t *timer) stop() {
	close(t.stopCh)
}

func (t *timer) trigger() {
	t.cb(t.id)
}

type TimerManager struct {
	idGenerator uint64
	timers      map[TimerID]*timer
	pid         *PID
}

func NewTimerManager(pid *PID) *TimerManager {
	return &TimerManager{
		idGenerator: 0,
		timers:      make(map[TimerID]*timer),
		pid:         pid,
	}
}

func (m *TimerManager) nextID() TimerID {
	atomic.AddUint64(&m.idGenerator, 1)
	return TimerID(m.idGenerator)
}

func (m *TimerManager) NewTimer(dur time.Duration, callback TimerCallback) TimerID {
	newTimer := &timer{
		id:     m.nextID(),
		cb:     callback,
		looped: false,
		stopCh: make(chan struct{}),
	}

	m.timers[newTimer.id] = newTimer

	go func() {
		select {
		case <-time.After(dur):
			Tell(m.pid, &timeout{newTimer.id})

		case <-newTimer.stopCh:
			return
		}
	}()

	return newTimer.id
}

func (m *TimerManager) NewLoopTimer(dur time.Duration, callback TimerCallback) TimerID {
	newTimer := &timer{
		id:     m.nextID(),
		cb:     callback,
		looped: true,
		stopCh: make(chan struct{}),
	}

	m.timers[newTimer.id] = newTimer

	go func() {
		ticker := time.NewTicker(dur)
		chTicker := ticker.C
		msg := &timeout{newTimer.id}

		for {
			select {
			case <-chTicker:
				Tell(m.pid, msg)

			case <-newTimer.stopCh:
				ticker.Stop()
				return
			}
		}
	}()

	return newTimer.id
}

func (m *TimerManager) Trigger(id TimerID) error {
	if t, ok := m.timers[id]; ok {
		t.trigger()
		if !t.looped {
			delete(m.timers, id)
		}
		return nil
	}
	return fmt.Errorf("timer[%d] not found", id)
}

func (m *TimerManager) Stop(id TimerID) error {
	if t, ok := m.timers[id]; ok {
		t.stop()
		delete(m.timers, id)
		return nil
	}
	return fmt.Errorf("timer[%d] not found", id)
}

func (m *TimerManager) StopAll() {
	for _, t := range m.timers {
		t.stop()
	}
	// 清空map
	m.timers = make(map[TimerID]*timer)
}
