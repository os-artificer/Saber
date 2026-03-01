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

import (
	"context"
	"encoding/json"

	"os-artificer/saber/internal/databus/sink/base"
	"os-artificer/saber/pkg/logger"
	"os-artificer/saber/pkg/proto"
	"os-artificer/saber/pkg/sbdb"
	"os-artificer/saber/pkg/sbmodels"
)

var _ base.Sink = (*MySQLSink)(nil)

// MySQLSink implements base.Sink by writing DatabusRequest to MySQL.
type MySQLSink struct {
	db *sbdb.MySQL
}

// NewMySQLSink returns a Sink that writes to the given MySQL database.
func NewMySQLSink(db *sbdb.MySQL) *MySQLSink {
	return &MySQLSink{db: db}
}

// Write implements base.Sink.
func (m *MySQLSink) Write(ctx context.Context, req *proto.DatabusRequest) error {
	logger.Debugf("write databus request to mysql, client_id=%s, payload_len=%d", req.GetClientID(), len(req.GetPayload()))

	if len(req.GetPayload()) == 0 {
		return nil
	}

	var p databusPayload
	if err := json.Unmarshal(req.GetPayload(), &p); err != nil {
		logger.Errorf("mysql sink: unmarshal payload failed: %v", err)
		return err
	}

	if p.PluginName != "host" {
		return nil
	}

	if req.GetClientID() == "" {
		logger.Warnf("mysql sink: client_id empty, skip host snapshot")
		return nil
	}

	if p.Data == nil {
		logger.Warnf("mysql sink: host data nil, skip")
		return nil
	}

	ips := collectIPs(p.Data.Networks)
	snapshot := sbmodels.HostSnapshot{
		MachineID: req.GetClientID(),
		HostName:  p.Data.Hostname,
		IPs:       sbmodels.JSONValueOf(&ips),
		Stats:     sbmodels.JSONValueOf(p.Data),
	}

	db := m.db.DB()
	// Upsert: insert when no row for machine_id, otherwise update host_name/ips/stats.
	err := db.WithContext(ctx).Where(sbmodels.HostSnapshotColMachineID+" = ?", snapshot.MachineID).
		Assign(map[string]any{
			sbmodels.HostSnapshotColHostName: snapshot.HostName,
			sbmodels.HostSnapshotColIPs:      snapshot.IPs,
			sbmodels.HostSnapshotColStats:    snapshot.Stats,
		}).FirstOrCreate(&snapshot).Error
	if err != nil {
		logger.Errorf("mysql sink: upsert host snapshot failed: %v", err)
		return err
	}
	return nil
}

// Close implements base.Sink.
func (m *MySQLSink) Close() error {
	if m.db == nil {
		return nil
	}

	if err := m.db.Close(); err != nil {
		logger.Warnf("close mysql sink failed: %v", err)
		return err
	}

	return nil
}
