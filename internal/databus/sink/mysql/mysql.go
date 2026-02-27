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

	"os-artificer/saber/internal/databus/sink"
	"os-artificer/saber/pkg/logger"
	"os-artificer/saber/pkg/proto"
	"os-artificer/saber/pkg/sbdb"
)

var _ sink.Sink = (*MySQLSink)(nil)

// MySQLSink implements sink.Sink by writing DatabusRequest to MySQL.
type MySQLSink struct {
	db *sbdb.MySQL
}

// NewMySQLSink returns a Sink that writes to the given MySQL database.
func NewMySQLSink(db *sbdb.MySQL) *MySQLSink {
	return &MySQLSink{db: db}
}

// Write implements sink.Sink.
func (m *MySQLSink) Write(ctx context.Context, req *proto.DatabusRequest) error {
	logger.Debugf("write databus request to mysql: %v", req)

	// TODO: write to mysql

	return nil
}

// Close implements sink.Sink.
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
