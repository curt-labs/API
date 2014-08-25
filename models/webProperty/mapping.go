package webProperty_model

func (w WebPropertyTypes) ToMap() map[interface{}]WebPropertyType {
	theMap := make(map[interface{}]WebPropertyType)
	for _, v := range w {
		theMap[v.ID] = v
	}
	return theMap
}

func (w WebPropertyNotes) ToMap() map[interface{}]WebPropertyNote {
	theMap := make(map[interface{}]WebPropertyNote)
	for _, v := range w {
		theMap[v.ID] = v
	}
	return theMap
}

func (w WebPropertyRequirements) ToMap() map[interface{}]WebPropertyRequirement {
	theMap := make(map[interface{}]WebPropertyRequirement)
	for _, v := range w {
		theMap[v.ID] = v
	}
	return theMap
}
