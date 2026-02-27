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

package source

import (
	"os-artificer/saber/pkg/gerrors"
	"os-artificer/saber/pkg/sbnet"

	"github.com/go-viper/mapstructure/v2"
)

var (
	ErrNoEndpoint = gerrors.New(gerrors.InvalidParameter, "no endpoint configured")
)

type Config struct {
	Endpoint sbnet.Endpoint `yaml:"endpoint" mapstructure:"endpoint"`
}

// ConfigFromMap decodes m into *Config using mapstructure with Endpoint hooks.
func ConfigFromMap(m map[string]any) (*Config, error) {
	var cfg Config
	dec, err := mapstructure.NewDecoder(&mapstructure.DecoderConfig{
		DecodeHook: sbnet.StringToEndpointHookFunc(),
		Result:     &cfg,
	})
	if err != nil {
		return nil, err
	}
	if err := dec.Decode(m); err != nil {
		return nil, err
	}
	return &cfg, nil
}

// ListenAddress returns the endpoint for use as the listen address.
func (c *Config) ListenAddress() (sbnet.Endpoint, error) {
	if c.Endpoint.Host == "" && c.Endpoint.Port == 0 {
		return sbnet.Endpoint{}, ErrNoEndpoint
	}
	return c.Endpoint, nil
}
