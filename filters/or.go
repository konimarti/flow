package filters

type or struct {
	left  Filter
	right Filter
	value interface{}
}

func (o *or) Check(v interface{}) bool {
	if o.left.Check(v) {
		o.value = o.left.Update(v)
		return true
	} else if o.right.Check(v) {
		o.value = o.right.Update(v)
		return true
	} else {
		return false
	}
}

func (o *or) Update(v interface{}) interface{} {
	return o.value
}

//NewOr evaluates the returned check value for
//two filters and performs the operation of the first
//filter that returns true.
func NewOr(left, right Filter) Filter {
	return &or{left: left, right: right}
}
