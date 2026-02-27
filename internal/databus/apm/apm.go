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

package apm

import (
	"net"
	"sync"

	pkgapm "os-artificer/saber/pkg/apm"
	"os-artificer/saber/pkg/logger"
	"os-artificer/saber/pkg/sbnet"
)

// APM APM service
type APM struct {
	enabled  bool
	endpoint sbnet.Endpoint

	mu       sync.Mutex
	listener net.Listener
	server   *sbnet.Server
}

// NewAPM creates a new APM service
func NewAPM(enabled bool, endpoint sbnet.Endpoint) *APM {
	return &APM{
		enabled:  enabled,
		endpoint: endpoint,
	}
}

// Run starts the APM service and blocks until Close() is called.
// Call it in a goroutine so the main process is not blocked.
func (a *APM) Run() error {
	if !a.enabled {
		return nil
	}

	srv := sbnet.NewServer(sbnet.WithRoutes(pkgapm.MetricsRoute()))
	lis, err := net.Listen(a.endpoint.Protocol, a.endpoint.HostPort())
	if err != nil {
		return err
	}

	a.mu.Lock()
	a.server = srv
	a.listener = lis
	a.mu.Unlock()

	logger.Infof("APM metrics server listening at %s", a.endpoint.String())
	return srv.Engine().RunListener(lis)
}

// Close closes the APM service and stops the metrics HTTP server.
func (a *APM) Close() error {
	a.mu.Lock()
	lis := a.listener
	a.listener = nil
	a.server = nil
	a.mu.Unlock()
	if lis != nil {
		return lis.Close()
	}
	return nil
}

// IsEnabled returns true if the APM service is enabled
func (a *APM) IsEnabled() bool {
	return a.enabled
}

// GetEndpoint returns the APM endpoint
func (a *APM) GetEndpoint() sbnet.Endpoint {
	return a.endpoint
}

// SetEndpoint sets the APM endpoint
func (a *APM) SetEndpoint(endpoint sbnet.Endpoint) {
	a.endpoint = endpoint
}
