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

package apm

import (
	pkgapm "os-artificer/saber/pkg/apm"
)

// Business metrics for the admin service, registered to the default APM registry.
// Add gauges/counters as needed for admin-specific metrics.
var (
	// AdminInfo is a placeholder gauge for admin service info (e.g. version or status).
	AdminInfo = pkgapm.NewGauge(
		"admin",
		"info",
		"Admin service info placeholder",
		[]string{"version"},
	)
)

func init() {
	AdminInfo.WithLabelValues("v1").Set(1)
}
