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
	"sync"
	"time"

	"os-artificer/saber/internal/agent/harvester/plugin"
	"os-artificer/saber/pkg/logger"
)

const filePluginVersion = "1.0.0"

func init() {
	plugin.RegisterPlugin("file", newFilePlugin)
}

// FilePlugin collects data from files.
type FilePlugin struct {
	plugin.UnimplementedPlugin

	wg   sync.WaitGroup
	done chan struct{}
	opts any
}

func newFilePlugin(ctx context.Context, opts any) (plugin.Plugin, error) {
	return &FilePlugin{opts: opts, done: make(chan struct{})}, nil
}

func (p *FilePlugin) Version() string {
	return filePluginVersion
}

func (p *FilePlugin) Name() string {
	return "file"
}

func (p *FilePlugin) Run(ctx context.Context) (plugin.EventC, error) {
	eventC := make(plugin.EventC)

	p.wg.Add(1)
	go func() {
		defer p.wg.Done()
		defer close(eventC)

		for {
			select {
			case <-p.done:
				logger.Infof("file plugin run exited: %s", p.Name())
				return

			case <-ctx.Done():
				logger.Infof("file plugin run exited: %s", p.Name())
				return

			case <-time.After(1 * time.Second):
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
	if p.done != nil {
		close(p.done)
		p.done = nil
	}

	p.wg.Wait()
	return nil
}
