package filters

type chain struct {
	fs    []Filter
	Value interface{}
	Flag  bool
}

func (c *chain) Check(v interface{}) bool {
	c.Value = v
	c.Flag = true
	for _, f := range c.fs {
		if f.Check(c.Value) {
			c.Value = f.Update(c.Value)
		} else {
			c.Flag = false
			break
		}
	}
	return c.Flag
}

func (c *chain) Update(v interface{}) interface{} {
	return c.Value
}

//NewChain chains together filters.
func NewChain(filters ...Filter) Filter {
	chainedFilters := make([]Filter, 0)
	for _, f := range filters {
		chainedFilters = append(chainedFilters, f)
	}
	return &chain{fs: chainedFilters}
}
