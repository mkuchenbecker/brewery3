package counters

type Counter struct {
	counters map[string]uint64
}

func New() *Counter {
	return &Counter{counters: make(map[string]uint64)}
}

func (c *Counter) Inc(key string) {
	c.counters[key]++
}

func (c *Counter) Get(key string) uint64 {
	count, ok := c.counters[key]
	if !ok {
		return 0
	}
	return count
}
