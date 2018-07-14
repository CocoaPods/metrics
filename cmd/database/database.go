package database

import (
	"github.com/spf13/cobra"

	"github.com/cocoapods/metrics/aggregator"
	"github.com/cocoapods/metrics/internal/config"
)

var DatabaseCommand = &cobra.Command{
	Use:   "database",
	Short: "A collection of commands for performing database actions",
}

func newAggregateDataCommand() *cobra.Command {
	return &cobra.Command{
		Use:   "aggregate",
		Short: "Rollup data from the warehouse database into the metrics db",
		RunE: func(cmd *cobra.Command, args []string) error {
			c, err := config.Parse("metrics", []string{".", "/etc/cocoapods-metrics"})
			if err != nil {
				return err
			}
			a := aggregator.NewAggregator(c)
			return a.Aggregate()
		},
	}
}

func init() {
	DatabaseCommand.AddCommand(newAggregateDataCommand())
}
