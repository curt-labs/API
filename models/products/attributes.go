package products

type Attribute struct {
	Key   string `json:"key" xml:"key,attr"`
	Value string `json:"value" xml:",chardata"`
	Sort  int    `json:"sort,omitempty" xml:"sort,omitempty"`
}
