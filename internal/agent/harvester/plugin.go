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

package harvester

import (
	"context"

	"os-artificer/saber/pkg/gerrors"
)

const (
	PluginVersionUnknown = "unknown"
)

var (
	ErrPluginUnimplemented = gerrors.New(gerrors.Unimplemented, "plugin is not implemented")
)

// Plugin is the interface for harvester plugins.
type Plugin interface {
	Version() string
	Name() string
	Run(ctx context.Context) error
	Close() error
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
