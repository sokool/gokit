package cqrs

import "github.com/sokool/gokit/log"

type cache struct {
	aggregates map[string]Aggregate
}

func (c *cache) restore(id string) (Aggregate, bool) {
	a, ok := c.aggregates[id]
	if ok {
		log.Info("cqrs.cache", "%s restored", a.Root().String())
	}

	return a, ok
}

func (c *cache) store(a Aggregate) {
	b := c.aggregates[a.Root().ID]
	if b != nil {
		return
	}
	log.Info("cqrs.cache", "%s stored", a.Root().String())

	c.aggregates[a.Root().ID] = a
}

func newCache() *cache {
	return &cache{
		aggregates: map[string]Aggregate{},
	}
}
