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
	insertBaseVehicle                 = `insert into BaseVehicle (AAIABaseVehicleID, YearID, MakeID, ModelID) 
		values(?,(select YearID from vcdb_Year where YearID = ?),(select AAIAMakeID from vcdb_Make where ID = ?), (select AAIAModelID from vcdb_Model where ID = ?))`
	insertSubmodel = `insert into Submodel (AAIASubmodelID, SubmodelName) values(?,?)`
)

func (c *CsvDatum) InsertPartIntoVehiclePart() error {
	var err error
	db, err := sql.Open("mysql", database.AriesConnectionString())
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
	db, err := sql.Open("mysql", database.AriesConnectionString())
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
	db, err := sql.Open("mysql", database.AriesConnectionString())
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

func (c *CsvDatum) InsertBaseVehicle() error {
	var err error
	db, err := sql.Open("mysql", database.AriesConnectionString())
	if err != nil {
		return err
	}
	defer db.Close()

	stmt, err := db.Prepare(insertBaseVehicle)
	if err != nil {
		return err
	}
	defer stmt.Close()
	res, err := stmt.Exec(c.CsvVehicle.BaseVehicleID, c.CsvVehicle.YearID, c.CsvVehicle.MakeID, c.CsvVehicle.ModelID)
	if err != nil {
		return err
	}
	id, err := res.LastInsertId()
	if err != nil {
		return err
	}
	c.CurtVehicle.CurtBaseID = int(id)
	return err
}

func (c *CsvDatum) InsertSubmodel() error {
	var err error
	db, err := sql.Open("mysql", database.AriesConnectionString())
	if err != nil {
		return err
	}
	defer db.Close()

	stmt, err := db.Prepare(insertSubmodel)
	if err != nil {
		return err
	}
	defer stmt.Close()
	res, err := stmt.Exec(c.CsvVehicle.SubmodelID, c.CsvVehicle.SubModel)
	if err != nil {
		return err
	}
	id, err := res.LastInsertId()
	if err != nil {
		return err
	}
	c.CurtVehicle.CurtSubmodelID = int(id)
	return err
}
