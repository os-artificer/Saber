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

package host

import (
	"encoding/json"
	"fmt"
	"time"
)

// Options is the option for the host plugin.
type Options struct {
	Interval time.Duration `yaml:"interval" json:"interval"`
	Timeout  time.Duration `yaml:"timeout" json:"timeout"`
}

// OptionsFromAny converts opts (any) to Options. Supports nil, Options, and map[string]any (via JSON).
func OptionsFromAny(opts any) (Options, error) {
	if opts == nil {
		return Options{}, nil
	}

	if o, ok := opts.(Options); ok {
		return o, nil
	}

	if m, ok := opts.(map[string]any); ok {
		data, err := json.Marshal(m)
		if err != nil {
			return Options{}, fmt.Errorf("host options marshal: %w", err)
		}

		var out Options
		if err := json.Unmarshal(data, &out); err != nil {
			return Options{}, fmt.Errorf("host options unmarshal: %w", err)
		}

		return out, nil
	}

	return Options{}, fmt.Errorf("unsupported options type: %T", opts)
}
