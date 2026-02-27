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

// Business metrics for databus, registered to the default APM registry.
// Update them from source/sink/handler code.
var (
	// ConnectionsActive is the number of active agent connections.
	// Labels: state (e.g. "connected", "disconnected").
	ConnectionsActive = pkgapm.NewGauge(
		"databus",
		"connections_active",
		"Number of active agent connections",
		[]string{"state"},
	)

	// RecordsWrittenTotal is the total number of records written to sink.
	// Labels: sink_type, status.
	RecordsWrittenTotal = pkgapm.NewCounter(
		"databus",
		"records_written_total",
		"Total number of records written to sink",
		[]string{"sink_type", "status"},
	)

	// BytesWrittenTotal is the total number of bytes written to sink.
	// Labels: sink_type.
	BytesWrittenTotal = pkgapm.NewCounter(
		"databus",
		"bytes_written_total",
		"Total number of bytes written to sink",
		[]string{"sink_type"},
	)
)

func init() {
	ConnectionsActive.WithLabelValues("connected").Set(0)
	RecordsWrittenTotal.WithLabelValues("kafka", "success").Add(0)
	BytesWrittenTotal.WithLabelValues("kafka").Add(0)
}
