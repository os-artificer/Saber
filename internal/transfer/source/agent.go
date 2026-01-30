/**
 * Copyright 2025 Saber authors.
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
**/

package source

import (
	"context"
	"io"
	"net"

	"os-artificer/saber/pkg/constant"
	"os-artificer/saber/pkg/logger"
	"os-artificer/saber/pkg/proto"

	"google.golang.org/grpc"
	"google.golang.org/grpc/keepalive"
	"google.golang.org/grpc/peer"
)

var _ Source = (*AgentSource)(nil)

// AgentSource implements Source by accepting gRPC PushData streams from agents
// and delivering each TransferRequest to the Handler.
type AgentSource struct {
	address string
}

// NewAgentSource returns a Source that listens on address and serves TransferService PushData.
func NewAgentSource(address string) *AgentSource {
	return &AgentSource{address: address}
}

// Run starts the gRPC server and blocks until ctx is done or server stops.
// Each received TransferRequest is passed to h.OnTransferRequest.
func (p *AgentSource) Run(ctx context.Context, h Handler) error {
	kasp := keepalive.ServerParameters{
		Time:    constant.DefaultServerPingTime,
		Timeout: constant.DefaultPingTimeout,
	}
	kacp := keepalive.EnforcementPolicy{
		MinTime:             constant.DefaultKeepaliveMiniTime,
		PermitWithoutStream: true,
	}

	svr := grpc.NewServer(
		grpc.KeepaliveParams(kasp),
		grpc.KeepaliveEnforcementPolicy(kacp),
		grpc.MaxRecvMsgSize(constant.DefaultMaxReceiveMessageSize),
		grpc.MaxSendMsgSize(constant.DefaultMaxSendMessageSize),
	)

	pushSrv := &pushDataServer{handler: h}
	proto.RegisterTransferServiceServer(svr, pushSrv)

	lis, err := net.Listen("tcp", p.address)
	if err != nil {
		return err
	}

	logger.Info("Server listening at %v", lis.Addr())

	go func() {
		<-ctx.Done()
		svr.GracefulStop()
	}()

	return svr.Serve(lis)
}

// pushDataServer implements proto.TransferServiceServer and forwards each request to Handler.
type pushDataServer struct {
	proto.UnimplementedTransferServiceServer
	handler Handler
}

func (s *pushDataServer) PushData(stream proto.TransferService_PushDataServer) error {
	ctx := stream.Context()
	addr, ok := peer.FromContext(ctx)
	clientID := ""
	if ok {
		clientID = addr.String()
	}

	for {
		select {
		case <-ctx.Done():
			logger.Error("service exited due to canceled context, client-id: %s", clientID)
			return nil
		default:
			req, err := stream.Recv()
			if err == io.EOF {
				logger.Error("receiver server exited, client-id(%s), errmsg: %v", clientID, err)
				return nil
			}
			if err != nil {
				logger.Error("receiver server exited, client-id(%s), errmsg: %v", clientID, err)
				return nil
			}
			if err := s.handler.OnTransferRequest(req); err != nil {
				logger.Warn("handle the client event data failed, client-id: %s, errmsg: %v", clientID, err)
			}
		}
	}
}
