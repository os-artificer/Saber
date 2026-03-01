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

package admin

import (
	"fmt"

	"os-artificer/saber/internal/admin/config"
	"os-artificer/saber/internal/admin/migration"
	"os-artificer/saber/pkg/sbnet"
)

// GetDBConfigForMigrate loads admin config and returns migration DBConfig for the migrate command.
// Uses service.storage when type=mysql; parses storage.config.url for host/port.
func GetDBConfigForMigrate() (*migration.DBConfig, error) {
	loadAdminConfig()
	s := config.Cfg.Service.Storage
	if s == nil || s.Type != "mysql" {
		return nil, fmt.Errorf("migrate requires service.storage with type=mysql in config file (e.g. -c admin.yaml)")
	}
	ep, err := sbnet.NewEndpointFromString(s.Config.URL)
	if err != nil {
		return nil, fmt.Errorf("service.storage.config.url %q: %w", s.Config.URL, err)
	}
	return &migration.DBConfig{
		Host:     ep.Host,
		Port:     ep.Port,
		User:     s.Config.Username,
		Password: s.Config.Password,
		Database: s.Config.Database,
	}, nil
}
