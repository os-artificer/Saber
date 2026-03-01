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

package sbmodels

import "time"

// HostSnapshot table column names (for raw SQL / Assign maps).
const (
	HostSnapshotColID        = "id"
	HostSnapshotColMachineID = "machine_id"
	HostSnapshotColHostName  = "host_name"
	HostSnapshotColIPs       = "ips"
	HostSnapshotColStats     = "stats"
	HostSnapshotColCreatedAt = "created_at"
	HostSnapshotColUpdatedAt = "updated_at"
	HostSnapshotColDeletedAt = "deleted_at"
)

// HostSnapshot is the model for the host snapshot table.
type HostSnapshot struct {
	ID        uint                `gorm:"column:id;type:bigint;not null;primaryKey;autoIncrement"`
	MachineID string              `gorm:"column:machine_id;type:varchar(64);not null;index:idx_machine_id"`
	HostName  string              `gorm:"column:host_name;type:varchar(255);not null"`
	IPs       JSONValue[[]string] `gorm:"column:ips;type:json;serializer:json"`
	Stats     JSONValue[Stats]    `gorm:"column:stats;type:json"`
	CreatedAt time.Time           `gorm:"column:created_at;type:datetime;not null;default:current_timestamp"`
	UpdatedAt time.Time           `gorm:"column:updated_at;type:datetime;not null;default:current_timestamp on update current_timestamp"`
	DeletedAt time.Time           `gorm:"column:deleted_at;type:datetime;default:null"`
}

// TableName is the table name for the host snapshot model.
func (HostSnapshot) TableName() string {
	return "t_host_snapshots"
}
