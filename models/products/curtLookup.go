package products

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"strconv"
	"strings"
	"time"

	"github.com/curt-labs/API/helpers/apicontext"
	"github.com/curt-labs/API/helpers/database"
	"github.com/curt-labs/API/helpers/redisNew"
)

var (
	GetYearsStmt = `select distinct year, '' from Year
						order by year desc`
	GetMakesStmt = `select distinct ma.make, group_concat(p.partID) from Make as ma
						join Vehicle as v on ma.makeID = v.makeID
						join Year as y on v.yearID = y.yearID
						join VehiclePart as vp on v.vehicleID = vp.vehicleID
						join Part as p on vp.partID = p.partID
						where y.year = ? && (p.status = 800 || p.status = 900)
						group by ma.make
						order by ma.make`
	GetModelsStmt = `select distinct mo.model, group_concat(p.partID) from Model as mo
						join Vehicle as v on mo.modelID = v.modelID
						join Year as y on v.yearID = y.yearID
						join Make as ma on v.makeID = ma.makeID
						join VehiclePart as vp on v.vehicleID = vp.vehicleID
						join Part as p on vp.partID = p.partID
						where y.year = ? && ma.make = ? && (p.status = 800 || p.status = 900)
						group by mo.model
						order by mo.model`
	GetStylesStmt = `select distinct s.style, group_concat(p.partID) from Style as s
						join Vehicle as v on s.styleID = v.styleID
						join Year as y on v.yearID = y.yearID
						join Make as ma on v.makeID = ma.makeID
						join Model as mo on v.modelID = mo.modelID
						join VehiclePart as vp on v.vehicleID = vp.vehicleID
						join Part as p on vp.partID = p.partID
						where y.year = ? && ma.make = ? && mo.model = ? && (p.status = 800 || p.status = 900)
						group by s.style
						order by s.style`
	GetPartNumbersStmt = `select distinct p.partID from Style as s
							join Vehicle as v on s.styleID = v.styleID
							join Year as y on v.yearID = y.yearID
							join Make as ma on v.makeID = ma.makeID
							join Model as mo on v.modelID = mo.modelID
							join VehiclePart as vp on v.vehicleID = vp.vehicleID
							join Part as p on vp.partID = p.partID
							where y.year = ? && ma.make = ? && mo.model = ? && s.style = ? && (p.status = 800 || p.status = 900)
							order by p.partID`
	GetPartNumbersWithoutStyleStmt = `select distinct p.partID from Style as s
							join Vehicle as v on s.styleID = v.styleID
							join Year as y on v.yearID = y.yearID
							join Make as ma on v.makeID = ma.makeID
							join Model as mo on v.modelID = mo.modelID
							join VehiclePart as vp on v.vehicleID = vp.vehicleID
							join Part as p on vp.partID = p.partID
							where y.year = ? && ma.make = ? && mo.model = ? && (p.status = 800 || p.status = 900)
							order by p.partID`
	VehicleAppsStmt = `select v.vehicleID, y.year, ma.make, mo.model, s.style, p.partID, p.shortDesc, p.dateModified
									from Part as p
									join VehiclePart as vp on p.partID = vp.partID
									join Vehicle as v on vp.vehicleID =  v.vehicleID
									join Year as y on v.yearID = y.yearID
									join Make as ma on v.makeID = ma.makeID
									join Model as mo on v.modelID = mo.modelID
									join Style as s on v.styleID = s.styleID
									where p.status != 999 && p.brandID = 1
									order by p.partID, ma.make, mo.model, s.style, y.year`
	VehicleAppsWithDate = `select v.vehicleID, y.year, ma.make, mo.model, s.style, p.partID, p.shortDesc, p.dateModified
										from Part as p
										join VehiclePart as vp on p.partID = vp.partID
										join Vehicle as v on vp.vehicleID =  v.vehicleID
										join Year as y on v.yearID = y.yearID
										join Make as ma on v.makeID = ma.makeID
										join Model as mo on v.modelID = mo.modelID
										join Style as s on v.styleID = s.styleID
										where p.status != 999 && p.brandID = 1 && p.dateModified >= ?
										order by p.partID, ma.make, mo.model, s.style, y.year`
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

type VehicleApp struct {
	Category      string
	VehicleID     int
	Year          float64
	Make          string
	Model         string
	Style         string
	PartID        int
	PartShortDesc string
	DateModified  time.Time
	GroupID       string
	Drilling      string
	Exposed       string
	InstallTime   int
	ClassID       int
	PartNumber    string
}

func CurtVehicleApps(date string) (vehicleApps []VehicleApp, err error) {
	vehicleApps = make([]VehicleApp, 0)
	redis_key := fmt.Sprintf("CurtVehicleApps:v4:%s", date)
	data, err := redis.RedisMaster.Get(redis_key)
	log.Println()
	if err == nil {
		err = json.Unmarshal(data, &vehicleApps)

		if err != nil {
			log.Println("error unmarshaling", err)
			log.Println(len(data))
		}
		log.Print("returning from redis")
		log.Println(len(vehicleApps))
		return vehicleApps, err
	} else {
		log.Println("REDIS ERROR Getting Vehicle apps:", err)
	}

	var stmt *sql.Stmt
	var res *sql.Rows

	db, err := sql.Open("mysql", database.ConnectionString())
	if err != nil {
		return vehicleApps, err
	}
	defer db.Close()

	if date == "" {
		stmt, err = db.Prepare(VehicleAppsStmt)
		if err != nil {
			return vehicleApps, err
		}
		defer stmt.Close()
		res, err = stmt.Query()
		if err != nil {
			return vehicleApps, err
		}
	} else {
		stmt, err = db.Prepare(VehicleAppsWithDate)
		if err != nil {
			return vehicleApps, err
		}
		defer stmt.Close()
		res, err = stmt.Query(date)
		if err != nil {
			return vehicleApps, err
		}
	}
	defer res.Close()

	for res.Next() {
		var v VehicleApp
		var dateMod *time.Time
		var year, ma, mo, st, desc *string
		var vID, partID *int
		err = res.Scan(
			&vID,
			&year,
			&ma,
			&mo,
			&st,
			&partID,
			&desc,
			&dateMod,
		)

		if err != nil {
			return vehicleApps, err
		}

		if vID != nil {
			v.VehicleID = *vID
		}
		if year != nil {
			v.Year, _ = strconv.ParseFloat(*year, 64)
		}
		if ma != nil {
			v.Make = *ma
		}
		if mo != nil {
			v.Model = *mo
		}
		if st != nil {
			v.Style = *st
		}
		if partID != nil {
			v.PartID = *partID
		}
		if desc != nil {
			v.PartShortDesc = *desc
		}
		if dateMod != nil {
			v.DateModified = *dateMod
		}

		vehicleApps = append(vehicleApps, v)
	}
	log.Println("setting vehicle apps in redis")
	log.Println(len(vehicleApps))
	if data_bytes, err := json.Marshal(&vehicleApps); err == nil {
		err = redis.RedisMaster.Setex(redis_key, 86400, data_bytes)
		log.Println(err)
	}

	return vehicleApps, err
}

func (c *CurtLookup) GetYears() error {

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

func (c *CurtLookup) GetMakes() error {

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

	rows, err := stmt.Query(c.Year)
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

func (c *CurtLookup) GetModels() error {

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

	rows, err := stmt.Query(c.Year, c.Make)
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

func (c *CurtLookup) GetStyles() error {

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

	rows, err := stmt.Query(c.Year, c.Make, c.Model)
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

func (c *CurtLookup) GetParts(dtx *apicontext.DataContext) error {
	db, err := sql.Open("mysql", database.ConnectionString())
	if err != nil {
		return err
	}
	defer db.Close()
	var rows *sql.Rows

	if c.Style != "" {
		rows, err = db.Query(GetPartNumbersStmt, c.Year, c.Make, c.Model, c.Style)
		if err != nil {
			return err
		}
		defer rows.Close()
	} else {
		rows, err = db.Query(GetPartNumbersWithoutStyleStmt, c.Year, c.Make, c.Model)
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
