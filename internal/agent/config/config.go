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

package config

import (
	"time"
)

var Cfg = Configuration{
	Name:       "Agent",
	Version:    "v1.0.0",
	Discovery:  DiscoveryConfig{},
	Controller: ControllerConfig{},
	Reporters: []ReporterEntry{
		{Type: "transfer", Config: map[string]any{"endpoints": "127.0.0.1:26689"}},
	},
	Harvester: HarvesterConfig{},
	Log:       LogConfig{},
}

// DiscoveryConfig discovery configuration
type DiscoveryConfig struct {
	Endpoints        string        `yaml:"endpoints"`
	SyncMetaInterval time.Duration `yaml:"syncMetaInterval"`
}

// ControllerConfig controller service configuration
type ControllerConfig struct {
	Endpoints        string        `yaml:"endpoints"`
	SyncMetaInterval time.Duration `yaml:"syncMetaInterval"`
}

// ReporterEntry config for one reporter (e.g. type + config in reporters list).
type ReporterEntry struct {
	Type   string         `yaml:"type"`
	Config map[string]any `yaml:"config"`
}

// ReporterOpts is passed to reporter.CreateReporter when creating a reporter from an entry.
type ReporterOpts struct {
	Config       map[string]any
	AgentName    string
	AgentVersion string
}

// HarvesterPluginEntry config for one harvester plugin
type HarvesterPluginEntry struct {
	Name    string         `yaml:"name"`
	Options map[string]any `yaml:"options"`
}

// HarvesterConfig harvester plugins configuration
type HarvesterConfig struct {
	Plugins []HarvesterPluginEntry `yaml:"plugins"`
}

// LogConfig log config
type LogConfig struct {
	FileName       string `yaml:"fileName"`
	FileSize       int    `yaml:"fileSize"`
	MaxBackupCount int    `yaml:"maxBackupCount"`
	MaxBackupAge   int    `yaml:"maxBackupAge"`
}

// Configuration agent's configuration
type Configuration struct {
	Name       string           `yaml:"name"`
	Version    string           `yaml:"version"`
	Discovery  DiscoveryConfig  `yaml:"discovery"`
	Controller ControllerConfig `yaml:"controller"`
	Reporters  []ReporterEntry  `yaml:"reporters"`
	Harvester  HarvesterConfig  `yaml:"harvester"`
	Log        LogConfig        `yaml:"log"`
}
