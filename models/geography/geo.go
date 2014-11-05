package geography

import (
	"database/sql"
	"github.com/curt-labs/GoAPI/helpers/database"
	"github.com/curt-labs/GoAPI/helpers/sortutil"
	_ "github.com/go-sql-driver/mysql"
)

var (
	getAllStatesStmt             = `select * from States`
	getAllCountriesStmt          = `select * from Country`
	getAllCountriesAndStatesStmt = `select C.*, S.stateID, S.state, S.abbr from Country C
									inner join States S on S.countryID = C.countryID
									order by C.countryID, S.state`
)

type States []State
type State struct {
	Id           int      `json:"state_id"`
	State        string   `json:"state"`
	Abbreviation string   `json:"abbreviation"`
	Country      *Country `json:"country,omitempty"`
}

type Countries []Country
type Country struct {
	Id           int     `json:"country_id"`
	Country      string  `json:"country"`
	Abbreviation string  `json:"abbreviation"`
	States       *States `json:"states,omitempty"`
}

func GetAllCountriesAndStates() (countries Countries, err error) {
	db, err := sql.Open("mysql", database.ConnectionString())
	if err != nil {
		return
	}
	defer db.Close()

	stmt, err := db.Prepare(getAllCountriesAndStatesStmt)
	if err != nil {
		return
	}
	defer stmt.Close()

	rows, err := stmt.Query()
	if err != nil {
		return
	}

	countryMap := make(map[int]Country, 0)

	for rows.Next() {
		var c Country
		var s State

		err = rows.Scan(
			&c.Id,
			&c.Country,
			&c.Abbreviation,
			&s.Id,
			&s.State,
			&s.Abbreviation,
		)
		if err != nil {
			return
		}

		country, exists := countryMap[c.Id]

		if !exists {
			c.States = &States{s}
			countryMap[c.Id] = c
		} else {
			*country.States = append(*country.States, s)
		}
	}
	defer rows.Close()

	for _, c := range countryMap {
		countries = append(countries, c)
	}

	sortutil.AscByField(countries, "Id")
	return
}

func GetAllCountries() (countries Countries, err error) {
	db, err := sql.Open("mysql", database.ConnectionString())
	if err != nil {
		return
	}
	defer db.Close()

	stmt, err := db.Prepare(getAllCountriesStmt)
	if err != nil {
		return
	}
	defer stmt.Close()

	rows, err := stmt.Query()
	if err != nil {
		return
	}

	for rows.Next() {
		var c Country
		err = rows.Scan(
			&c.Id,
			&c.Country,
			&c.Abbreviation,
		)
		if err != nil {
			return
		}
		countries = append(countries, c)
	}
	defer rows.Close()

	sortutil.AscByField(countries, "Id")

	return
}

func GetAllStates() (states States, err error) {
	db, err := sql.Open("mysql", database.ConnectionString())
	if err != nil {
		return
	}
	defer db.Close()

	stmt, err := db.Prepare(getAllStatesStmt)
	if err != nil {
		return
	}
	defer stmt.Close()

	rows, err := stmt.Query()
	if err != nil {
		return
	}

	for rows.Next() {
		var state State
		state.Country = &Country{}
		err = rows.Scan(
			&state.Id,
			&state.State,
			&state.Abbreviation,
			&state.Country.Id,
		)
		if err != nil {
			return
		}
		states = append(states, state)
	}
	defer rows.Close()

	return
}

func GetStateMap() (map[int]State, error) {
	stateMap := make(map[int]State)
	states, err := GetAllStates()
	for _, state := range states {
		stateMap[state.Id] = state
	}
	return stateMap, err
}

func GetCountryMap() (map[int]Country, error) {
	countryMap := make(map[int]Country)
	countries, err := GetAllCountries()
	for _, country := range countries {
		countryMap[country.Id] = country
	}
	return countryMap, err
}
