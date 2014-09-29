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
