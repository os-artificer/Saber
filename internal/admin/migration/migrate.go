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

package migration

import (
	"fmt"

	"os-artificer/saber/pkg/sbdb"
	"os-artificer/saber/pkg/sbmodels"

	"github.com/spf13/cobra"
	"gorm.io/gorm"
)

const bootstrapDB = "mysql"

// DBConfig holds MySQL connection parameters for migration (no dependency on admin config).
type DBConfig struct {
	Host     string
	Port     int
	User     string
	Password string
	Database string
	Charset  string
}

// NewMigrateCmd returns a cobra command that runs migration using config from getConfig.
func NewMigrateCmd(getConfig func() (*DBConfig, error)) *cobra.Command {
	return &cobra.Command{
		Use:   "migrate",
		Short: "Create database and run migrations",
		RunE: func(cmd *cobra.Command, args []string) error {
			cfg, err := getConfig()
			if err != nil {
				return err
			}
			return Run(cfg)
		},
	}
}

// Run creates the database if not exists, then runs AutoMigrate for all sbmodels.
func Run(cfg *DBConfig) error {
	if cfg == nil {
		return fmt.Errorf("migration: DBConfig is required")
	}
	dbName := cfg.Database
	if dbName == "" {
		dbName = sbmodels.DatabaseName
	}

	opts := []sbdb.Option{
		sbdb.OptionUser(cfg.User),
		sbdb.OptionPassword(cfg.Password),
		sbdb.OptionHost(cfg.Host),
		sbdb.OptionPort(cfg.Port),
	}
	if cfg.Charset != "" {
		opts = append(opts, sbdb.OptionCharset(cfg.Charset))
	}

	// Connect to system DB and create target database if not exists
	bootstrap, err := sbdb.NewMySQL(append(opts, sbdb.OptionDatabase(bootstrapDB))...)
	if err != nil {
		return fmt.Errorf("connect to %s: %w", bootstrapDB, err)
	}
	defer bootstrap.Close()

	sql := fmt.Sprintf("CREATE DATABASE IF NOT EXISTS `%s`", dbName)
	if err := bootstrap.DB().Exec(sql).Error; err != nil {
		return fmt.Errorf("create database %q: %w", dbName, err)
	}

	// Connect to target database and run migrations
	target, err := sbdb.NewMySQL(append(opts, sbdb.OptionDatabase(dbName))...)
	if err != nil {
		return fmt.Errorf("connect to %q: %w", dbName, err)
	}
	defer target.Close()

	if err := autoMigrate(target.DB()); err != nil {
		return err
	}
	return nil
}

func autoMigrate(db *gorm.DB) error {
	return db.AutoMigrate(
		&sbmodels.HostSnapshot{},
	)
}
