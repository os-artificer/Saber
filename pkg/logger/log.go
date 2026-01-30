/**
 * Copyright 2025 saber authors.
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

package logger

import "log"

// Logger is a universal logging interface that can be
// flexibly replaced with other logging libraies
// without interfering with the operational logic
// of business code.
type Logger interface {
	Debug(format string, args ...any)
	Info(format string, args ...any)
	Warn(format string, args ...any)
	Error(format string, args ...any)
	Fatal(format string, args ...any)
}

type Level string

var (
	DebugLevel Level = "debug"
	InfoLevel  Level = "info"
	WarnLevel  Level = "warn"
	ErrorLevel Level = "error"
	FatalLevel Level = "fatal"
)

// Config logger config
type Config struct {
	Filename   string
	LogLevel   Level
	MaxSizeMB  int
	MaxBackups int
	MaxAge     int
}

var l Logger

func SetLogger(log Logger) {
	l = log
}

func Debug(format string, args ...interface{}) {
	if l == nil {
		log.Printf(format, args...)
		return
	}

	l.Debug(format, args...)
}

func Info(format string, args ...interface{}) {
	if l == nil {
		log.Printf(format, args...)
		return
	}

	l.Info(format, args...)
}

func Warn(format string, args ...interface{}) {
	if l == nil {
		log.Printf(format, args...)
		return
	}

	l.Warn(format, args...)
}

func Error(format string, args ...interface{}) {
	if l == nil {
		log.Printf(format, args...)
		return
	}

	l.Error(format, args...)
}

func Fatal(format string, args ...interface{}) {
	if l == nil {
		log.Fatalf(format, args...)
		return
	}

	l.Fatal(format, args...)
}
