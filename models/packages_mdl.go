package models

import (
	"github.com/curt-labs/GoAPI/helpers/database"
	"strconv"
	"strings"
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

	partPackageStmt_Grouped = `select pp.partID, pp.height as height, pp.length as length, pp.width as width, pp.weight as weight, pp.quantity as quantity,
				um_dim.code as dimensionUnit, um_dim.name as dimensionUnitLabel, um_wt.code as weightUnit, um_wt.name as weightUnitLabel,
				um_pkg.code as packageUnit, um_pkg.name as packageUnitLabel
				from PartPackage as pp
				join UnitOfMeasure as um_dim on pp.dimensionUOM = um_dim.ID
				join UnitOfMeasure as um_wt on pp.weightUOM = um_wt.ID
				join UnitOfMeasure as um_pkg on pp.packageUOM = um_pkg.ID
				where pp.partID IN (%s)`
)

type Package struct {
	Height, Width, Length, Quantity   float64
	Weight                            float64
	DimensionUnit, DimensionUnitLabel string
	WeightUnit, WeightUnitLabel       string
	PackageUnit, PackageUnitLabel     string
}

func (part *Part) GetPartPackaging() error {
	qry, err := database.Db.Prepare(partPackageStmt)
	if err != nil {
		return err
	}

	rows, res, err := qry.Exec(part.PartId)
	if database.MysqlError(err) {
		return err
	}

	height := res.Map("height")
	length := res.Map("length")
	width := res.Map("width")
	weight := res.Map("weight")
	qty := res.Map("quantity")
	dimUnit := res.Map("dimensionUnit")
	dimUnitLabel := res.Map("dimensionUnitLabel")
	weightUnit := res.Map("weightUnit")
	weightUnitLabel := res.Map("weightUnitLabel")
	pkgUnit := res.Map("packageUnit")
	pkgUnitLabel := res.Map("packageUnitLabel")

	var pkgs []Package
	for _, row := range rows {
		p := Package{
			Height:             row.Float(height),
			Width:              row.Float(width),
			Length:             row.Float(length),
			Quantity:           row.Float(qty),
			Weight:             row.Float(weight),
			DimensionUnit:      row.Str(dimUnit),
			DimensionUnitLabel: row.Str(dimUnitLabel),
			WeightUnit:         row.Str(weightUnit),
			WeightUnitLabel:    row.Str(weightUnitLabel),
			PackageUnit:        row.Str(pkgUnit),
			PackageUnitLabel:   row.Str(pkgUnitLabel),
		}
		pkgs = append(pkgs, p)
	}

	part.Packages = pkgs
	return nil
}

func (lookup *Lookup) GetPartPackaging() error {

	var ids []string
	for _, p := range lookup.Parts {
		ids = append(ids, strconv.Itoa(p.PartId))
	}
	if len(ids) == 0 {
		return nil
	}

	rows, res, err := database.Db.Query(partPackageStmt_Grouped, strings.Join(ids, ","))
	if database.MysqlError(err) || len(rows) == 0 {
		return err
	}

	partID := res.Map("partID")
	height := res.Map("height")
	length := res.Map("length")
	width := res.Map("width")
	weight := res.Map("weight")
	qty := res.Map("quantity")
	dimUnit := res.Map("dimensionUnit")
	dimUnitLabel := res.Map("dimensionUnitLabel")
	weightUnit := res.Map("weightUnit")
	weightUnitLabel := res.Map("weightUnitLabel")
	pkgUnit := res.Map("packageUnit")
	pkgUnitLabel := res.Map("packageUnitLabel")

	packages := make(map[int][]Package, len(lookup.Parts))

	for _, row := range rows {
		pId := row.Int(partID)

		p := Package{
			Height:             row.Float(height),
			Width:              row.Float(width),
			Length:             row.Float(length),
			Quantity:           row.Float(qty),
			Weight:             row.Float(weight),
			DimensionUnit:      row.Str(dimUnit),
			DimensionUnitLabel: row.Str(dimUnitLabel),
			WeightUnit:         row.Str(weightUnit),
			WeightUnitLabel:    row.Str(weightUnitLabel),
			PackageUnit:        row.Str(pkgUnit),
			PackageUnitLabel:   row.Str(pkgUnitLabel),
		}

		packages[pId] = append(packages[pId], p)
	}

	for _, p := range lookup.Parts {
		p.Packages = packages[p.PartId]
	}

	return nil
}
