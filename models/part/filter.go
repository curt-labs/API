package part

type FilterSet struct {
	Name    string
	Options []FilterOption
}

type FilterOption struct {
	Name     string
	Id       int
	Selected bool
}

func (p Part) FilterAttributes() (sets map[string]FilterSet, err error) {
	sets = make(map[string]FilterSet, 0)

	for _, attr := range p.Attributes {

		if _, ok := sets[attr.Key]; !ok {
			fs := FilterSet{
				Name:    attr.Key,
				Options: make([]FilterOption, 0),
			}
			sets[fs.Name] = fs
		}
	}
	return
}
