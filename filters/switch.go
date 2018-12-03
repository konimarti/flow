package filters

type switchElem struct {
	filters []Filter
	value   interface{}
}

func (s *switchElem) Check(v interface{}) bool {
	for _, f := range s.filters {
		if f.Check(v) {
			s.value = f.Update(v)
			return true
		}
	}
	return false
}

func (s *switchElem) Update(v interface{}) interface{} {
	return s.value
}

//NewSwitch accepts a list of filters and returns Switch Filter.
//The Switch Filter evaluates all filters in sequence and
//returns true if any of the Filters is true.
func NewSwitch(fs ...Filter) Filter {
	filters := []Filter{}
	for _, f := range fs {
		filters = append(filters, f)
	}
	return &switchElem{filters: filters}
}
