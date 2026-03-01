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

package admin

import (
	"os-artificer/saber/pkg/version"

	"github.com/spf13/cobra"
)

var VersionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print Version Information",
	Run: func(cmd *cobra.Command, args []string) {
		version.Print("saber Admin Server")
	},
}

// StartCmd starts the admin server
var StartCmd = &cobra.Command{
	Use:   "start",
	Short: "Start Admin Server",
	RunE:  runStart,
}

// StopCmd stops the admin server
var StopCmd = &cobra.Command{
	Use:   "stop",
	Short: "Stop Admin Server",
	RunE:  runStop,
}

// RestartCmd restarts the admin server
var RestartCmd = &cobra.Command{
	Use:   "restart",
	Short: "Restart Admin Server",
	RunE:  runRestart,
}

// ReloadCmd reloads the admin server config
var ReloadCmd = &cobra.Command{
	Use:   "reload",
	Short: "Reload Admin Server",
	RunE:  runReload,
}

var HealthCheckCmd = &cobra.Command{
	Use:   "health",
	Short: "Health Check Admin Server",
	Run: func(cmd *cobra.Command, args []string) {
	},
}
