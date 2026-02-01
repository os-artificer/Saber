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
		logger.Errorf("failed to start agent server. errmsg:%s", err.Error())
		return
	}
}
