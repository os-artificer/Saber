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

package main

import (
	"os-artificer/saber/internal/agent"
	"os-artificer/saber/pkg/logger"

	"github.com/spf13/cobra"
)

func main() {
	rootCmd := &cobra.Command{
		Use:          "Agent",
		Short:        "saber Agent",
		SilenceUsage: true,
		RunE:         agent.Run,
	}

	rootCmd.PersistentFlags().StringVarP(&agent.ConfigFilePath, "config", "c", "/var/ylg/saber/etc/agent.yaml", "")
	rootCmd.CompletionOptions.DisableDefaultCmd = true
	rootCmd.AddCommand(agent.VersionCmd)
	if err := rootCmd.Execute(); err != nil {
		logger.Error("failed to start agent server. errmsg:%s", err.Error())
		return
	}

}
