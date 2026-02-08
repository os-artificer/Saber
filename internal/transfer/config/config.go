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

import "time"

type DiscoveryConfig struct {
	EtcdEndpoints         []string      `yaml:"etcdEndpoints"`
	EtcdUser              string        `yaml:"etcdUser"`
	EtcdPassword          string        `yaml:"etcdPassword"`
	DialTimeout           time.Duration `yaml:"dialTimeout"`
	AutoSyncInterval      time.Duration `yaml:"autoSyncInterval"`
	DialKeepAliveTime     time.Duration `yaml:"dialKeepAliveTime"`
	DialKeepAliveTimeout  time.Duration `yaml:"dialKeepAliveTimeout"`
	RegistryRootKeyPrefix string        `yaml:"registryRootKeyPrefix"` // default: /os-artificer/saber/registry
}

type ServiceConfig struct {
	ListenAddress string `yaml:"listenAddress"`
}

// LogConfig log config
type LogConfig struct {
	Path        string `yaml:"path"`
	FileSize    int    `yaml:"MaxSizeMB"`
	BackupCount int    `yaml:"Maxbackups"`
	BackupAge   int    `yaml:"MaxAge"`
}

// KafkaConfig Kafka producer configuration
type KafkaConfig struct {
	Brokers []string `yaml:"brokers"`
	Topic   string   `yaml:"topic"`
}

// Configuration transfer's configuration
type Configuration struct {
	Service ServiceConfig `yaml:"service"`
	Kafka   KafkaConfig   `yaml:"kafka"`
}

// Cfg global config (loaded by run)
var Cfg Configuration
