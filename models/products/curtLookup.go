package products

import (
	"database/sql"
	"strconv"
	"strings"

	"github.com/curt-labs/API/helpers/apicontext"
	"github.com/curt-labs/API/helpers/database"
)

var (
	GetYearsStmt = `select distinct year, '' from Vehicle v
					join Year as y on v.yearID = y.yearID
					join VehiclePart vp on v.vehicleID = vp.vehicleID
					join Part as p on vp.partID = p.partID
					where p.status = 800 || p.status = 900
					order by year desc`
	GetMakesStmt = `select distinct ma.make, group_concat(p.partID) from Make as ma
						join Vehicle as v on ma.makeID = v.makeID
						join Year as y on v.yearID = y.yearID
						join VehiclePart as vp on v.vehicleID = vp.vehicleID
						join Part as p on vp.partID = p.partID
						where y.year = ? && (p.status = 800 || p.status = 900) && p.classID > ?
						group by ma.make
						order by ma.make`
	GetModelsStmt = `select distinct mo.model, group_concat(p.partID) from Model as mo
						join Vehicle as v on mo.modelID = v.modelID
						join Year as y on v.yearID = y.yearID
						join Make as ma on v.makeID = ma.makeID
						join VehiclePart as vp on v.vehicleID = vp.vehicleID
						join Part as p on vp.partID = p.partID
						where y.year = ? && ma.make = ? && (p.status = 800 || p.status = 900) && p.classID > ?
						group by mo.model
						order by mo.model`
	GetStylesStmt = `select distinct s.style, group_concat(p.partID) from Style as s
						join Vehicle as v on s.styleID = v.styleID
						join Year as y on v.yearID = y.yearID
						join Make as ma on v.makeID = ma.makeID
						join Model as mo on v.modelID = mo.modelID
						join VehiclePart as vp on v.vehicleID = vp.vehicleID
						join Part as p on vp.partID = p.partID
						where y.year = ? && ma.make = ? && mo.model = ? && (p.status = 800 || p.status = 900) && p.classID > ?
						group by s.style
						order by s.style`
	GetPartNumbersStmt = `select distinct p.partID from Style as s
							join Vehicle as v on s.styleID = v.styleID
							join Year as y on v.yearID = y.yearID
							join Make as ma on v.makeID = ma.makeID
							join Model as mo on v.modelID = mo.modelID
							join VehiclePart as vp on v.vehicleID = vp.vehicleID
							join Part as p on vp.partID = p.partID
							where y.year = ? && ma.make = ? && mo.model = ? && s.style = ? && (p.status = 800 || p.status = 900) && p.classID > ?
							order by p.partID`
	GetPartNumbersWithoutStyleStmt = `select distinct p.partID from Style as s
							join Vehicle as v on s.styleID = v.styleID
							join Year as y on v.yearID = y.yearID
							join Make as ma on v.makeID = ma.makeID
							join Model as mo on v.modelID = mo.modelID
							join VehiclePart as vp on v.vehicleID = vp.vehicleID
							join Part as p on vp.partID = p.partID
							where y.year = ? && ma.make = ? && mo.model = ? && (p.status = 800 || p.status = 900) && p.classID > ?
							order by p.partID`
)

type CurtVehicle struct {
	Year            string      `json:"year,omitempty" xml:"year, omitempty"`
	Make            string      `json:"make,omitempty" xml:"make, omitempty"`
	Model           string      `json:"model,omitempty" xml:"model, omitempty"`
	Style           string      `json:"style,omitempty" xml:"style, omitempty"`
	Parts           []BasicPart `json:"parts,omitempty" xml:"parts, omitempty"`
	PartIdentifiers []int       `json:"parts_ids" xml:"-"`
}

type CurtLookup struct {
	Years  []string `json:"available_years,omitempty" xml:"available_years, omitempty"`
	Makes  []string `json:"available_makes,omitempty" xml:"available_makes, omitempty"`
	Models []string `json:"available_models,omitempty" xml:"available_models, omitempty"`
	Styles []string `json:"available_styles,omitempty" xml:"available_styles, omitempty"`
	Parts  []Part   `json:"parts,omitempty" xml:"parts, omitempty"`
	CurtVehicle
}

func (c *CurtLookup) GetYears(heavyduty bool) error {

	db, err := sql.Open("mysql", database.ConnectionString())
	if err != nil {
		return err
	}
	defer db.Close()

	stmt, err := db.Prepare(GetYearsStmt)
	if err != nil {
		return err
	}
	defer stmt.Close()
	rows, err := stmt.Query()
	if err != nil {
		return err
	}
	defer rows.Close()

	var ids []string
	for rows.Next() {
		var val string
		var idStr string
		if err = rows.Scan(&val, &idStr); err != nil {
			continue
		}

		arr := strings.Split(idStr, ",")
		if err == nil && len(arr) > 0 {
			ids = append(ids, arr...)
		}

		c.Years = append(c.Years, val)

	}

	existing := make(map[int]int, 0)
	for _, i := range ids {
		intID, err := strconv.Atoi(i)
		if err == nil {
			if _, ok := existing[intID]; !ok {
				c.PartIdentifiers = append(c.PartIdentifiers, intID)
				existing[intID] = intID
			}
		}
	}

	return rows.Err()
}

func (c *CurtLookup) GetMakes(heavyduty bool) error {

	db, err := sql.Open("mysql", database.ConnectionString())
	if err != nil {
		return err
	}
	defer db.Close()

	stmt, err := db.Prepare(GetMakesStmt)
	if err != nil {
		return err
	}
	defer stmt.Close()
	class := 0
	if heavyduty {
		class = -1
	}
	rows, err := stmt.Query(c.Year, class)
	if err != nil {
		return err
	}
	defer rows.Close()

	var ids []string
	for rows.Next() {
		var val string
		var idStr string
		if err = rows.Scan(&val, &idStr); err != nil {
			continue
		}

		arr := strings.Split(idStr, ",")
		if err == nil && len(arr) > 0 {
			ids = append(ids, arr...)
		}

		c.Makes = append(c.Makes, val)

	}

	existing := make(map[int]int, 0)
	for _, i := range ids {
		intID, err := strconv.Atoi(i)
		if err == nil {
			if _, ok := existing[intID]; !ok {
				c.PartIdentifiers = append(c.PartIdentifiers, intID)
				existing[intID] = intID
			}
		}
	}

	return rows.Err()
}

func (c *CurtLookup) GetModels(heavyduty bool) error {

	db, err := sql.Open("mysql", database.ConnectionString())
	if err != nil {
		return err
	}
	defer db.Close()

	stmt, err := db.Prepare(GetModelsStmt)
	if err != nil {
		return err
	}
	defer stmt.Close()
	class := 0
	if heavyduty {
		class = -1
	}
	rows, err := stmt.Query(c.Year, c.Make, class)
	if err != nil {
		return err
	}
	defer rows.Close()

	var ids []string
	for rows.Next() {
		var val string
		var idStr string
		if err = rows.Scan(&val, &idStr); err != nil {
			continue
		}

		arr := strings.Split(idStr, ",")
		if err == nil && len(arr) > 0 {
			ids = append(ids, arr...)
		}

		c.Models = append(c.Models, val)

	}

	existing := make(map[int]int, 0)
	for _, i := range ids {
		intID, err := strconv.Atoi(i)
		if err == nil {
			if _, ok := existing[intID]; !ok {
				c.PartIdentifiers = append(c.PartIdentifiers, intID)
				existing[intID] = intID
			}
		}
	}

	return rows.Err()
}

func (c *CurtLookup) GetStyles(heavyduty bool) error {

	db, err := sql.Open("mysql", database.ConnectionString())
	if err != nil {
		return err
	}
	defer db.Close()

	stmt, err := db.Prepare(GetStylesStmt)
	if err != nil {
		return err
	}
	defer stmt.Close()
	class := 0
	if heavyduty {
		class = -1
	}
	rows, err := stmt.Query(c.Year, c.Make, c.Model, class)
	if err != nil {
		return err
	}
	defer rows.Close()

	var ids []string
	for rows.Next() {
		var val string
		var idStr string
		if err = rows.Scan(&val, &idStr); err != nil {
			continue
		}

		arr := strings.Split(idStr, ",")
		if err == nil && len(arr) > 0 {
			ids = append(ids, arr...)
		}

		c.Styles = append(c.Styles, val)

	}

	existing := make(map[int]int, 0)
	for _, i := range ids {
		intID, err := strconv.Atoi(i)
		if err == nil {
			if _, ok := existing[intID]; !ok {
				c.PartIdentifiers = append(c.PartIdentifiers, intID)
				existing[intID] = intID
			}
		}
	}

	return rows.Err()
}

func (c *CurtLookup) GetParts(dtx *apicontext.DataContext, heavyduty bool) error {
	db, err := sql.Open("mysql", database.ConnectionString())
	if err != nil {
		return err
	}
	defer db.Close()
	var rows *sql.Rows
	class := 0
	if heavyduty {
		class = -1
	}

	if c.Style != "" {
		rows, err = db.Query(GetPartNumbersStmt, c.Year, c.Make, c.Model, c.Style, class)
		if err != nil {
			return err
		}
		defer rows.Close()
	} else {
		rows, err = db.Query(GetPartNumbersWithoutStyleStmt, c.Year, c.Make, c.Model, class)
		if err != nil {
			return err
		}
		defer rows.Close()
	}
	ch := make(chan *Part)
	iter := 0
	for rows.Next() {
		iter++
		var p Part
		if err = rows.Scan(&p.ID); err != nil {
			ch <- nil
			continue
		}
		go func(prt Part) {
			er := prt.Get(dtx)
			if er != nil {
				ch <- nil
			} else {
				ch <- &prt
			}

		}(p)

	}

	for i := 0; i < iter; i++ {
		tmp := <-ch
		if tmp != nil {
			c.Parts = append(c.Parts, *tmp)
		}
	}
	return rows.Err()
}
