package webProperty_model

func (w WebPropertyTypes) ToMap() map[int]WebPropertyType {
	theMap := make(map[int]WebPropertyType)
	for _, v := range w {
		theMap[v.ID] = v
	}
	return theMap
}

func (w WebPropertyNotes) ToMap() map[int]WebPropertyNote {
	theMap := make(map[int]WebPropertyNote)
	for _, v := range w {
		theMap[v.WebPropID] = v
	}
	return theMap
}

func (w WebPropertyRequirements) ToMap() map[int]WebPropertyRequirement {
	theMap := make(map[int]WebPropertyRequirement)
	for _, v := range w {
		theMap[v.WebPropID] = v
	}
	return theMap
}
