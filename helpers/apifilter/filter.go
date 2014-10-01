package apifilter

type FilteredOptions []Options

type Options struct {
	Key     string   `json:"key" xml:"key,attr"`
	Options []Option `json:"options" xml:"options"`
}

type Option struct {
	Value    string `json:"value" xml:"value,attr"`
	Selected bool   `json:"selected" xml:"selected,attr"`
	Products []int  `json:"products" xml:"products"`
}

type Decision struct {
	Field  string                 `json:"field" xml:"field,attr"`
	Values map[string]interface{} `json:"values" xml:"values"`
}

func (o *Options) AppendValue(newOpt Option) {
	if len(o.Options) == 0 {
		o.Options = append(o.Options, newOpt)
		return
	}

	for i, opt := range o.Options {
		if opt.Value == newOpt.Value { // exists
			opt.Products = append(opt.Products, newOpt.Products...)
			opt.Products = removeDuplicates(opt.Products)
			o.Options[i] = opt
		} else {
			o.Options = append(o.Options, newOpt)
		}
	}
}

func removeDuplicates(a []int) []int {
	result := []int{}
	seen := map[int]int{}
	for _, val := range a {
		if _, ok := seen[val]; !ok {
			result = append(result, val)
			seen[val] = val
		}
	}
	return result
}
