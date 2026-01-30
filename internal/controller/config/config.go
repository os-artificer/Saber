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

// DiscoveryConfig discovery's config
type DiscoveryConfig struct {
	Endpoints string `yaml:"endpoints"`
	User      string `yaml:"user"`
	Password  string `yaml:"password"`
}

// ServiceConfig service local config
type ServiceConfig struct {
	ListenAddress string `yaml:"listenAddress"`
}

// LogConfig log config
type LogConfig struct {
	FileName       string `yaml:"fileName"`
	FileSizeMB     int    `yaml:"fileSize"`
	MaxBackupCount int    `yaml:"maxBackupCount"`
	MaxBackupAge   int    `yaml:"maxBackupAge"`
}

// Configuration controller's configuration
type Configuration struct {
	Name      string          `yaml:"name"`
	Version   string          `yaml:"version"`
	Discovery DiscoveryConfig `yaml:"discovery"`
	Service   ServiceConfig   `yaml:"service"`
	Log       LogConfig       `yaml:"log"`
}
