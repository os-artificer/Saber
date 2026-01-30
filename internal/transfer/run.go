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

package transfer

import (
	"context"
	"net"
	"strings"

	"os-artificer/saber/internal/transfer/config"
	"os-artificer/saber/internal/transfer/service"
	"os-artificer/saber/pkg/logger"

	"github.com/segmentio/kafka-go"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func Run(cmd *cobra.Command, args []string) error {
	ctx := context.Background()

	address := ":26689"
	var kafkaWriter *kafka.Writer

	if ConfigFilePath != "" {
		viper.SetConfigFile(ConfigFilePath)
		viper.SetConfigType("yaml")
		if err := viper.ReadInConfig(); err == nil {
			var cfg config.Configuration
			if err := viper.Unmarshal(&cfg); err == nil {
				if cfg.Service.ListenAddress != "" {
					address = parseListenAddress(cfg.Service.ListenAddress)
				}
				if len(cfg.Kafka.Brokers) > 0 && cfg.Kafka.Topic != "" {
					kafkaWriter = &kafka.Writer{
						Addr:     kafka.TCP(cfg.Kafka.Brokers...),
						Topic:     cfg.Kafka.Topic,
						Balancer:  &kafka.LeastBytes{},
					}
					defer func() {
						if err := kafkaWriter.Close(); err != nil {
							logger.Warn("kafka writer close: %v", err)
						}
					}()
				}
			}
		}
	}

	svr := service.New(ctx, address, "", kafkaWriter)
	return svr.Run()
}

func parseListenAddress(addr string) string {
	if strings.HasPrefix(addr, "tcp://") {
		addr = strings.TrimPrefix(addr, "tcp://")
	}
	if host, port, err := net.SplitHostPort(addr); err == nil && host != "" && port != "" {
		return ":" + port
	}
	return addr
}
