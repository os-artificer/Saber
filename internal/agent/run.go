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

package agent

import (
	"context"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"os-artificer/saber/internal/agent/config"
	"os-artificer/saber/internal/agent/reporter"
	"os-artificer/saber/pkg/logger"

	"github.com/spf13/cobra"
)

func setupGracefulShutdown() {
	sigC := make(chan os.Signal, 1)
	signal.Notify(sigC, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		<-sigC
		os.Exit(0)
	}()
}

func Run(cmd *cobra.Command, args []string) error {
	ctx := context.Background()
	cfg := config.Cfg

	transferClient, err := reporter.NewTransferClient(ctx, cfg.Transfer.Endpoints, cfg.Name+"-"+cfg.Version)
	if err != nil {
		logger.Fatalf("Failed to create transfer client: %v", err)
	}

	setupGracefulShutdown()

	var wg sync.WaitGroup
	wg.Add(1)

	go func() {
		defer wg.Done()
		if err := transferClient.Run(); err != nil {
			logger.Warnf("Transfer client exited: %v", err)
		}
	}()

	// Host metrics collector: periodically collect and send via transfer client
	interval := time.Duration(cfg.Collector.Interval) * time.Second
	if interval <= 0 {
		interval = 10 * time.Second
	}

	go func() {
		time.Sleep(interval) // wait for transfer stream to be ready
		ticker := time.NewTicker(interval)
		defer ticker.Stop()

		for range ticker.C {
			// payload, err := collector.Collect()
			// if err != nil {
			// 	logger.Warnf("Collect metrics failed: %v", err)
			// 	continue
			// }

			// if err := transferClient.SendMessage(payload); err != nil {
			// 	logger.Warnf("Send metrics to transfer failed: %v", err)
			// }
		}
	}()

	wg.Wait()
	return nil
}
