package main

import (
	"os-artificer/saber/internal/probe"
	"os-artificer/saber/pkg/logger"

	"github.com/spf13/cobra"
)

func main() {
	rootCmd := &cobra.Command{
		Use:          "Probe",
		Short:        "saber Probe",
		SilenceUsage: true,
		RunE:         probe.Run,
	}

	rootCmd.PersistentFlags().StringVarP(&probe.ConfigFilePath, "config", "c", "/var/ylg/saber/etc/probe.yaml", "")
	rootCmd.CompletionOptions.DisableDefaultCmd = true
	rootCmd.AddCommand(probe.VersionCmd)
	if err := rootCmd.Execute(); err != nil {
		logger.Error("failed to start probe server. errmsg:%s", err.Error())
		return
	}
}
