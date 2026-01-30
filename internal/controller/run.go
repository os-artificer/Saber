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
	"net"
	"os"
	"os/signal"
	"strings"
	"syscall"

	"os-artificer/saber/internal/controller/config"
	"os-artificer/saber/internal/controller/service"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func setupGracefulShutdown(svr *service.Service) {
	sigC := make(chan os.Signal, 1)
	signal.Notify(sigC, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		<-sigC

		os.Exit(0)
	}()
}

func parseListenAddress(addr string) string {
	addr = strings.TrimPrefix(addr, "tcp://")
	if _, port, err := net.SplitHostPort(addr); err == nil && port != "" {
		return ":" + port
	}
	return addr
}

func Run(cmd *cobra.Command, args []string) error {
	ctx := context.Background()

	address := ":26688"
	if ConfigFilePath != "" {
		viper.SetConfigFile(ConfigFilePath)
		viper.SetConfigType("yaml")
		if err := viper.ReadInConfig(); err == nil {
			var cfg config.Configuration
			if err := viper.Unmarshal(&cfg); err == nil && cfg.Service.ListenAddress != "" {
				address = parseListenAddress(cfg.Service.ListenAddress)
			}
		}
	}

	svr := service.New(ctx, address, "")
	setupGracefulShutdown(svr)

	return svr.Run()
}
