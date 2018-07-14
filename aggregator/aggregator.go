package aggregator

type Aggregator struct {
}

func NewAggregator() *Aggregator {
	return &Aggregator{}
}

func (a *Aggregator) Aggregate() error {
	return nil
}
