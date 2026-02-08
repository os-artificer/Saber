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
package service

import (
	"context"
	"fmt"
	"net"
	"sync"
	"time"

	"os-artificer/saber/pkg/constant"
	"os-artificer/saber/pkg/logger"
	"os-artificer/saber/pkg/proto"

	"google.golang.org/grpc"
	"google.golang.org/grpc/keepalive"
)

type Service struct {
	proto.UnimplementedControllerServiceServer

	ctx         context.Context
	address     string
	serviceID   string
	connections map[string]*Connection
	mu          sync.RWMutex
}

func New(ctx context.Context, address string, serviceID string) *Service {
	return &Service{
		ctx:         ctx,
		address:     address,
		serviceID:   serviceID,
		connections: make(map[string]*Connection),
	}
}

// extractClientInfo is unused; clientID is read from first AgentRequest in Connect.
func (s *Service) extractClientInfo(ctx context.Context) (string, map[string]string, error) {
	clientID := fmt.Sprintf("client-%d", time.Now().UnixNano())
	metadata := make(map[string]string)
	_ = ctx
	return clientID, metadata, nil
}

func (s *Service) registerConnection(clientID string, conn *Connection) {
	s.mu.Lock()
	defer s.mu.Unlock()

	if existingConn, exists := s.connections[clientID]; exists {
		existingConn.close()
	}

	s.connections[clientID] = conn
}

func (s *Service) unregisterConnection(clientID string) {
	s.mu.Lock()
	defer s.mu.Unlock()

	if conn, exists := s.connections[clientID]; exists {
		conn.close()
		delete(s.connections, clientID)
	}
}

func (s *Service) getConnection(clientID string) (*Connection, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	conn, exists := s.connections[clientID]
	return conn, exists
}

func (s *Service) Connect(stream proto.ControllerService_ConnectServer) error {
	// Read first message to get clientID; same clientID reconnecting will close previous session.
	req, err := stream.Recv()
	if err != nil {
		return err
	}
	clientID := req.GetClientID()
	if clientID == "" {
		clientID = fmt.Sprintf("client-%d", time.Now().UnixNano())
	}
	metadata := make(map[string]string)
	if h := req.GetHeaders(); h != nil {
		for k, v := range h {
			metadata[k] = v
		}
	}

	conn := &Connection{
		ClientID:   clientID,
		Stream:     stream,
		SendChan:   make(chan *proto.AgentResponse, 100),
		LastActive: time.Now(),
		Metadata:   metadata,
		FirstReq:   req,
	}

	s.registerConnection(clientID, conn) // closes any existing connection with same clientID
	defer s.unregisterConnection(clientID)

	wg := &sync.WaitGroup{}

	wg.Add(2)

	go func() {
		defer wg.Done()
		conn.sendMessages(s.ctx)
	}()

	go func() {
		defer wg.Done()
		conn.receiveMessages(s.ctx)
	}()

	wg.Wait()
	return nil
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

	proto.RegisterControllerServiceServer(svr, s)
	lis, err := net.Listen("tcp", s.address)
	if err != nil {
		return err
	}

	logger.Infof("Server listening at %v", lis.Addr())
	return svr.Serve(lis)
}
