package products

type Attribute struct {
	Key   string `json:"name" xml:"name,attr" bson:"name"`
	Value string `json:"value" xml:",chardata" bson:"value"`
	Sort  int    `json:"sort,omitempty" xml:"sort,omitempty" bson:"sort"`
}
