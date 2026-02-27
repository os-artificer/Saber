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
	"fmt"

	"os-artificer/saber/internal/databus/config"
	"os-artificer/saber/pkg/proto"
)

var (
	ErrSourceTypeNotSupported = fmt.Errorf("source type not supported")
)

// Handler processes a single DatabusRequest (e.g. write to Kafka, log).
// All source implementations deliver received requests to a Handler.
type Handler interface {
	OnDatabusRequest(req *proto.DatabusRequest) error
}

// Source is a data source that runs and delivers DatabusRequests to the handler
// until context is done. Implementations may start a gRPC server or other acceptor.
type Source interface {
	Run(ctx context.Context, h Handler) error
}

// NewSourceFromConfig creates a new source from a config.SourceConfig.
// Supported source types: agent.
func NewSourceFromConfig(cfg *config.SourceConfig) (Source, error) {
	switch cfg.Type {
	case config.SourceTypeAgent:
		cfg, err := ConfigFromMap(cfg.Config)
		if err != nil {
			return nil, err
		}

		address, err := cfg.ListenAddress()
		if err != nil {
			return nil, err
		}
		return NewAgentSource(address, nil), nil

	default:
		return nil, ErrSourceTypeNotSupported
	}
}

// NewSourcesFromConfig creates sources from a slice of config.SourceConfig.
// Returns (nil, error) if any config fails to create a source.
func NewSourcesFromConfig(cfgs []config.SourceConfig) ([]Source, error) {
	if len(cfgs) == 0 {
		return []Source{}, nil
	}

	sources := make([]Source, 0, len(cfgs))
	for i := range cfgs {
		src, err := NewSourceFromConfig(&cfgs[i])
		if err != nil {
			return nil, err
		}
		sources = append(sources, src)
	}
	return sources, nil
}
