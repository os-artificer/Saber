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

	"os-artificer/saber/internal/admin"
	"os-artificer/saber/pkg/logger"

	"github.com/spf13/cobra"
)

func main() {
	if os.Getenv("SABER_ADMIN_SUPERVISOR") == "1" {
		if err := admin.RunSupervisor(); err != nil {
			os.Exit(1)
		}
		os.Exit(0)
		return
	}

	rootCmd := &cobra.Command{
		Use:          "Admin",
		Short:        "Saber Admin Server",
		SilenceUsage: true,
		RunE:         admin.Run,
	}

	rootCmd.PersistentFlags().StringVarP(&admin.ConfigFilePath, "config", "c", "./etc/admin.yaml", "")
	rootCmd.CompletionOptions.DisableDefaultCmd = true

	rootCmd.AddCommand(admin.StartCmd)
	rootCmd.AddCommand(admin.StopCmd)
	rootCmd.AddCommand(admin.RestartCmd)
	rootCmd.AddCommand(admin.ReloadCmd)
	rootCmd.AddCommand(admin.HealthCheckCmd)
	rootCmd.AddCommand(admin.VersionCmd)
	rootCmd.AddCommand(admin.MigrateCmd)

	if err := rootCmd.Execute(); err != nil {
		logger.Errorf("failed to start admin server. errmsg:%s", err.Error())
		os.Exit(1)
	}
}
