package geography

type State struct {
	Id                  int
	State, Abbreviation string
	Country             *Country
}

type Country struct {
	Id                    int
	Country, Abbreviation string
}

type State_New struct {
	Id                  int
	State, Abbreviation string
	Country             Country
}
