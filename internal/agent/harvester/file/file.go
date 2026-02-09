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

package file

import (
	"context"

	"os-artificer/saber/internal/agent/harvester/plugin"
	"os-artificer/saber/pkg/logger"
)

const filePluginVersion = "1.0.0"

func init() {
	plugin.RegisterPlugin("file", newFilePlugin)
}

// FilePlugin collects data from files.
type FilePlugin struct {
	opts any
}

func newFilePlugin(ctx context.Context, opts any) (plugin.Plugin, error) {
	return &FilePlugin{opts: opts}, nil
}

func (p *FilePlugin) Version() string {
	return filePluginVersion
}

func (p *FilePlugin) Name() string {
	return "file"
}

func (p *FilePlugin) Run(ctx context.Context) (plugin.EventC, error) {
	eventC := make(plugin.EventC)

	go func() {
		defer close(eventC)

		for {
			select {
			case <-ctx.Done():
				logger.Infof("file plugin run exited: %s", p.Name())
				return

			default:
				eventC <- &plugin.Event{
					PluginName: p.Name(),
					EventName:  "file",
					Data:       nil,
				}
			}
		}
	}()

	return eventC, nil
}

func (p *FilePlugin) Close() error {
	return nil
}
