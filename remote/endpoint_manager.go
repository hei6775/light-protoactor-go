package remote

import (
	"gitee.com/lwj8507/light-protoactor-go/actor"
	"gitee.com/lwj8507/light-protoactor-go/mailbox"
)

var endpointManagerPID *actor.PID

func newEndpointManager(config *remoteConfig) actor.Producer {
	return func() actor.Actor {
		return &endpointManager{
			config: config,
		}
	}
}

func spawnEndpointManager(config *remoteConfig) {
	props := actor.
		FromProducer(newEndpointManager(config)).
		WithMailbox(mailbox.Bounded(config.endpointManagerQueueSize)).
		WithSupervisor(actor.RestartingSupervisorStrategy())

	endpointManagerPID = actor.Spawn(props)
}

type endpoint struct {
	writer  *actor.PID
	watcher *actor.PID
}

type endpointManager struct {
	config    *remoteConfig
	endpoints map[string]*endpoint
}

func (mgr *endpointManager) Receive(ctx actor.Context) {
	switch msg := ctx.Message().(type) {
	case *actor.Started:
		mgr.endpoints = make(map[string]*endpoint)
		logger.Info("Started EndpointManager")

	case *actor.Stopped:
		logger.Info("Stopped EndpointManager")

	case *EndpointTerminatedEvent:
		//fmt.Println("*EndpointTerminatedEvent")
		logger.Warn("EndpointWriter lost connection to address, address=[%v]", msg.Address)
		address := msg.Address
		//endpoint := mgr.ensureConnected(address, ctx)
		//ctx.Tell(endpoint.watcher, msg)
		endpoint, ok := mgr.getEndpoint(address)
		if ok {
			ctx.Tell(endpoint.watcher, msg)
			actor.StopActor(endpoint.writer)
			actor.StopActor(endpoint.watcher)
			delete(mgr.endpoints, address)
		}

	case *EndpointReaderFailedToRead:
		logger.Error("EndpointReader failed to read, %v", msg.Err)

	case *remoteTerminate:
		//fmt.Println("*remoteTerminate")
		address := msg.Watchee.Address
		//endpoint := mgr.ensureConnected(address, ctx)
		//ctx.Tell(endpoint.watcher, msg)
		endpoint, ok := mgr.getEndpoint(address)
		if ok {
			ctx.Tell(endpoint.watcher, msg)
		}

	case *remoteWatch:
		//fmt.Println("*remoteWatch")
		address := msg.Watchee.Address
		//endpoint := mgr.ensureConnected(address, ctx)
		//ctx.Tell(endpoint.watcher, msg)
		endpoint, ok := mgr.getEndpoint(address)
		if ok {
			ctx.Tell(endpoint.watcher, msg)
		} else {
			ctx.Tell(msg.Watcher, &actor.Terminated{
				Who:               msg.Watchee,
				AddressTerminated: true,
			})
		}

	case *remoteUnwatch:
		//fmt.Println("*remoteUnwatch")
		address := msg.Watchee.Address
		//endpoint := mgr.ensureConnected(address, ctx)
		//ctx.Tell(endpoint.watcher, msg)
		endpoint, ok := mgr.getEndpoint(address)
		if ok {
			ctx.Tell(endpoint.watcher, msg)
		}

	case *remoteDeliver:
		//fmt.Println("*remoteDeliver")
		address := msg.target.Address
		endpoint := mgr.ensureConnected(address, ctx)
		ctx.Tell(endpoint.writer, msg)
	}
}

func (mgr *endpointManager) getEndpoint(address string) (e *endpoint, ok bool) {
	e, ok = mgr.endpoints[address]
	return
}

func (mgr *endpointManager) spawnEndpointWriter(address string, ctx actor.Context) *actor.PID {
	props := actor.
		FromProducer(newEndpointWriter(address, mgr.config)).
		WithMailbox(newEndpointWriterMailbox(mgr.config.endpointWriterBatchSize, mgr.config.endpointWriterQueueSize))
	pid := ctx.Spawn(props)
	return pid
}

func (mgr *endpointManager) spawnEndpointWatcher(address string, ctx actor.Context) *actor.PID {
	props := actor.
		FromProducer(newEndpointWatcher(address))
	pid := ctx.Spawn(props)
	return pid
}

func (mgr *endpointManager) newEndpoint(address string, ctx actor.Context) *endpoint {
	e := &endpoint{
		writer:  mgr.spawnEndpointWriter(address, ctx),
		watcher: mgr.spawnEndpointWatcher(address, ctx),
	}
	return e
}

func (mgr *endpointManager) ensureConnected(address string, ctx actor.Context) *endpoint {
	e, ok := mgr.getEndpoint(address)
	if !ok {
		e = mgr.newEndpoint(address, ctx)
		mgr.endpoints[address] = e
	}
	return e
}
