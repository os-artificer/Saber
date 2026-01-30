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

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
)

type ZapLogger struct {
	logger        *zap.Logger
	sugaredLogger *zap.SugaredLogger
}

func (z *ZapLogger) Debug(format string, args ...interface{}) {
	z.sugaredLogger.Debugf(format, args...)
}

func (z *ZapLogger) Info(format string, args ...interface{}) {
	z.sugaredLogger.Infof(format, args...)
}

func (z *ZapLogger) Warn(format string, args ...interface{}) {
	z.sugaredLogger.Warnf(format, args...)
}

func (z *ZapLogger) Error(format string, args ...interface{}) {
	z.sugaredLogger.Errorf(format, args...)
}

func (z *ZapLogger) Fatal(format string, args ...interface{}) {
	z.sugaredLogger.Fatalf(format, args...)
}

func (z *ZapLogger) Sync() error {
	return z.logger.Sync()
}

func convertLevel(level Level) zapcore.Level {
	switch level {
	case DebugLevel:
		return zapcore.DebugLevel
	case InfoLevel:
		return zapcore.InfoLevel
	case ErrorLevel:
		return zapcore.ErrorLevel
	case FatalLevel:
		return zapcore.FatalLevel
	default:
		return zapcore.InfoLevel
	}
}

// NewZapLogger create a zap logger
func NewZapLogger(config Config) Logger {

	logRotator := &lumberjack.Logger{
		Filename:   config.Filename,
		MaxSize:    config.MaxSizeMB,
		MaxBackups: config.MaxBackups,
		MaxAge:     config.MaxAge,
		Compress:   true,
	}

	cfg := zap.NewProductionConfig()
	cfg.EncoderConfig = zapcore.EncoderConfig{
		MessageKey:     "msg",
		LevelKey:       "level",
		TimeKey:        "time",
		NameKey:        "name",
		FunctionKey:    "func",
		StacktraceKey:  "stacktrace",
		CallerKey:      "caller",
		SkipLineEnding: false,
		LineEnding:     zapcore.DefaultLineEnding,
		EncodeLevel:    zapcore.LowercaseLevelEncoder,
		EncodeTime:     zapcore.RFC3339TimeEncoder,
		EncodeDuration: zapcore.SecondsDurationEncoder,
		EncodeName:     zapcore.FullNameEncoder,
		EncodeCaller:   zapcore.ShortCallerEncoder,
	}

	encoder := zapcore.NewConsoleEncoder(cfg.EncoderConfig)
	core := zapcore.NewCore(encoder, zapcore.AddSync(logRotator), convertLevel(config.LogLevel))
	logger := zap.New(core, zap.AddCaller())

	return &ZapLogger{logger: logger, sugaredLogger: logger.Sugar()}
}
