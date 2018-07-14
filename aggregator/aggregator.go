package aggregator

import (
	"github.com/cocoapods/metrics/internal/config"
)

type Aggregator struct {
}

func NewAggregator(c *config.Config) *Aggregator {
	return &Aggregator{}
}

func (a *Aggregator) Aggregate() error {
	return nil
}
