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

package client

import (
	"context"
)

// Service is the client service (no gRPC server, no discovery, no APM; Run blocks until Close).
type Service struct {
	runCtx    context.Context
	runCancel context.CancelFunc
}

// CreateService creates a new client service.
func CreateService(ctx context.Context) *Service {
	runCtx, runCancel := context.WithCancel(context.Background())
	return &Service{
		runCtx:    runCtx,
		runCancel: runCancel,
	}
}

// Run starts the client service and blocks until Close.
func (s *Service) Run() error {
	<-s.runCtx.Done()
	return nil
}

// Close stops the client service. Cancels runCtx so Run() returns.
func (s *Service) Close() error {
	if s.runCancel != nil {
		s.runCancel()
		s.runCancel = nil
	}
	return nil
}
