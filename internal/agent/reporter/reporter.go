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
	"sync"
)

// Reporter is the interface for data reporters (e.g. transfer, kafka).
type Reporter interface {
	SendMessage(ctx context.Context, content []byte) error
	Run() error
	Close() error
}

// ReporterFactory creates a Reporter from options (e.g. *config.Configuration).
type ReporterFactory func(ctx context.Context, opts any) (Reporter, error)

var (
	registry   = make(map[string]ReporterFactory)
	registryMu sync.RWMutex
)

// RegisterReporter registers a reporter implementation by type name.
func RegisterReporter(typeName string, factory ReporterFactory) {
	registryMu.Lock()
	defer registryMu.Unlock()
	registry[typeName] = factory
}

// CreateReporter creates a Reporter by type name and options.
func CreateReporter(ctx context.Context, typeName string, opts any) (Reporter, error) {
	registryMu.RLock()
	factory, ok := registry[typeName]
	registryMu.RUnlock()
	if !ok {
		return nil, errUnknownReporterType(typeName)
	}
	return factory(ctx, opts)
}

func errUnknownReporterType(name string) error {
	return fmt.Errorf("unknown reporter type: %s", name)
}
