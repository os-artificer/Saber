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
)

const hostPluginVersion = "1.0.0"

func init() {
	RegisterPlugin("host", newHostPlugin)
}

// HostPlugin collects host metrics/info.
type HostPlugin struct {
	opts any
}

func newHostPlugin(ctx context.Context, opts any) (Plugin, error) {
	return &HostPlugin{opts: opts}, nil
}

func (p *HostPlugin) Version() string {
	return hostPluginVersion
}

func (p *HostPlugin) Name() string {
	return "host"
}

func (p *HostPlugin) Run(ctx context.Context) error {
	<-ctx.Done()
	return ctx.Err()
}

func (p *HostPlugin) Close() error {
	return nil
}
