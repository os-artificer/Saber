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
	"sync"

	"os-artificer/saber/internal/agent/harvester/plugin"
	"os-artificer/saber/pkg/logger"
	"os-artificer/saber/pkg/tools"
)

type Harvester struct {
	plugins map[string]plugin.Plugin
	mu      sync.RWMutex
	runWg   sync.WaitGroup // used only by Run()
	closeWg sync.WaitGroup // used only by Close(); separate to avoid WaitGroup contract violation
}

// NewHarvester creates a new harvester with the given plugins.
func NewHarvester(plugins []plugin.Plugin) *Harvester {
	pluginsMap := make(map[string]plugin.Plugin)
	for _, plugin := range plugins {
		if _, ok := pluginsMap[plugin.Name()]; ok {
			continue
		}
		pluginsMap[plugin.Name()] = plugin
	}
	return &Harvester{plugins: pluginsMap}
}

func (h *Harvester) Run(ctx context.Context) error {
	h.mu.RLock()
	defer h.mu.RUnlock()

	for _, plugin := range h.plugins {
		h.runWg.Add(1)
		plugin := plugin

		tools.Go(func() {
			defer h.runWg.Done()

			eventC, err := plugin.Run(ctx)
			if err != nil {
				logger.Errorf("failed to run plugin: %s, errmsg: %v", plugin.Name(), err)
				return
			}

			for {
				select {
				case <-ctx.Done():
					logger.Infof("harvester run exited: %s", plugin.Name())
					return

				case event := <-eventC:
					logger.Infof("harvester received event: %s, event: %v", plugin.Name(), event)
				}
			}
		})
	}

	h.runWg.Wait()
	return nil
}

func (h *Harvester) Close() error {
	h.mu.RLock()
	defer h.mu.RUnlock()

	for _, plugin := range h.plugins {
		h.closeWg.Add(1)
		plugin := plugin
		tools.Go(func() {
			defer h.closeWg.Done()

			if err := plugin.Close(); err != nil {
				logger.Errorf("failed to close plugin: %s, errmsg: %v", plugin.Name(), err)
			}
		})
	}

	h.closeWg.Wait()
	return nil
}
