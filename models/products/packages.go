package products

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"github.com/curt-labs/GoAPI/helpers/apicontext"
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
	createPackage  = `INSERT INTO PartPackage (partID, height, width, length, weight, dimensionUOM, weightUOM, packageUOM, quantity, typeID)`
	deletePackages = `DELETE FROM PartPackages WHERE partID = ?`
)

type Package struct {
	ID                 int         `json:"id,omitempty" xml:"id,omitempty"`
	PartID             int         `json:"partId,omitempty" xml:"partId,omitempty"`
	Height             float64     `json:"height,omitempty" xml:"height,omitempty"`
	Width              float64     `json:"width,omitempty" xml:"width,omitempty"`
	Length             float64     `json:"length,omitempty" xml:"length,omitempty"`
	Weight             float64     `json:"weight,omitempty" xml:"weight,omitempty"`
	DimensionUnit      string      `json:"dimensionUnit,omitempty" xml:"dimensionUnit,omitempty"`
	DimensionUnitLabel string      `json:"dimensionUnitLabel,omitempty" xml:"dimensionUnitLabel,omitempty"`
	WeightUnit         string      `json:"weightUnit,omitempty" xml:"weightUnit,omitempty"`
	WeightUnitLabel    string      `json:"weightUnitLabel,omitempty" xml:"weightUnitLabel,omitempty"`
	PackageUnit        string      `json:"packageUnit,omitempty" xml:"packageUnit,omitempty"`
	PackageUnitLabel   string      `json:"packageUnitLabel,omitempty" xml:"packageUnitLabel,omitempty"`
	Quantity           int         `json:"quantity,omitempty" xml:"quantity,omitempty"`
	PackageType        PackageType `json:"packageType,omitempty" xml:"packageType,omitempty"`
}

type PackageType struct {
	ID   int    `json:"id,omitempty" xml:"id,omitempty"`
	Name string `json:"name,omitempty" xml:"name,omitempty"`
}

func (p *Part) GetPartPackaging(dtx *apicontext.DataContext) error {
	redis_key := fmt.Sprintf("part:%d:packages:%s", p.ID, dtx.BrandString)

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

	rows, err := qry.Query(p.ID)
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
	defer rows.Close()

	p.Packages = pkgs
	if dtx.BrandString != "" {
		go redis.Setex(redis_key, p.Packages, redis.CacheTimeout)
	}
	return nil
}

func (p *Package) Create() (err error) {
	db, err := sql.Open("mysql", database.ConnectionString())
	if err != nil {
		return err
	}
	defer db.Close()
	stmt, err := db.Prepare(createPackage)
	if err != nil {
		return err
	}
	defer stmt.Close()
	res, err := stmt.Exec(p.PartID, p.Height, p.Width, p.Length, p.DimensionUnit, p.WeightUnit, p.PackageUnit, p.PackageType.ID)
	if err != nil {
		return err
	}
	id, err := res.LastInsertId()
	if err != nil {
		return err
	}
	p.ID = int(id)
	return
}

func (p *Package) DeleteByPart() (err error) {
	db, err := sql.Open("mysql", database.ConnectionString())
	if err != nil {
		return err
	}
	defer db.Close()
	stmt, err := db.Prepare(deletePackages)
	if err != nil {
		return err
	}
	defer stmt.Close()
	_, err = stmt.Exec(p.PartID)
	if err != nil {
		return err
	}
	return nil
}
