package remote

import (
	"io/ioutil"
	"log"
	"net"
	"fmt"

	"google.golang.org/grpc"
	"google.golang.org/grpc/grpclog"

	"github.com/hei6775/light-protoactor-go/actor"
)

// Start the remote server
func Start(address string, options ...RemotingOption) error {
	grpclog.SetLogger(log.New(ioutil.Discard, "", 0))
	lis, err := net.Listen("tcp", address)
	if err != nil {
		//logger.Error("failed to listen, %s", err)
		return fmt.Errorf("failed to listen, %s", err)
	}
	config := defaultRemoteConfig()
	for _, option := range options {
		option(config)
	}

	address = lis.Addr().String()
	actor.ProcessRegistry.RegisterAddressResolver(remoteHandler)
	actor.ProcessRegistry.Address = address

	spawnActivatorActor()
	spawnEndpointManager(config)

	s := grpc.NewServer(config.serverOptions...)
	RegisterRemotingServer(s, &server{})
	logger.Info("Starting Proto.Actor server, address=[%v]", address)
	go s.Serve(lis)

	return nil
}
