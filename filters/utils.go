package filters

//GetFloat64 parses interface to float64, if not a number it returns 0
func GetFloat64(v interface{}) float64 {
	var ret float64
	switch v.(type) {
	case int:
		ret = float64(v.(int))
	case int16:
		ret = float64(v.(int16))
	case int32:
		ret = float64(v.(int32))
	case int64:
		ret = float64(v.(int64))
	case float32:
		ret = float64(v.(float32))
	case float64:
		ret = float64(v.(float64))
	}
	return ret

}
