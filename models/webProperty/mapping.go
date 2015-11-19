package webProperty_model

// creates a map so you can more easily fetch WebPropertyType's by their ID
func (w WebPropertyTypes) ToMap() map[int]WebPropertyType {
	theMap := make(map[int]WebPropertyType)
	for _, v := range w {
		theMap[v.ID] = v
	}
	return theMap
}

// creates a map so you can more easily fetch WebPropertyNote's by their ID
func (w WebPropertyNotes) ToMap() map[int]WebPropertyNote {
	theMap := make(map[int]WebPropertyNote)
	for _, v := range w {
		theMap[v.WebPropID] = v
	}
	return theMap
}

// creates a map so you can more easily fetch WebPropertyRequirement's by their ID
func (w WebPropertyRequirements) ToMap() map[int]WebPropertyRequirement {
	theMap := make(map[int]WebPropertyRequirement)
	for _, v := range w {
		theMap[v.WebPropID] = v
	}
	return theMap
}
