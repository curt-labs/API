package products

import (
	"encoding/json"
	"fmt"
	"strconv"
	"strings"

	"github.com/curt-labs/API/helpers/apicontext"
	"github.com/curt-labs/API/helpers/database"
	"github.com/curt-labs/API/helpers/redis"
	_ "github.com/go-sql-driver/mysql"
)

func (l *Lookup) GetYears(dtx *apicontext.DataContext) error {
	//hit redis first
	redis_key := fmt.Sprintf("lookup:years:%s", dtx.BrandString)
	data, err := redis.Get(redis_key)
	if err == nil {
		err = json.Unmarshal(data, &l.Years)
		if len(l.Years) > 0 {
			return nil
		}
	}

	stmtBeginning := `
		select distinct y.YearID from vcdb_Year as y
		join BaseVehicle as bv on y.YearID = bv.YearID
		join vcdb_Vehicle as v on bv.ID = v.BaseVehicleID
		join vcdb_VehiclePart as vp on v.ID = vp.VehicleID
		join Part as p on vp.PartNumber = p.partID
		where p.status in (700, 800, 810, 815, 850, 870, 888, 900, 910, 950) `
	stmtEnd := ` order by y.YearID desc`
	brandStmt := " && p.brandID in ("

	for _, b := range l.Brands {
		brandStmt += strconv.Itoa(b) + ","
	}
	brandStmt = strings.TrimRight(brandStmt, ",") + ")"
	wholeStmt := stmtBeginning + brandStmt + stmtEnd

	err = database.Init()
	if err != nil {
		return err
	}

	stmt, err := database.DB.Prepare(wholeStmt)
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
	if dtx.BrandString != "" {
		go redis.Setex(redis_key, l.Years, 86400)
	}
	return nil
}
