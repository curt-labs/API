package products

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"github.com/curt-labs/GoAPI/helpers/database"
	"github.com/curt-labs/GoAPI/helpers/redis"
	_ "github.com/go-sql-driver/mysql"
)

var (
	partPackageStmt = `select pp.height as height, pp.length as length, pp.width as width, pp.weight as weight, pp.quantity as quantity,
				um_dim.code as dimensionUnit, um_dim.name as dimensionUnitLabel, um_wt.code as weightUnit, um_wt.name as weightUnitLabel,
				um_pkg.code as packageUnit, um_pkg.name as packageUnitLabel
				from PartPackage as pp
				join UnitOfMeasure as um_dim on pp.dimensionUOM = um_dim.ID
				join UnitOfMeasure as um_wt on pp.weightUOM = um_wt.ID
				join UnitOfMeasure as um_pkg on pp.packageUOM = um_pkg.ID
				where pp.partID = ?`
)

type Package struct {
	Height, Width, Length, Quantity   float64
	Weight                            float64
	DimensionUnit, DimensionUnitLabel string
	WeightUnit, WeightUnitLabel       string
	PackageUnit, PackageUnitLabel     string
}

func (p *Part) GetPartPackaging() error {
	redis_key := fmt.Sprintf("part:%d:packages", p.PartId)

	data, err := redis.Get(redis_key)
	if err == nil && len(data) > 0 {
		if err = json.Unmarshal(data, &p.Packages); err != nil {
			return nil
		}
	}

	db, err := sql.Open("mysql", database.ConnectionString())
	if err != nil {
		return err
	}
	defer db.Close()

	qry, err := db.Prepare(partPackageStmt)
	if err != nil {
		return err
	}
	defer qry.Close()

	rows, err := qry.Query(p.PartId)
	if err != nil {
		return err
	}

	var pkgs []Package
	for rows.Next() {
		var pkg Package
		err = rows.Scan(
			&pkg.Height,
			&pkg.Length,
			&pkg.Width,
			&pkg.Weight,
			&pkg.Quantity,
			&pkg.DimensionUnit,
			&pkg.DimensionUnitLabel,
			&pkg.WeightUnit,
			&pkg.WeightUnitLabel,
			&pkg.PackageUnit,
			&pkg.PackageUnitLabel)
		if err == nil {
			pkgs = append(pkgs, pkg)
		}
	}

	p.Packages = pkgs

	go redis.Setex(redis_key, p.Packages, redis.CacheTimeout)

	return nil
}
