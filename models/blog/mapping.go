package blog_model

func (cs Categories) ToMap() map[interface{}]Category {
	catMap := make(map[interface{}]Category)
	for _, v := range cs {
		catMap[v.ID] = v
	}
	return catMap
}

func (bcs BlogCategories) ToMap() map[interface{}]BlogCategory {
	bcMap := make(map[interface{}]BlogCategory)
	for _, v := range bcs {
		bcMap[v.ID] = v
	}
	return bcMap
}
