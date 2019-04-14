package remote

import (
	"github.com/hei6775/light-protoactor-go/actor"
)

type server struct{}

func (s *server) Receive(stream Remoting_ReceiveServer) error {
	for {
		batch, err := stream.Recv()
		if err != nil {
			//actor.Tell(endpointManagerPID, &EndpointReaderFailedToRead{Err:err})
			return err
		}
		for _, envelope := range batch.Envelopes {
			targetName := batch.TargetNames[envelope.Target]
			pid := actor.NewLocalPID(targetName)
			message := deserialize(envelope, batch.TypeNames[envelope.TypeId])
			//if message is system message send it as sysmsg instead of usermsg

			sender := envelope.Sender

			switch msg := message.(type) {
			case *actor.Terminated:
				rt := &remoteTerminate{
					Watchee: msg.Who,
					Watcher: pid,
				}
				actor.Tell(endpointManagerPID, rt)
			case actor.SystemMessage:
				ref, _ := actor.ProcessRegistry.GetLocal(pid.Id)
				ref.SendSystemMessage(pid, msg)
			default:
				actor.Request(pid, message, sender)
			}
		}
	}
}
