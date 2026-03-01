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

// HostSnapshot is the model for the host snapshot table.
type HostSnapshot struct {
	ID        uint      `gorm:"column:id;type:bigint;not null;primaryKey;autoIncrement"`
	MachineID string    `gorm:"column:machine_id;type:varchar(64);not null;index:idx_machine_id" json:"machine_id"`
	HostName  string    `gorm:"column:host_name;type:varchar(255);not null"                      json:"host_name"`
	IPs       []string  `gorm:"column:ips;type:json;serializer:json"                             json:"ips"`
	MACs      []string  `gorm:"column:macs;type:json;serializer:json"                            json:"macs"`
	CPUs      int       `gorm:"column:cpu_count;not null;default:0"                              json:"cpu_count"`
	Memory    int       `gorm:"column:memory_size;not null;default:0"                            json:"memory_size"`
	Disk      int       `gorm:"column:disk_size;not null;default:0"                              json:"disk_size"`
	Network   int       `gorm:"column:network_count;not null;default:0"                          json:"network_count"`
	CreatedAt time.Time `gorm:"column:created_at;type:datetime;not null;default:current_timestamp"`
	UpdatedAt time.Time `gorm:"column:updated_at;type:datetime;not null;default:current_timestamp on update current_timestamp"`
	DeletedAt time.Time `gorm:"column:deleted_at;type:datetime;default:null"`
}

// TableName is the table name for the host snapshot model.
func (HostSnapshot) TableName() string {
	return "t_host_snapshots"
}
