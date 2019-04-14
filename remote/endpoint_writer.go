package remote

import (
	"github.com/hei6775/light-protoactor-go/actor"

	"golang.org/x/net/context"
	"google.golang.org/grpc"
)

func newEndpointWriter(address string, config *remoteConfig) actor.Producer {
	return func() actor.Actor {
		return &endpointWriter{
			address: address,
			config:  config,
		}
	}
}

type endpointWriter struct {
	config  *remoteConfig
	address string
	conn    *grpc.ClientConn
	stream  Remoting_ReceiveClient
}

func (state *endpointWriter) initialize() {
	err := state.initializeInternal()
	if err != nil {
		logger.Error("EndpointWriter failed to connect, address=[%v], %s", state.address, err)
	}
}

func (state *endpointWriter) initializeInternal() error {
	logger.Info("Started EndpointWriter, address=[%v]", state.address)
	logger.Info("EndpointWatcher connecting, address=[%v]", state.address)
	conn, err := grpc.Dial(state.address, state.config.dialOptions...)
	if err != nil {
		return err
	}
	state.conn = conn
	c := NewRemotingClient(conn)
	//	log.Printf("Getting stream from address %v", mgr.address)
	stream, err := c.Receive(context.Background(), state.config.callOptions...)
	if err != nil {
		return err
	}
	go func() {
		_, err := stream.Recv()
		if err != nil {
			//logger.Info("EndpointWriter lost connection to address, address=[%v]", state.address)
			actor.Tell(endpointManagerPID, &EndpointTerminatedEvent{Address: state.address})
		}
	}()
	logger.Info("EndpointWriter connected, address=[%v]", state.address)
	state.stream = stream
	return nil
}

func (state *endpointWriter) sendEnvelopes(msg []interface{}, ctx actor.Context) {
	envelopes := make([]*MessageEnvelope, len(msg))

	//type name uniqueness map name string to type index
	typeNames := make(map[string]int32)
	typeNamesArr := make([]string, 0)
	targetNames := make(map[string]int32)
	targetNamesArr := make([]string, 0)
	var typeID int32
	var targetID int32
	for i, tmp := range msg {
		rd := tmp.(*remoteDeliver)
		bytes, typeName, _ := serialize(rd.message)
		typeID, typeNamesArr = addToLookup(typeNames, typeName, typeNamesArr)
		targetID, targetNamesArr = addToLookup(targetNames, rd.target.Id, targetNamesArr)

		envelopes[i] = &MessageEnvelope{
			MessageData: bytes,
			Sender:      rd.sender,
			Target:      targetID,
			TypeId:      typeID,
		}
	}

	batch := &MessageBatch{
		TypeNames:   typeNamesArr,
		TargetNames: targetNamesArr,
		Envelopes:   envelopes,
	}
	err := state.stream.Send(batch)
	if err != nil {
		//ctx.Stash()
		logger.Debug("gRPC Failed to send, address=[%v]", state.address)
		//panic("restart it")
		actor.Tell(endpointManagerPID, &EndpointTerminatedEvent{Address: state.address})
	}
}

func addToLookup(m map[string]int32, name string, a []string) (int32, []string) {
	max := int32(len(m))
	id, ok := m[name]
	if !ok {
		m[name] = max
		id = max
		a = append(a, name)
	}
	return id, a
}

func (state *endpointWriter) Receive(ctx actor.Context) {
	switch msg := ctx.Message().(type) {
	case *actor.Started:
		state.initialize()

	case *actor.Stopping:
		// pass

	case *actor.Stopped:
		state.conn.Close()
		logger.Info("Stopped EndpointWriter, address=[%v]", state.address)

	case *actor.Restarting:
		state.conn.Close()

	case []interface{}:
		state.sendEnvelopes(msg, ctx)

	default:
		logger.Error("Unknown message[%#v]", msg, msg)
	}
}
