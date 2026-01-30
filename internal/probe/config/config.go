/**
 * Copyright 2025 saber authors.
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

var Cfg = Configuration{
	Name:       "Probe",
	Version:    "v1.0.0",
	Discovery:  DiscoveryConfig{},
	Controller: ControllerConfig{},
	Transfer:   TransferConfig{},
	Collector:  CollectorConfig{Interval: 10},
	Engine: EngineConfig{
		Script: ScriptEngine{},
		Data:   DataEngine{},
	},
	Log: LogConfig{},
}

// DiscoveryConfig discovery configuration
type DiscoveryConfig struct {
	Endpoints        string `yaml:"endpoints"`
	SyncMetaInterval int    `yaml:"syncMetaInterval"`
}

// ControllerConfig controller service configuration
type ControllerConfig struct {
	Endpoints        string `yaml:"endpoints"`
	SyncMetaInterval int    `yaml:"syncMetaInterval"`
}

// TransferConfig transfer service configuration
type TransferConfig struct {
	Endpoints string `yaml:"endpoints"`
}

// CollectorConfig host metrics collector configuration
type CollectorConfig struct {
	Interval int `yaml:"interval"` // seconds
}

// ScriptEngine script engine
type ScriptEngine struct {
	Enable      bool   `yaml:"enable"`
	Concurrency int    `yaml:"concurrency"`
	TmpDir      string `yaml:"tmpDir"`
}

// DataEngine data engine
type DataEngine struct {
	Enable   bool   `yaml:"enable"`
	SendMode string `yaml:"sendMode"`
}

// EngineConfig engine config
type EngineConfig struct {
	Script ScriptEngine `yaml:"script"`
	Data   DataEngine  `yaml:"data"`
}

// LogConfig log config
type LogConfig struct {
	FileName       string `yaml:"fileName"`
	FileSize       int    `yaml:"fileSize"`
	MaxBackupCount int    `yaml:"maxBackupCount"`
	MaxBackupAge   int    `yaml:"maxBackupAge"`
}

// Configuration probe's configuration
type Configuration struct {
	Name       string           `yaml:"name"`
	Version    string           `yaml:"version"`
	Discovery  DiscoveryConfig  `yaml:"discovery"`
	Controller ControllerConfig `yaml:"controller"`
	Transfer   TransferConfig   `yaml:"transfer"`
	Collector  CollectorConfig  `yaml:"collector"`
	Engine     EngineConfig     `yaml:"engine"`
	Log        LogConfig        `yaml:"log"`
}
