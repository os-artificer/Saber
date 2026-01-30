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
	"sync"

	"os-artificer/saber/internal/transfer/sink"
	"os-artificer/saber/pkg/constant"
	"os-artificer/saber/pkg/gerrors"
	"os-artificer/saber/pkg/logger"
	"os-artificer/saber/pkg/proto"
)

var _ Handler = (*ConnectionHandler)(nil)

type requestEventC chan *proto.TransferRequest

// ConnectionHandler buffers TransferRequests and writes them to a Sink.
// It implements Handler for use with Source (e.g. AgentSource).
type ConnectionHandler struct {
	eventC requestEventC
	quit   chan struct{}
	wg     sync.WaitGroup
	sink   sink.Sink
}

// NewConnectionHandler returns a Handler that buffers requests and writes to the given Sink.
// If sink is nil, received requests are dropped (no write).
func NewConnectionHandler(sink sink.Sink) *ConnectionHandler {
	c := &ConnectionHandler{
		eventC: make(requestEventC, constant.DefaultMaxReceiveMessageSize),
		quit:   make(chan struct{}, 1),
		sink:   sink,
	}
	c.run()
	return c
}

// OnTransferRequest implements Handler.
func (c *ConnectionHandler) OnTransferRequest(req *proto.TransferRequest) error {
	return c.postEvent(req)
}

func (c *ConnectionHandler) readEvent() {
	for {
		select {
		case <-c.quit:
			return

		case msg := <-c.eventC:
			if msg != nil && c.sink != nil {
				if err := c.sink.Write(context.Background(), msg); err != nil {
					logger.Warn("sink write failed: %v", err)
				}
			}
		}
	}
}

func (c *ConnectionHandler) postEvent(event *proto.TransferRequest) error {
	select {
	case c.eventC <- event:
		return nil

	default:
		return gerrors.New(gerrors.QueueFull, "connection handler event channel is full")
	}
}

func (c *ConnectionHandler) run() {
	if c.eventC == nil {
		c.eventC = make(chan *proto.TransferRequest, constant.DefaultMaxReceiveMessageSize)
	}
	if c.quit == nil {
		c.quit = make(chan struct{}, 1)
	}
	c.wg.Add(1)
	go func() {
		c.readEvent()
		c.wg.Done()
	}()
}

// Close stops the handler and waits for in-flight events.
func (c *ConnectionHandler) Close() {
	c.close()
}

func (c *ConnectionHandler) close() {
	if c.quit != nil {
		close(c.quit)
	}

	if c.eventC != nil {
		close(c.eventC)
	}

	c.wg.Wait()
	c.quit = nil
	c.eventC = nil
}
