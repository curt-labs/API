package apifilter

type FilteredOptions []Options

type Options struct {
	Key     string
	Options []Option
}

type Option struct {
	Value    string
	Selected bool
	Products []int
}

type Decision struct {
	Field  string
	Values map[string]interface{}
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
