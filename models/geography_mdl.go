package models

type State struct {
	State, Abbreviation string
	Country             *Country
}

type Country struct {
	Country, Abbreviation string
}
