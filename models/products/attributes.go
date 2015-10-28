package products

type Attribute struct {
	Key   string `json:"key" xml:"key,attr" bson:"key"`
	Value string `json:"value" xml:",chardata" bson:"value"`
	Sort  int    `json:"sort,omitempty" xml:"sort,omitempty" bson:"sort"`
}
