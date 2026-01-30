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

	"os-artificer/saber/pkg/constant"
	"os-artificer/saber/pkg/gerrors"
	"os-artificer/saber/pkg/logger"
	"os-artificer/saber/pkg/proto"

	"github.com/segmentio/kafka-go"
)

type requestEventC chan *proto.TransferRequest

// connectionHandler service connection handler
type connectionHandler struct {
	eventC      requestEventC
	quit        chan struct{}
	wg          sync.WaitGroup
	kafkaWriter *kafka.Writer
}

func (c *connectionHandler) readEvent() {
	for {
		select {
		case <-c.quit:
			return

		case msg := <-c.eventC:
			if msg != nil && len(msg.Payload) > 0 && c.kafkaWriter != nil {
				key := []byte(msg.ClientID)
				if err := c.kafkaWriter.WriteMessages(context.Background(), kafka.Message{
					Key:   key,
					Value: msg.Payload,
				}); err != nil {
					logger.Warn("write to kafka failed: %v", err)
				}
			}
		}
	}
}

func (c *connectionHandler) postEvent(event *proto.TransferRequest) error {
	select {
	case c.eventC <- event:
		return nil

	default:
		return gerrors.New(gerrors.QueueFull, "connection handler event channel is full")
	}
}

func (c *connectionHandler) run() {
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

func (c *connectionHandler) close() {
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
