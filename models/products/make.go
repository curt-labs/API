package products

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"github.com/curt-labs/API/helpers/apicontext"
	"github.com/curt-labs/API/helpers/database"
	"github.com/curt-labs/API/helpers/redis"
	_ "github.com/go-sql-driver/mysql"
	"strconv"
	"strings"
)

func (l *Lookup) GetMakes(dtx *apicontext.DataContext) error {
	redis_key := fmt.Sprintf("lookup:year:%d:makes:%s", l.Vehicle.Base.Year, dtx.BrandString)
	data, err := redis.Get(redis_key)
	if err == nil {
		err = json.Unmarshal(data, &l.Makes)
		return nil
	}

	stmtBeginning := `
		select distinct m.MakeName from vcdb_Make as m
		join BaseVehicle as bv on m.ID = bv.MakeID
		join vcdb_Vehicle as v on bv.ID = v.BaseVehicleID
		join vcdb_VehiclePart as vp on v.ID = vp.VehicleID
		join Part as p on vp.PartNumber = p.partID
		where (p.status = 800 || p.status = 900) && bv.YearID = ? `
	stmtEnd := `	order by m.MakeName`
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

	res, err := stmt.Query(l.Vehicle.Base.Year)
	if err != nil {
		return err
	}

	l.Makes = make([]string, 0)
	for res.Next() {
		var ma string
		err = res.Scan(&ma)
		if err == nil {
			l.Makes = append(l.Makes, ma)
		}
	}
	defer res.Close()

	l.Pagination = Pagination{
		TotalItems:    len(l.Makes),
		ReturnedCount: len(l.Makes),
		Page:          1,
		PerPage:       len(l.Makes),
		TotalPages:    1,
	}
	if dtx.BrandString != "" {
		redis.Setex(redis_key, l.Makes, 86400)
	}
	return nil
}
