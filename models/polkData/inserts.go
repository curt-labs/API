package polkData

import (
	"database/sql"
	"github.com/curt-labs/GoAPI/helpers/database"
	_ "github.com/go-sql-driver/mysql"
)

var (
	insertPartIntoVehiclePart         = `insert into vcdb_VehiclePart (VehicleID, PartNumber) values(?,?)`
	insertBaseVehicleIntoVcdbVehicles = `insert into vcdb_Vehicle (BaseVehicleID, SubmodelID, ConfigID, AppID, RegionID) values (?,0,0,0,0)`
	insertSubmodelIntoVcdbVehicles    = `insert into vcdb_Vehicle (BaseVehicleID, SubmodelID, ConfigID, AppID, RegionID) values (?,?,0,0,0)`
)

func (c *CsvDatum) InsertPartIntoVehiclePart() error {
	var err error
	db, err := sql.Open("mysql", database.ConnectionString())
	if err != nil {
		return err
	}
	defer db.Close()

	stmt, err := db.Prepare(insertPartIntoVehiclePart)
	if err != nil {
		return err
	}
	defer stmt.Close()
	_, err = stmt.Exec(c.CurtVehicle.ID, c.Part.ID)
	return err
}

func (c *CsvDatum) InsertBaseVehicleIntoVcdbVehicles() error {
	db, err := sql.Open("mysql", database.ConnectionString())
	if err != nil {
		return err
	}
	defer db.Close()

	stmt, err := db.Prepare(insertBaseVehicleIntoVcdbVehicles)
	if err != nil {
		return err
	}
	defer stmt.Close()
	res, err := stmt.Exec(c.CurtVehicle.CurtBaseID)
	if err != nil {
		return err
	}
	id, err := res.LastInsertId()
	c.CurtVehicle.ID = int(id)
	return err
}

func (c *CsvDatum) InsertSubmodelIntoVcdbVehicles() error {
	db, err := sql.Open("mysql", database.ConnectionString())
	if err != nil {
		return err
	}
	defer db.Close()

	stmt, err := db.Prepare(insertSubmodelIntoVcdbVehicles)
	if err != nil {
		return err
	}
	defer stmt.Close()
	res, err := stmt.Exec(c.CurtVehicle.CurtBaseID, c.CurtVehicle.CurtSubmodelID)
	if err != nil {
		return err
	}
	id, err := res.LastInsertId()
	c.CurtVehicle.ID = int(id)
	return err
}
