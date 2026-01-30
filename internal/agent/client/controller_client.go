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

package client

import (
	"context"
	"fmt"
	"math/rand"
	"sync"
	"time"

	"os-artificer/saber/pkg/constant"
	"os-artificer/saber/pkg/logger"
	"os-artificer/saber/pkg/proto"

	"google.golang.org/grpc"
	"google.golang.org/grpc/connectivity"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/keepalive"
)

type ControllerClient struct {
	conn                 *grpc.ClientConn
	client               proto.ControllerServiceClient
	stream               proto.ControllerService_ConnectClient
	ctx                  context.Context
	cancel               context.CancelFunc
	clientId             string
	mu                   sync.RWMutex
	closed               bool
	reconnecting         bool
	reconnectInterval    time.Duration
	maxReconnectAttempts int
	reconnectAttempts    int
}

func NewControllerClient(ctx context.Context, serverAddr string, clientId string) (*ControllerClient, error) {

	kacp := keepalive.ClientParameters{
		Time:                constant.DefaultClientPingTime,
		Timeout:             constant.DefaultPingTimeout,
		PermitWithoutStream: true,
	}

	conn, err := grpc.NewClient(
		serverAddr,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithKeepaliveParams(kacp),
		grpc.WithDefaultCallOptions(
			grpc.MaxCallRecvMsgSize(constant.DefaultMaxReceiveMessageSize),
			grpc.MaxCallSendMsgSize(constant.DefaultMaxSendMessageSize),
		),
	)

	if err != nil {
		return nil, err
	}

	ctxCancel, cancel := context.WithCancel(ctx)

	return &ControllerClient{
		conn:                 conn,
		client:               proto.NewControllerServiceClient(conn),
		ctx:                  ctxCancel,
		cancel:               cancel,
		clientId:             clientId,
		reconnectInterval:    constant.DefaultClientReconnectInterval,
		maxReconnectAttempts: constant.DefaultClientMaxReconnectAttempts,
	}, nil
}

func (c *ControllerClient) connect() error {
	c.mu.Lock()
	if c.closed {
		c.mu.Unlock()
		return fmt.Errorf("client is closed")
	}
	c.mu.Unlock()

	stream, err := c.client.Connect(c.ctx)
	if err != nil {
		return err
	}

	c.mu.Lock()
	c.stream = stream
	c.reconnectAttempts = 0
	c.mu.Unlock()

	go c.receiveMessages()

	go c.sendConnectionEstablished()

	go c.monitorConnection()

	return nil
}

func (c *ControllerClient) handleDisconnect() {
	c.mu.Lock()
	if c.closed || c.reconnecting {
		c.mu.Unlock()
		return
	}

	c.reconnecting = true
	c.mu.Unlock()

	defer func() {
		c.mu.Lock()
		c.reconnecting = false
		c.mu.Unlock()
	}()

	c.mu.Lock()
	c.reconnectAttempts++
	reconnectAttempts := c.reconnectAttempts
	maxAttempts := c.maxReconnectAttempts
	c.mu.Unlock()

	if maxAttempts > 0 && reconnectAttempts > maxAttempts {
		logger.Warn("Max reconnect attempts (%d) reached, giving up", maxAttempts)
		return
	}

	// The exponential backoff algorithm calculates the reconnection interval.
	backoffInterval := c.reconnectInterval * time.Duration(1<<uint(reconnectAttempts-1))

	// Add some randomness to avoid the stampede effect.
	backoffInterval = backoffInterval + time.Duration(rand.Int63n(int64(backoffInterval/2)))

	logger.Info("Reconnect attempt %d in %v", reconnectAttempts, backoffInterval)
	time.Sleep(backoffInterval)

	// retry
	logger.Info("Attempting to reconnect...")
	err := c.connect()
	if err != nil {
		logger.Warn("Reconnect failed: %v", err)
		go c.handleDisconnect()
	} else {
		logger.Info("Reconnect successful")
	}
}

func (c *ControllerClient) monitorConnection() {
	ticker := time.NewTicker(5 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			c.mu.RLock()
			if c.closed {
				c.mu.RUnlock()
				return
			}
			c.mu.RUnlock()
			state := c.GetConnectionState()
			if state == connectivity.TransientFailure || state == connectivity.Shutdown {
				logger.Warn("Connection state: %s, starting reconnect", state.String())
				go c.handleDisconnect()
				return
			}

		case <-c.ctx.Done():
			return
		}
	}
}

func (c *ControllerClient) receiveMessages() {
	for {
		select {
		case <-c.ctx.Done():
			return

		default:
			c.mu.RLock()
			if c.closed || c.stream == nil {
				c.mu.RUnlock()
				return
			}
			stream := c.stream
			c.mu.RUnlock()

			msg, err := stream.Recv()
			if err != nil {
				logger.Warn("Error receiving message: %v, starting reconnect", err)
				go c.handleDisconnect()
				return
			}

			logger.Debug("msg:%v", msg)
		}
	}
}

func (c *ControllerClient) sendConnectionEstablished() {
	msg := &proto.AgentRequest{
		ClientID: c.clientId,
	}

	c.mu.RLock()
	defer c.mu.RUnlock()

	if c.stream == nil {
		return
	}

	if err := c.stream.Send(msg); err != nil {
		logger.Error("Error sending connection established message: %v", err)
	}
}

func (c *ControllerClient) Run() error {
	err := c.connect()
	if err != nil {
		logger.Error("failed to connect remote server. errmsg:%v", err)
		return err
	}

outerLoop:
	for {
		select {
		case <-c.ctx.Done():
			break outerLoop
		}
	}

	return nil
}

func (c *ControllerClient) SendMessage(content string) error {
	c.mu.RLock()
	defer c.mu.RUnlock()

	if c.closed {
		return fmt.Errorf("client is closed")
	}

	if c.stream == nil {
		return fmt.Errorf("stream not initialized")
	}

	msg := &proto.AgentRequest{}

	return c.stream.Send(msg)
}

func (c *ControllerClient) Close() {
	c.mu.Lock()
	defer c.mu.Unlock()

	if c.closed {
		return
	}

	c.closed = true

	if c.cancel != nil {
		c.cancel()
	}

	if c.conn != nil {
		c.conn.Close()
	}
}

// GetConnectionState Get the state of the connection.
func (c *ControllerClient) GetConnectionState() connectivity.State {
	c.mu.RLock()
	defer c.mu.RUnlock()

	if c.conn == nil {
		return connectivity.Shutdown
	}

	return c.conn.GetState()
}
