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

package transfer

import (
	"context"
	"net"
	"strings"

	"os-artificer/saber/internal/transfer/config"
	"os-artificer/saber/internal/transfer/sink"
	sinkkafka "os-artificer/saber/internal/transfer/sink/kafka"
	"os-artificer/saber/internal/transfer/source"
	"os-artificer/saber/pkg/logger"

	"github.com/segmentio/kafka-go"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func Run(cmd *cobra.Command, args []string) error {
	ctx := context.Background()
	address, out, closeSink := loadTransferConfig()
	if closeSink != nil {
		defer closeSink()
	}

	handler := source.NewConnectionHandler(out)
	defer handler.Close()

	agentSource := source.NewAgentSource(address, nil)
	return agentSource.Run(ctx, handler)
}

// loadTransferConfig reads config from ConfigFilePath and returns listen address, optional sink, and optional cleanup.
func loadTransferConfig() (address string, out sink.Sink, closeSink func()) {
	address = ":26689"

	if ConfigFilePath == "" {
		return address, nil, nil
	}

	viper.SetConfigFile(ConfigFilePath)
	viper.SetConfigType("yaml")

	if err := viper.ReadInConfig(); err != nil {
		return address, nil, nil
	}

	var cfg config.Configuration

	if err := viper.Unmarshal(&cfg); err != nil {
		return address, nil, nil
	}

	if cfg.Service.ListenAddress != "" {
		address = parseListenAddress(cfg.Service.ListenAddress)
	}

	if len(cfg.Kafka.Brokers) == 0 || cfg.Kafka.Topic == "" {
		return address, nil, nil
	}

	writer := &kafka.Writer{
		Addr:     kafka.TCP(cfg.Kafka.Brokers...),
		Topic:    cfg.Kafka.Topic,
		Balancer: &kafka.LeastBytes{},
	}

	out = sinkkafka.New(writer)

	closeSink = func() {
		if err := out.Close(); err != nil {
			logger.Warnf("sink close: %v", err)
		}
	}

	return address, out, closeSink
}

func parseListenAddress(addr string) string {
	addr = strings.TrimPrefix(addr, "tcp://")
	if _, port, err := net.SplitHostPort(addr); err == nil && port != "" {
		return ":" + port
	}
	return addr
}
