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

package reporter

import (
	"context"
	"fmt"
	"math/rand"
	"sync"
	"time"

	"os-artificer/saber/internal/agent/config"
	"os-artificer/saber/pkg/constant"
	"os-artificer/saber/pkg/logger"
	"os-artificer/saber/pkg/proto"
	"os-artificer/saber/pkg/tools"

	"google.golang.org/grpc"
	"google.golang.org/grpc/connectivity"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/keepalive"
)

var _ Reporter = (*TransferReporter)(nil)

func init() {
	RegisterReporter("transfer", func(ctx context.Context, opts any) (Reporter, error) {
		o, ok := opts.(*config.ReporterOpts)
		if !ok {
			return nil, fmt.Errorf("transfer reporter expects *config.ReporterOpts, got %T", opts)
		}
		endpoints, _ := o.Config["endpoints"].(string)
		if endpoints == "" {
			return nil, fmt.Errorf("transfer reporter config missing or invalid endpoints")
		}

		clientID, err := tools.MachineID("saber-agent")
		if err != nil {
			return nil, fmt.Errorf("failed to generate machine-id: %v", err)
		}

		return NewTransferReporter(ctx, endpoints, clientID)
	})
}

// TransferReporter is the reporter for transfer.
type TransferReporter struct {
	conn                 *grpc.ClientConn
	client               proto.TransferServiceClient
	stream               proto.TransferService_PushDataClient
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

// NewTransferReporter creates a new transfer reporter.
func NewTransferReporter(ctx context.Context, serverAddr string, clientId string) (*TransferReporter, error) {
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

	return &TransferReporter{
		conn:                 conn,
		client:               proto.NewTransferServiceClient(conn),
		ctx:                  ctxCancel,
		cancel:               cancel,
		clientId:             clientId,
		reconnectInterval:    constant.DefaultClientReconnectInterval,
		maxReconnectAttempts: constant.DefaultClientMaxReconnectAttempts,
	}, nil
}

func (c *TransferReporter) connect() error {
	c.mu.Lock()
	if c.closed {
		c.mu.Unlock()
		return fmt.Errorf("client is closed")
	}
	c.mu.Unlock()

	stream, err := c.client.PushData(c.ctx)
	if err != nil {
		return err
	}

	c.mu.Lock()
	c.stream = stream
	c.reconnectAttempts = 0
	c.mu.Unlock()

	go c.sendConnectionEstablished()

	go c.monitorConnection()

	return nil
}

func (c *TransferReporter) handleDisconnect() {
	c.mu.Lock()
	if c.closed || c.reconnecting {
		c.mu.Unlock()
		return
	}

	c.reconnecting = true
	c.stream = nil // so SendMessage returns "stream not initialized" until reconnect succeeds
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
		logger.Warnf("Max reconnect attempts (%d) reached, giving up", maxAttempts)
		return
	}

	// The exponential backoff algorithm calculates the reconnection interval.
	backoffInterval := c.reconnectInterval * time.Duration(1<<uint(reconnectAttempts-1))

	// Add some randomness to avoid the stampede effect.
	backoffInterval = backoffInterval + time.Duration(rand.Int63n(int64(backoffInterval/2)))

	logger.Infof("Reconnect attempt %d in %v", reconnectAttempts, backoffInterval)
	time.Sleep(backoffInterval)

	// Do not attempt connect if client was closed during backoff
	c.mu.Lock()
	closed := c.closed
	c.mu.Unlock()
	if closed {
		return
	}

	// retry
	logger.Infof("Attempting to reconnect...")
	err := c.connect()
	if err != nil {
		logger.Warnf("Reconnect failed: %v", err)
		go c.handleDisconnect()
	} else {
		logger.Infof("Reconnect successful")
	}
}

func (c *TransferReporter) monitorConnection() {
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
				logger.Warnf("Connection state: %s, starting reconnect", state.String())
				go c.handleDisconnect()
				return
			}

		case <-c.ctx.Done():
			return
		}
	}
}

func (c *TransferReporter) sendConnectionEstablished() {
	msg := &proto.TransferRequest{}

	c.mu.RLock()
	stream := c.stream
	c.mu.RUnlock()

	if stream == nil {
		return
	}

	if err := stream.Send(msg); err != nil {
		logger.Errorf("Error sending connection established message: %v", err)
	}
}

func (c *TransferReporter) Run() error {
	err := c.connect()
	if err != nil {
		logger.Errorf("failed to connect remote server. errmsg:%v", err)
		return err
	}

	<-c.ctx.Done()

	return c.ctx.Err()
}

func (c *TransferReporter) SendMessage(ctx context.Context, content []byte) error {
	c.mu.RLock()
	closed := c.closed
	stream := c.stream
	clientId := c.clientId
	c.mu.RUnlock()

	if closed {
		return fmt.Errorf("client is closed")
	}

	if stream == nil {
		return fmt.Errorf("stream not initialized")
	}

	msg := &proto.TransferRequest{
		ClientID: clientId,
		Payload:  content,
	}

	err := stream.Send(msg)
	if err != nil {
		// Trigger reconnect on send failure so we recover without waiting for monitor tick
		c.mu.Lock()
		if !c.closed && !c.reconnecting {
			c.stream = nil
			go c.handleDisconnect()
		}
		c.mu.Unlock()
		return err
	}
	return nil
}

func (c *TransferReporter) Close() error {
	c.mu.Lock()
	defer c.mu.Unlock()

	if c.closed {
		return nil
	}

	c.closed = true

	if c.cancel != nil {
		c.cancel()
	}

	if c.conn != nil {
		return c.conn.Close()
	}
	return nil
}

// GetConnectionState Get the state of the connection.
func (c *TransferReporter) GetConnectionState() connectivity.State {
	c.mu.RLock()
	defer c.mu.RUnlock()

	if c.conn == nil {
		return connectivity.Shutdown
	}

	return c.conn.GetState()
}
