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

package controller

import (
	"context"

	"os-artificer/saber/internal/controller/apm"
	"os-artificer/saber/internal/controller/config"
	"os-artificer/saber/internal/controller/server"
	"os-artificer/saber/pkg/logger"
	"os-artificer/saber/pkg/sbnet"

	"github.com/go-viper/mapstructure/v2"
	"github.com/spf13/viper"
)

// controllerUnmarshalOpt composes default viper hooks with string->Endpoint so
// service.listenAddress (string) unmarshals into sbnet.Endpoint.
var controllerUnmarshalOpt = viper.DecodeHook(mapstructure.ComposeDecodeHookFunc(
	mapstructure.StringToTimeDurationHookFunc(),
	mapstructure.StringToSliceHookFunc(","),
	sbnet.StringToEndpointHookFunc(),
))

// Service is the controller service.
type Service struct {
	svr *server.AgentServer
	apm *apm.APM
}

// CreateService creates a new controller service. apmSvc may be nil if APM is disabled.
func CreateService(ctx context.Context, address sbnet.Endpoint, serviceID string, apmSvc *apm.APM) *Service {
	svr := server.New(ctx, address, serviceID)
	return &Service{
		svr: svr,
		apm: apmSvc,
	}
}

// Run starts the controller service. If APM is enabled, it is started in a goroutine before the gRPC server runs.
func (s *Service) Run() error {
	if s.apm != nil && s.apm.IsEnabled() {
		go func() {
			_ = s.apm.Run()
		}()
	}
	return s.svr.Run()
}

// Close stops the controller service (APM first, then the gRPC server).
func (s *Service) Close() error {
	if s.apm != nil {
		_ = s.apm.Close()
	}
	return s.svr.Close()
}

// loadControllerConfig reads config from ConfigFilePath and returns the listen address.
func loadControllerConfig() {
	if ConfigFilePath == "" {
		return
	}

	viper.SetConfigFile(ConfigFilePath)
	viper.SetConfigType("yaml")

	if err := viper.ReadInConfig(); err != nil {
		logger.Errorf("Failed to read config file: %v", err)
		return
	}

	if err := viper.Unmarshal(&config.Cfg, controllerUnmarshalOpt); err != nil {
		logger.Errorf("Failed to unmarshal config: %v", err)
		return
	}

	logger.Infof("Loaded controller config: %+v", config.Cfg)
}
