package models

type State struct {
	Id                  int
	State, Abbreviation string
	Country             *Country
}

type Country struct {
	Id                    int
	Country, Abbreviation string
}
