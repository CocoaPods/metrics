package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var (
	Version   = "dev"
	BuildDate = ""
)

var VersionCmd = &cobra.Command{
	Use:   "version",
	Short: "show version",
	RunE:  versionRunE,
}

func versionRunE(_ *cobra.Command, _ []string) error {
	fmt.Printf("Metrics version: %s\nbuilt: %s\n", Version, BuildDate)
	return nil
}
