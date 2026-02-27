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

package main

import (
	"os"

	"os-artificer/saber/internal/databus"
	"os-artificer/saber/pkg/logger"

	"github.com/spf13/cobra"
)

func main() {
	if os.Getenv("SABER_DATABUS_SUPERVISOR") == "1" {
		if err := databus.RunSupervisor(); err != nil {
			os.Exit(1)
		}
		os.Exit(0)
		return
	}

	rootCmd := &cobra.Command{
		Use:          "Databus",
		Short:        "Saber Databus Server",
		SilenceUsage: true,
		RunE:         databus.Run,
	}

	rootCmd.PersistentFlags().StringVarP(&databus.ConfigFilePath, "config", "c", "./etc/databus.yaml", "")
	rootCmd.CompletionOptions.DisableDefaultCmd = true

	rootCmd.AddCommand(databus.StartCmd)
	rootCmd.AddCommand(databus.StopCmd)
	rootCmd.AddCommand(databus.RestartCmd)
	rootCmd.AddCommand(databus.ReloadCmd)
	rootCmd.AddCommand(databus.HealthCheckCmd)
	rootCmd.AddCommand(databus.VersionCmd)

	if err := rootCmd.Execute(); err != nil {
		logger.Errorf("failed to start databus server. errmsg:%s", err.Error())
		return
	}

}
