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

package sbdb

import (
	"context"
	"errors"
	"fmt"
	"time"

	"os-artificer/saber/pkg/logger"

	gormlogger "gorm.io/gorm/logger"
)

// gormLogger adapts the project logger to gorm.Logger.Interface so GORM logs go through saber's logger.
type gormLogger struct {
	SlowThreshold             time.Duration
	LogLevel                  gormlogger.LogLevel
	IgnoreRecordNotFoundError bool
}

// newGormLogger returns a GORM logger that writes to the project logger.
func newGormLogger() gormlogger.Interface {
	return &gormLogger{
		SlowThreshold:             200 * time.Millisecond,
		LogLevel:                  gormlogger.Warn,
		IgnoreRecordNotFoundError: true,
	}
}

// LogMode implements gorm.Logger.Interface.
func (l *gormLogger) LogMode(level gormlogger.LogLevel) gormlogger.Interface {
	nl := *l
	nl.LogLevel = level
	return &nl
}

// Info implements gorm.Logger.Interface.
func (l *gormLogger) Info(ctx context.Context, msg string, data ...interface{}) {
	if l.LogLevel >= gormlogger.Info {
		logger.Infof("[gorm] "+msg, data...)
	}
}

// Warn implements gorm.Logger.Interface.
func (l *gormLogger) Warn(ctx context.Context, msg string, data ...interface{}) {
	if l.LogLevel >= gormlogger.Warn {
		logger.Warnf("[gorm] "+msg, data...)
	}
}

// Error implements gorm.Logger.Interface.
func (l *gormLogger) Error(ctx context.Context, msg string, data ...interface{}) {
	if l.LogLevel >= gormlogger.Error {
		logger.Errorf("[gorm] "+msg, data...)
	}
}

// Trace implements gorm.Logger.Interface.
func (l *gormLogger) Trace(ctx context.Context, begin time.Time, fc func() (sql string, rowsAffected int64), err error) {
	if l.LogLevel <= gormlogger.Silent {
		return
	}

	elapsed := time.Since(begin)
	sql, rows := fc()
	ms := float64(elapsed.Nanoseconds()) / 1e6
	rowsStr := fmt.Sprintf("%d", rows)
	if rows == -1 {
		rowsStr = "-"
	}

	switch {
	case err != nil && l.LogLevel >= gormlogger.Error && (!errors.Is(err, gormlogger.ErrRecordNotFound) || !l.IgnoreRecordNotFoundError):
		logger.Errorf("[gorm] %v | %.3fms | rows:%s | %s", err, ms, rowsStr, sql)

	case elapsed > l.SlowThreshold && l.SlowThreshold != 0 && l.LogLevel >= gormlogger.Warn:
		slowLog := fmt.Sprintf("SLOW SQL >= %v", l.SlowThreshold)
		logger.Warnf("[gorm] %s | %.3fms | rows:%s | %s", slowLog, ms, rowsStr, sql)

	case l.LogLevel >= gormlogger.Info:
		logger.Infof("[gorm] %.3fms | rows:%s | %s", ms, rowsStr, sql)
	}
}
