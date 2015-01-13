package products

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"github.com/curt-labs/GoAPI/helpers/apicontext"
	"github.com/curt-labs/GoAPI/helpers/database"
	"github.com/curt-labs/GoAPI/helpers/redis"
	_ "github.com/go-sql-driver/mysql"
	"strconv"
	"strings"
)

func (l *Lookup) GetYears(dtx *apicontext.DataContext) error {
	//hit redis first
	var brands string
	if dtx.Globals["brandsString"] != nil {
		brands = dtx.Globals["brandsString"].(string)
	}
	redis_key := fmt.Sprintf("lookup:years:%s", brands)
	data, err := redis.Get(redis_key)
	if err == nil {
		err = json.Unmarshal(data, &l.Years)
		return nil
	}

	stmtBeginning := `
		select distinct y.YearID from vcdb_Year as y
		join BaseVehicle as bv on y.YearID = bv.YearID
		join vcdb_Vehicle as v on bv.ID = v.BaseVehicleID
		join vcdb_VehiclePart as vp on v.ID = vp.VehicleID
		join Part as p on vp.PartNumber = p.partID
		where (p.status = 800 || p.status = 900) `
	stmtEnd := ` order by y.YearID desc`
	brandStmt := " && p.brandID in ("

	for _, b := range l.Brands {
		brandStmt += strconv.Itoa(b) + ","
	}
	brandStmt = strings.TrimRight(brandStmt, ",") + ")"
	wholeStmt := stmtBeginning + brandStmt + stmtEnd

	db, err := sql.Open("mysql", database.ConnectionString())
	if err != nil {
		return err
	}
	defer db.Close()

	stmt, err := db.Prepare(wholeStmt)
	if err != nil {
		return err
	}
	defer stmt.Close()

	res, err := stmt.Query()
	if err != nil {
		return err
	}

	l.Years = make([]int, 0)
	for res.Next() {
		var year int
		err = res.Scan(&year)
		if err == nil {
			l.Years = append(l.Years, year)
		}
	}
	defer res.Close()

	l.Pagination = Pagination{
		TotalItems:    len(l.Years),
		ReturnedCount: len(l.Years),
		Page:          1,
		PerPage:       len(l.Years),
		TotalPages:    1,
	}

	go redis.Setex(redis_key, l.Years, 86400)

	return nil
}
