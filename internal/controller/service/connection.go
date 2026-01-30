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
	"sync"
	"time"

	"os-artificer/saber/pkg/logger"
	"os-artificer/saber/pkg/proto"
)

type Connection struct {
	ClientID   string
	Stream     proto.ControllerService_ConnectServer
	SendChan   chan *proto.ProbeResponse
	LastActive time.Time
	Metadata   map[string]string
	mu         sync.RWMutex
	closed     bool
}

func (c *Connection) sendMessages(ctx context.Context) {
	for {
		select {
		case <-ctx.Done():
			return

		case msg := <-c.SendChan:
			if c.isClosed() {
				return
			}

			c.updateLastActive()

			if err := c.Stream.Send(msg); err != nil {
				c.close()
				return
			}
		}
	}
}

func (c *Connection) receiveMessages(ctx context.Context) {
	for {
		select {
		case <-ctx.Done():
			return

		default:
			if c.isClosed() {
				return
			}

			msg, err := c.Stream.Recv()
			if err != nil {
				c.close()
				return
			}

			c.updateLastActive()

			logger.Debug("Received from %s: %v\n", c.ClientID, msg.GetPayload())
		}
	}
}

func (c *Connection) updateLastActive() {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.LastActive = time.Now()
}

func (c *Connection) isClosed() bool {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return c.closed
}

func (c *Connection) close() {
	c.mu.Lock()
	defer c.mu.Unlock()

	if !c.closed {
		c.closed = true
		close(c.SendChan)
	}
}
