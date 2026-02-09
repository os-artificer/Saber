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

package plugin

import (
	"context"
	"fmt"
	"sync"

	"os-artificer/saber/pkg/gerrors"
)

const (
	PluginVersionUnknown = "unknown"
)

var (
	ErrPluginUnimplemented = gerrors.New(gerrors.Unimplemented, "plugin is not implemented")
)

// Event is the event for harvester plugins.
type Event struct {
	PluginName string
	EventName  string
	Data       any
}

// EventC is the channel for sending events to harvester plugins.
type EventC chan *Event

// Plugin is the interface for harvester plugins.
type Plugin interface {
	Version() string
	Name() string
	Run(ctx context.Context) (EventC, error)
	Close() error
}

// PluginFactory creates a Plugin from options (e.g. map from config).
type PluginFactory func(ctx context.Context, opts any) (Plugin, error)

var (
	pluginRegistry   = make(map[string]PluginFactory)
	pluginRegistryMu sync.RWMutex
)

// PluginConfig is the config for creating one plugin instance.
type PluginConfig struct {
	Name    string
	Options any
}

// RegisterPlugin registers a harvester plugin by name.
func RegisterPlugin(name string, factory PluginFactory) {
	pluginRegistryMu.Lock()
	defer pluginRegistryMu.Unlock()
	pluginRegistry[name] = factory
}

// CreatePlugins creates plugins from config entries.
func CreatePlugins(ctx context.Context, configs []PluginConfig) ([]Plugin, error) {
	if len(configs) == 0 {
		return nil, nil
	}

	out := make([]Plugin, 0, len(configs))
	for _, c := range configs {
		pluginRegistryMu.RLock()
		factory, ok := pluginRegistry[c.Name]
		pluginRegistryMu.RUnlock()

		if !ok {
			return nil, fmt.Errorf("unknown harvester plugin: %s", c.Name)
		}

		p, err := factory(ctx, c.Options)
		if err != nil {
			return nil, err
		}

		out = append(out, p)
	}

	return out, nil
}

// UnimplementedPlugin is the default implementation of the Plugin interface.
type UnimplementedPlugin struct {
}

func (p *UnimplementedPlugin) Version() string {
	return ""
}

func (p *UnimplementedPlugin) Name() string {
	return ""
}

func (p *UnimplementedPlugin) Run(ctx context.Context) error {
	return ErrPluginUnimplemented
}

func (p *UnimplementedPlugin) Close() error {
	return ErrPluginUnimplemented
}
