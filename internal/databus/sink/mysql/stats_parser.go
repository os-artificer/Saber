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

package mysql

import "os-artificer/saber/pkg/sbmodels"

// databusPayload matches the JSON shape of plugin.Event sent by the harvester.
type databusPayload struct {
	PluginName string          `json:"PluginName"`
	EventName  string          `json:"EventName"`
	Data       *sbmodels.Stats `json:"Data"`
}

// collectIPs flattens IPs from all Networks into a single slice (order preserved, no dedup).
func collectIPs(networks []sbmodels.NetworkStats) []string {
	var out []string
	for _, n := range networks {
		out = append(out, n.IPs...)
	}
	return out
}
