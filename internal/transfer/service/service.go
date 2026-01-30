/**
 * Copyright 2025 saber authors.
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

package service

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

type Service struct {
	proto.UnimplementedTransferServiceServer

	ctx       context.Context
	address   string
	serviceID string
}

func New(ctx context.Context, address string, serviceID string) *Service {
	return &Service{
		ctx:       ctx,
		address:   address,
		serviceID: serviceID,
	}
}

func (s *Service) PushData(stream proto.TransferService_PushDataServer) error {
	ctx := stream.Context()
	addr, ok := peer.FromContext(ctx)

	clientID := ""
	if ok {
		clientID = addr.String()
	}

	connHandler := &connectionHandler{
		eventC: make(requestEventC, constant.DefaultMaxReceiveMessageSize),
		quit:   make(chan struct{}),
	}

	connHandler.run()
	defer connHandler.close()

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

			if err := connHandler.postEvent(req); err != nil {
				logger.Warn("handle the client event data failed, client-id: %s, errmsg: %v", clientID, err)
			}
		}
	}
}

func (s *Service) Run() error {
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

	proto.RegisterTransferServiceServer(svr, s)
	lis, err := net.Listen("tcp", s.address)
	if err != nil {
		return err
	}

	logger.Info("Server listening at %v", lis.Addr())
	return svr.Serve(lis)
}
