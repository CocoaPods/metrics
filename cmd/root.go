package cmd

import (
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"

	"github.com/cocoapods/metrics/cmd/database"
)

var RootCmd = &cobra.Command{
	Use:   "metrics",
	Short: "CocoaPods Metrics Utilities",
}

func init() {
	RootCmd.AddCommand(VersionCmd)
	RootCmd.AddCommand(database.DatabaseCommand)
}

func Execute() {
	if err := RootCmd.Execute(); err != nil {
		logrus.Fatal(err.Error())
	}
}
