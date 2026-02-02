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
	"runtime/debug"
	"sync"

	"os-artificer/saber/pkg/logger"
)

type Harvester struct {
	plugins map[string]Plugin
	mu      sync.RWMutex
	wg      sync.WaitGroup
}

// NewHarvester creates a new harvester with the given plugins.
func NewHarvester(plugins []Plugin) *Harvester {
	pluginsMap := make(map[string]Plugin)
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
		h.wg.Add(1)
		go func(plugin Plugin) {
			defer h.wg.Done()
			defer func() {
				if r := recover(); r != nil {
					logger.Errorf("plugin %s panic in Run: %v\n%s", plugin.Name(), r, string(debug.Stack()))
				}
			}()

			if err := plugin.Run(ctx); err != nil {
				logger.Errorf("failed to run plugin: %s, errmsg: %v", plugin.Name(), err)
			}
		}(plugin)
	}

	h.wg.Wait()
	return nil
}

func (h *Harvester) Close() error {
	h.mu.RLock()
	defer h.mu.RUnlock()

	for _, plugin := range h.plugins {
		h.wg.Add(1)
		go func(plugin Plugin) {
			defer h.wg.Done()
			defer func() {
				if r := recover(); r != nil {
					logger.Errorf("plugin %s panic in Close: %v\n%s", plugin.Name(), r, string(debug.Stack()))
				}
			}()

			if err := plugin.Close(); err != nil {
				logger.Errorf("failed to close plugin: %s, errmsg: %v", plugin.Name(), err)
			}
		}(plugin)
	}

	h.wg.Wait()
	return nil
}
