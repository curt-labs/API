package lifestyle

//map contents
func (cs Contents) ToMap() map[interface{}]Content {
	zeeMap := make(map[interface{}]Content)
	for _, v := range cs {
		zeeMap[v.ID] = v
	}
	return zeeMap
}

//map towables
func (cs Towables) ToMap() map[interface{}]Towable {
	zeeMap := make(map[interface{}]Towable)
	for _, v := range cs {
		zeeMap[v.ID] = v
	}
	return zeeMap
}
