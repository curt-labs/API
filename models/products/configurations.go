package products

import (
	"database/sql"
	"github.com/curt-labs/GoAPI/helpers/database"
	_ "github.com/go-sql-driver/mysql"
	"strconv"
	"strings"
)

var (
	getAllVehicleConfigOptionsStmt = `
		select distinct cat.name, cat.AcesTypeID from vcdb_Vehicle as v
		join VehicleConfigAttribute as vca on v.ConfigID = vca.VehicleConfigID
		join ConfigAttribute as ca on vca.AttributeID = ca.ID
		join ConfigAttributeType as cat on ca.ConfigAttributeTypeID = cat.ID
		join Submodel as s on v.SubModelID = s.ID
		join BaseVehicle as bv on v.BaseVehicleID = bv.ID
		join vcdb_Model as mo on bv.ModelID = mo.ID
		join vcdb_Make as ma on bv.MakeID = ma.ID
		join vcdb_VehiclePart as vp on v.ID = vp.VehicleID
		join Part as p on vp.PartNumber = p.partID
		where (p.status = 800 || p.status = 900) &&
		bv.YearID = ? && ma.MakeName = ? &&
		mo.ModelName = ? && s.SubmodelName = ?
		order by cat.sort`
	getDefinedConfigurationsForVehicleStmt = `
		select distinct cat.name, ca.value,
		(
			select GROUP_CONCAT(vp1.PartNumber order by vp1.PartNumber)
			from vcdb_VehiclePart as vp1
			join Part as p1 on vp1.PartNumber = p1.partID
			where vp1.VehicleID = v.ID && (p1.status = 800 || p1.status = 900)
		) as parts, vca.VehicleConfigID
		from vcdb_Vehicle as v
		join VehicleConfigAttribute as vca on v.ConfigID = vca.VehicleConfigID
		join ConfigAttribute as ca on vca.AttributeID = ca.ID
		join ConfigAttributeType as cat on ca.ConfigAttributeTypeID = cat.ID
		join Submodel as s on v.SubModelID = s.ID
		join BaseVehicle as bv on v.BaseVehicleID = bv.ID
		join vcdb_Model as mo on bv.ModelID = mo.ID
		join vcdb_Make as ma on bv.MakeID = ma.ID
		join vcdb_VehiclePart as vp on v.ID = vp.VehicleID
		join Part as p on vp.PartNumber = p.partID
		join ApiKeyToBrand as atb on atb.brandID = p.brandID
		join ApiKey as ak on ak.id = atb.keyID
		where (p.status = 800 || p.status = 900) &&
		ak.api_key = ? &&
		bv.YearID = ? && ma.MakeName = ? &&
		mo.ModelName = ? && s.SubmodelName = ?
		order by vca.VehicleConfigID`
	getAllOptionsForType = `
		select distinct ca.value from vcdb_Vehicle as v
		join VehicleConfigAttribute as vca on v.ConfigID = vca.VehicleConfigID
		join ConfigAttribute as ca on vca.AttributeID = ca.ID
		join ConfigAttributeType as cat on ca.ConfigAttributeTypeID = cat.ID
		join Submodel as s on v.SubModelID = s.ID
		join BaseVehicle as bv on v.BaseVehicleID = bv.ID
		join vcdb_Model as mo on bv.ModelID = mo.ID
		join vcdb_Make as ma on bv.MakeID = ma.ID
		join vcdb_VehiclePart as vp on v.ID = vp.VehicleID
		join Part as p on vp.PartNumber = p.partID
		where (p.status = 800 || p.status = 900) &&
		bv.YearID = ? && ma.MakeName = ? &&
		mo.ModelName = ? && s.SubmodelName = ? && cat.name = ?
		order by ca.value`
	vcdb_GetAspirationForVehicle = `
		select distinct a.AspirationID as ID, a.AspirationName as value from VehicleToEngineConfig as vec
		join EngineConfig as ec on vec.EngineConfigID = ec.EngineConfigID
		join Aspiration as a on ec.AspirationID = a.AspirationID
		join Vehicle as v on vec.VehicleID = v.VehicleID
		where v.VehicleID = ?
		order by value`
	vcdb_GetBedLengthForVehicle = `
		select distinct bl.BedLengthID as ID, bl.BedLength as value from VehicleToBedConfig as vbc
		join BedConfig as bc on vbc.BedConfigID = bc.BedConfigID
		join BedLength as bl on bc.BedLengthID = bl.BedLengthID
		join Vehicle as v on vbc.VehicleID = v.VehicleID
		where v.VehicleID = ?
		order by value`
	vcdb_GetBedTypeForVehicle = `
		select distinct bt.BedTypeID as ID, bt.BedTypeName as value from VehicleToBedConfig as vbc
		join BedConfig as bc on vbc.BedConfigID = bc.BedConfigID
		join BedType as bt on bc.BedTypeID = bt.BedTypeID
		join Vehicle as v on vbc.VehicleID = v.VehicleID
		where v.VehicleID = ?
		order by value`
	vcdb_GetBodyTypeForVehicle = `
		select distinct bt.BodyTypeID as ID, bt.BodyTypename as value from VehicleToBodyStyleConfig as vbsc
		join BodyStyleConfig as bsc on vbsc.BodyStyleConfigID = bsc.BodyStyleConfigID
		join BodyType as bt on bsc.BodyTypeID = bt.BodyTypeID
		join Vehicle as v on vbsc.VehicleID = v.VehicleID
		where v.VehicleID = ?
		order by value`
	vcdb_GetBodyNumDoorsForVehicle = `
		select distinct bnd.BodyNumDoorsID as ID, bnd.BodyNumDoors as value from VehicleToBodyStyleConfig as vbsc
		join BodyStyleConfig as bsc on vbsc.BodyStyleConfigID = bsc.BodyStyleConfigID
		join BodyNumDoors as bnd on bsc.BodyNumDoorsID = bnd.BodyNumDoorsID
		join Vehicle as v on vbsc.VehicleID = v.VehicleID
		where v.VehicleID = ?
		order by value`
	vcdb_GetBrakeABSForVehicle = `
		select distinct ba.BrakeABSID as ID, ba.BrakeABSName as value from VehicleToBrakeConfig as vbc
		join BrakeConfig as bc on vbc.BrakeConfigID = bc.BrakeConfigID
		join BrakeABS as ba on bc.BrakeABSID = ba.BrakeABSID
		join Vehicle as v on vbc.VehicleID = v.VehicleID
		where v.VehicleID = ?
		order by value`
	vcdb_GetBrakeSystemForVehicle = `
		select distinct bs.BrakeSystemID as ID, bs.BrakeSystemName as value from VehicleToBrakeConfig as vbc
		join BrakeConfig as bc on vbc.BrakeConfigID = bc.BrakeConfigID
		join BrakeSystem as bs on bc.BrakeSystemID = bs.BrakeSystemID
		join Vehicle as v on vbc.VehicleID = v.VehicleID
		where v.VehicleID = ?
		order by value`
	vcdb_GetFrontBrakeTypeForVehicle = `
		select distinct bt.BrakeTypeID as ID, bt.BrakeTypeName as value from VehicleToBrakeConfig as vbc
		join BrakeConfig as bc on vbc.BrakeConfigID = bc.BrakeConfigID
		join BrakeType as bt on bc.FrontBrakeTypeID = bt.BrakeTypeID
		join Vehicle as v on vbc.VehicleID = v.VehicleID
		where v.VehicleID = ?
		order by value`
	vcdb_GetRearBrakeTypeForVehicle = `
		select distinct bt.BrakeTypeID as ID, bt.BrakeTypeName as value from VehicleToBrakeConfig as vbc
		join BrakeConfig as bc on vbc.BrakeConfigID = bc.BrakeConfigID
		join BrakeType as bt on bc.RearBrakeTypeID = bt.BrakeTypeID
		join Vehicle as v on vbc.VehicleID = v.VehicleID
		where v.VehicleID = ?
		order by value`
	vcdb_GetCylinderHeadTypeForVehicle = `
		select distinct cht.CylinderHeadTypeID as ID, cht.CylinderHeadTypeName as value from VehicleToEngineConfig as vec
		join EngineConfig as ec on vec.EngineConfigID = ec.EngineConfigID
		join CylinderHeadType as cht on ec.CylinderHeadTypeID = cht.CylinderHeadTypeID
		join Vehicle as v on vec.VehicleID = v.VehicleID
		where v.VehicleID = ?
		order by value`
	vcdb_GetDriveTypeForVehicle = `
		select distinct dt.DriveTypeID as ID, dt.DriveTypeName as value from VehicleToDriveType as vdt
		join DriveType as dt on vdt.DriveTypeID = dt.DriveTypeID
		join Vehicle as v on vdt.VehicleID = v.VehicleID
		where v.VehicleID = ?
		order by value`
	vcdb_GetEngineDesignationForVehicle = `
		select distinct ed.EngineDesignationID as ID, ed.EngineDesignationName as value from VehicleToEngineConfig as vec
		join EngineConfig as ec on vec.EngineConfigID = ec.EngineConfigID
		join EngineDesignation as ed on ec.EngineDesignationID = ed.EngineDesignationID
		join Vehicle as v on vec.VehicleID = v.VehicleID
		where v.VehicleID = ?
		order by value`
	vcdb_GetEngineVersionForVehicle = `
		select distinct ev.EngineVersionID as ID, ev.EngineVersion as value from VehicleToEngineConfig as vec
		join EngineConfig as ec on vec.EngineConfigID = ec.EngineConfigID
		join EngineVersion as ev on ec.EngineVersionID = ev.EngineVersionID
		join Vehicle as v on vec.VehicleID = v.VehicleID
		where v.VehicleID = ?
		order by value`
	vcdb_GetEngineVINForVehicle = `
		select distinct ev.EngineVINID as ID, ev.EngineVINName as value from VehicleToEngineConfig as vec
		join EngineConfig as ec on vec.EngineConfigID = ec.EngineConfigID
		join EngineVIN as ev on ec.EngineVINID = ev.EngineVINID
		join Vehicle as v on vec.VehicleID = v.VehicleID
		where v.VehicleID = ?
		order by value`
	vcdb_GetFuelDeliverySubTypeForVehicle = `
		select distinct fdst.FuelDeliverySubTypeID as ID, fdst.FuelDeliverySubTypeName as value from VehicleToEngineConfig as vec
		join EngineConfig as ec on vec.EngineConfigID = ec.EngineConfigID
		join FuelDeliveryConfig as fdc on ec.FuelDeliveryConfigID
		join FuelDeliverySubType as fdst on fdc.FuelDeliverySubTypeID = fdst.FuelDeliverySubTypeID
		join Vehicle as v on vec.VehicleID = v.VehicleID
		where v.VehicleID = ?
		order by value`
	vcdb_GetFuelDeliveryTypeForVehicle = `
		select distinct fdt.FuelDeliveryTypeID as ID, fdt.FuelDeliveryTypeName as value from VehicleToEngineConfig as vec
		join EngineConfig as ec on vec.EngineConfigID = ec.EngineConfigID
		join FuelDeliveryConfig as fdc on ec.FuelDeliveryConfigID
		join FuelDeliveryType as fdt on fdc.FuelDeliveryTypeID = fdt.FuelDeliveryTypeID
		join Vehicle as v on vec.VehicleID = v.VehicleID
		where v.VehicleID = ?
		order by value`
	vcdb_GetFuelSystemControlTypeForVehicle = `
		select distinct fsct.FuelSystemControlTypeID as ID, fsct.FuelSystemControlTypeName as value from VehicleToEngineConfig as vec
		join EngineConfig as ec on vec.EngineConfigID = ec.EngineConfigID
		join FuelDeliveryConfig as fdc on ec.FuelDeliveryConfigID
		join FuelSystemControlType as fsct on fdc.FuelSystemControlTypeID = fsct.FuelSystemControlTypeID
		join Vehicle as v on vec.VehicleID = v.VehicleID
		where v.VehicleID = ?
		order by value`
	vcdb_GetFuelSystemDesignForVehicle = `
		select distinct fsd.FuelSystemDesignID as ID, fsd.FuelSystemDesignName as value from VehicleToEngineConfig as vec
		join EngineConfig as ec on vec.EngineConfigID = ec.EngineConfigID
		join FuelDeliveryConfig as fdc on ec.FuelDeliveryConfigID
		join FuelSystemDesign as fsd on fdc.FuelSystemDesignID = fsd.FuelSystemDesignID
		join Vehicle as v on vec.VehicleID = v.VehicleID
		where v.VehicleID = ?
		order by value`
	vcdb_GetFuelTypeForVehicle = `
		select distinct ft.FuelTypeID as ID, ft.FuelTypeName as value from VehicleToEngineConfig as vec
		join EngineConfig as ec on vec.EngineConfigID = ec.EngineConfigID
		join FuelType as ft on ec.FuelTypeID = ft.FuelTypeID
		join Vehicle as v on vec.VehicleID = v.VehicleID
		where v.VehicleID = ?
		order by value`
	vcdb_GetIgnitionSystemForVehicle = `
		select distinct ist.IgnitionSystemTypeID as ID, ist.IgnitionSystemTypeName as value from VehicleToEngineConfig as vec
		join EngineConfig as ec on vec.EngineConfigID = ec.EngineConfigID
		join IgnitionSystemType as ist on ec.IgnitionSystemTypeID = ist.IgnitionSystemTypeID
		join Vehicle as v on vec.VehicleID = v.VehicleID
		where v.VehicleID = ?
		order by value`
	vcdb_GetMfrBodyCodeForVehicle = `
		select distinct mbc.MfrBodyCodeID as ID, mbc.MfrBodyCodeName as value from VehicleToMfrBodyCode as vmbc
		join MfrBodyCode as mbc on vmbc.MfrBodyCodeID = mbc.MfrBodyCodeID
		join Vehicle as v on vmbc.VehicleID = v.VehicleID
		where v.VehicleID = ?
		order by value`
	vcdb_GetFrontSpringTypeForVehicle = `
		select distinct st.SpringTypeID as ID, st.SpringTypeName as value from VehicleToSpringTypeConfig as vstc
		join SpringTypeConfig as stc on vstc.SpringTypeConfigID = stc.SpringTypeConfigID
		join SpringType as st on stc.FrontSpringTypeID = st.SpringTypeID
		join Vehicle as v on vstc.VehicleID = v.VehicleID
		where v.VehicleID = ?
		order by value`
	vcdb_GetRearSpringTypeForVehicle = `
		select distinct st.SpringTypeID as ID, st.SpringTypeName as value from VehicleToSpringTypeConfig as vstc
		join SpringTypeConfig as stc on vstc.SpringTypeConfigID = stc.SpringTypeConfigID
		join SpringType as st on stc.RearSpringTypeID = st.SpringTypeID
		join Vehicle as v on vstc.VehicleID = v.VehicleID
		where v.VehicleID = ?
		order by value`
	vcdb_GetSteeringSystemForVehicle = `
		select distinct ss.SteeringSystemID as ID, ss.SteeringSystemName as value from VehicleToSteeringConfig as vsc
		join SteeringConfig as sc on vsc.SteeringConfigID = sc.SteeringConfigID
		join SteeringSystem as ss on sc.SteeringSystemID = ss.SteeringSystemID
		join Vehicle as v on vsc.VehicleID = v.VehicleID
		where v.VehicleID = ?
		order by value`
	vcdb_GetSteeringTypeForVehicle = `
		select distinct st.SteeringTypeID as ID, st.SteeringTypeName as value from VehicleToSteeringConfig as vsc
		join SteeringConfig as sc on vsc.SteeringConfigID = sc.SteeringConfigID
		join SteeringType as st on sc.SteeringTypeID = st.SteeringTypeID
		join Vehicle as v on vsc.VehicleID = v.VehicleID
		where v.VehicleID = ?
		order by value`
	vcdb_GetTransmissionControlTypeForVehicle = `
		select distinct tct.TransmissionControlTypeID as ID, tct.TransmissionControlTypeName as value from VehicleToTransmission as vt
		join Transmission as t on vt.TransmissionID = t.TransmissionID
		join TransmissionBase as tb on t.TransmissionBaseID = tb.TransmissionBaseID
		join TransmissionControlType as tct on tb.TransmissionControlTypeID = tct.TransmissionControlTypeID
		join Vehicle as v on vt.VehicleID = v.VehicleID
		where v.VehicleID = ?
		order by value`
	vcdb_GetElecControlledForVehicle = `
		select distinct ec.ElecControlledID as ID, ec.ElecControlled as value from VehicleToTransmission as vt
		join Transmission as t on vt.TransmissionID = t.TransmissionID
		join ElecControlled as ec on t.TransmissionElecControlledID = ec.ElecControlledID
		join Vehicle as v on vt.VehicleID = v.VehicleID
		where v.VehicleID = ?
		order by value`
	vcdb_GetTransmissionMfrCodeForVehicle = `
		select distinct tmc.TransmissionMfrCodeID as ID, tmc.TransmissionMfrCode as value from VehicleToTransmission as vt
		join Transmission as t on vt.TransmissionID = t.TransmissionID
		join TransmissionMfrCode as tmc on t.TransmissionMfrCodeID = tmc.TransmissionMfrCodeID
		join Vehicle as v on vt.VehicleID = v.VehicleID
		where v.VehicleID = ?
		order by value`
	vcdb_GetTransmissionNumSpeedsForVehicle = `
		select distinct tns.TransmissionNumSpeedsID as ID, tns.TransmissionNumSpeeds as value from VehicleToTransmission as vt
		join Transmission as t on vt.TransmissionID = t.TransmissionID
		join TransmissionBase as tb on t.TransmissionBaseID = tb.TransmissionBaseID
		join TransmissionNumSpeeds as tns on tb.TransmissionNumSpeedsID = tns.TransmissionNumSpeedsID
		join Vehicle as v on vt.VehicleID = v.VehicleID
		where v.VehicleID = ?
		order by value`
	vcdb_GetTransmissionTypeForVehicle = `
		select distinct tt.TransmissionTypeID as ID, tt.TransmissionTypeName as value from VehicleToTransmission as vt
		join Transmission as t on vt.TransmissionID = t.TransmissionID
		join TransmissionBase as tb on t.TransmissionBaseID = tb.TransmissionBaseID
		join TransmissionType as tt on tb.TransmissionTypeID = tt.TransmissionTypeID
		join Vehicle as v on vt.VehicleID = v.VehicleID
		where v.VehicleID = ?
		order by value`
	vcdb_GetValvesForVehicle = `
		select distinct va.ValvesID as ID, va.ValvesPerEngine as value from VehicleToEngineConfig as vec
		join EngineConfig as ec on vec.EngineConfigID = ec.EngineConfigID
		join Valves as va on ec.ValvesID = va.ValvesID
		join Vehicle as v on vec.VehicleID = v.VehicleID
		where v.VehicleID = ?
		order by value`
	vcdb_GetWheelBaseForVehicle = `
		select distinct wb.WheelbaseID as ID, wb.WheelBase as value from VehicleToWheelbase as vwb
		join WheelBase as wb on vwb.WheelbaseID = wb.WheelBaseID
		join Vehicle as v on vwb.VehicleID = v.VehicleID
		where v.VehicleID = ?
		order by value`
)

type Configuration struct {
	Key   string `json:"key" xml:"key"`
	Value string `json:"value" xml:"value"`
}

type ConfigurationOption struct {
	Type    string   `json:"type" xml:"type"`
	Options []string `json:"options" xml:"options"`
}

type DefinedConfiguration struct {
	Type     string
	Value    string
	Parts    []int
	ConfigID int
}

func (l *Lookup) GetConfigurations() error {
	db, err := sql.Open("mysql", database.ConnectionString())
	if err != nil {
		return err
	}
	defer db.Close()

	stmt, err := db.Prepare(getAllVehicleConfigOptionsStmt)
	if err != nil {
		return err
	}
	defer stmt.Close()

	res, err := stmt.Query(l.Vehicle.Base.Year, l.Vehicle.Base.Make, l.Vehicle.Base.Model, l.Vehicle.Submodel)
	if err != nil {
		return err
	}

	l.Configurations = make([]ConfigurationOption, 0)
	count := 0
	ch := make(chan error)
	for res.Next() {
		var conf Configuration
		var acesType int
		err = res.Scan(&conf.Key, &acesType)
		if err == nil {
			go conf.allOptions(l, acesType, ch)
			count++
		}
	}
	defer res.Close()

	for i := 0; i < count; i++ {
		<-ch
	}

	l.Pagination = Pagination{
		TotalItems:    len(l.Submodels),
		ReturnedCount: len(l.Submodels),
		Page:          1,
		PerPage:       len(l.Submodels),
		TotalPages:    1,
	}

	return nil
}

func (v Vehicle) getDefinedConfigurations(apiKey string) (*map[int][]DefinedConfiguration, error) {
	configs := make(map[int][]DefinedConfiguration, 0)

	db, err := sql.Open("mysql", database.ConnectionString())
	if err != nil {
		return nil, err
	}
	defer db.Close()

	stmt, err := db.Prepare(getDefinedConfigurationsForVehicleStmt)
	if err != nil {
		return nil, err
	}
	defer stmt.Close()

	rows, err := stmt.Query(apiKey, v.Base.Year, v.Base.Make, v.Base.Model, v.Submodel)
	if err != nil || rows == nil {
		return nil, err
	}

	for rows.Next() {
		var name string
		var value string
		var parts string
		var configID int
		err = rows.Scan(&name, &value, &parts, &configID)
		if err == nil {
			dc := DefinedConfiguration{
				Type:     name,
				Value:    value,
				ConfigID: configID,
				Parts:    make([]int, 0),
			}
			partArr := strings.Split(parts, ",")
			for _, p := range partArr {
				if partID, err := strconv.Atoi(p); err == nil {
					dc.Parts = append(dc.Parts, partID)
				}
			}
			if _, ok := configs[dc.ConfigID]; !ok {
				configs[dc.ConfigID] = make([]DefinedConfiguration, 0)
			}
			configs[dc.ConfigID] = append(configs[dc.ConfigID], dc)
		}
	}
	defer rows.Close()

	return &configs, nil
}

// allOptions will populate the available configuration options
// from preferrable the VCDB, or CurtDev if it's a custom configuration
// type that match the current vehicle.
// This is meant to be run in a goroutine, hence the channel.
func (c Configuration) allOptions(l *Lookup, acesType int, ch chan error) {
	var opts []string
	var err error

	if acesType == 0 { // Custom configuraton (implicit)
		opts, err = c.getCurtOptions(l.Vehicle)
		if err != nil {
			ch <- err
			return
		}
	} else { // VCDB configuration (explicit)
		opts, err = c.getVcdbOptions(l.Vehicle)
		if err != nil {
			ch <- err
			return
		}
	}

	co := ConfigurationOption{
		Type:    c.Key,
		Options: opts,
	}

	// Check if this configuration option has already been selected
	exists := false
	for _, opt := range l.Vehicle.Configurations {
		if strings.ToLower(opt.Key) == strings.ToLower(co.Type) {
			exists = true
			break
		}
	}

	if !exists {
		l.Configurations = append(l.Configurations, co)
	}

	ch <- err
}

// getVcdbOptions will return the configuration options
// that fit the provided vehicel and the provided type from the VCDB.
func (c *Configuration) getVcdbOptions(v Vehicle) ([]string, error) {
	var opts []string
	var err error

	id, err := v.GetVcdbID()
	if err != nil || id == 0 {
		return opts, err
	}

	db, err := sql.Open("mysql", database.VcdbConnectionString())
	if err != nil {
		return opts, err
	}
	defer db.Close()

	var stmt *sql.Stmt
	switch strings.ToLower(strings.Replace(c.Key, " ", "", -1)) {
	case "aspiration":
		stmt, err = db.Prepare(vcdb_GetAspirationForVehicle)
	case "bedlength":
		stmt, err = db.Prepare(vcdb_GetBedLengthForVehicle)
	case "bedtype":
		stmt, err = db.Prepare(vcdb_GetBedTypeForVehicle)
	case "bodytype":
		stmt, err = db.Prepare(vcdb_GetBodyTypeForVehicle)
	case "brakeabs":
		stmt, err = db.Prepare(vcdb_GetBrakeABSForVehicle)
	case "brakesystem":
		stmt, err = db.Prepare(vcdb_GetBrakeSystemForVehicle)
	case "frontbraketype":
		stmt, err = db.Prepare(vcdb_GetFrontBrakeTypeForVehicle)
	case "rearbraketype":
		stmt, err = db.Prepare(vcdb_GetRearBrakeTypeForVehicle)
	case "cylinderheadtype":
		stmt, err = db.Prepare(vcdb_GetCylinderHeadTypeForVehicle)
	case "drivetype":
		stmt, err = db.Prepare(vcdb_GetDriveTypeForVehicle)
	case "enginedesignation":
		stmt, err = db.Prepare(vcdb_GetEngineDesignationForVehicle)
	case "engineversion":
		stmt, err = db.Prepare(vcdb_GetEngineVersionForVehicle)
	case "enginevin":
		stmt, err = db.Prepare(vcdb_GetEngineVINForVehicle)
	case "fueldeliverysubtype":
		stmt, err = db.Prepare(vcdb_GetFuelDeliverySubTypeForVehicle)
	case "fueldeliverytype":
		stmt, err = db.Prepare(vcdb_GetFuelDeliveryTypeForVehicle)
	case "fuelsystemcontroltype":
		stmt, err = db.Prepare(vcdb_GetFuelSystemControlTypeForVehicle)
	case "fuelsystemdesign":
		stmt, err = db.Prepare(vcdb_GetFuelSystemDesignForVehicle)
	case "fueltype":
		stmt, err = db.Prepare(vcdb_GetFuelTypeForVehicle)
	case "ignitionsystem":
		stmt, err = db.Prepare(vcdb_GetIgnitionSystemForVehicle)
	case "mfrbodycode":
		stmt, err = db.Prepare(vcdb_GetMfrBodyCodeForVehicle)
	case "numberofdoors":
		stmt, err = db.Prepare(vcdb_GetBodyNumDoorsForVehicle)
	case "frontspringtype":
		stmt, err = db.Prepare(vcdb_GetFrontSpringTypeForVehicle)
	case "rearspringtype":
		stmt, err = db.Prepare(vcdb_GetRearSpringTypeForVehicle)
	case "steeringsystem":
		stmt, err = db.Prepare(vcdb_GetSteeringSystemForVehicle)
	case "steeringtype":
		stmt, err = db.Prepare(vcdb_GetSteeringTypeForVehicle)
	case "transmissioncontroltype":
		stmt, err = db.Prepare(vcdb_GetTransmissionControlTypeForVehicle)
	case "transmissionelectroniccontrolled":
		stmt, err = db.Prepare(vcdb_GetElecControlledForVehicle)
	case "transmissionmanufacturercode":
		stmt, err = db.Prepare(vcdb_GetTransmissionMfrCodeForVehicle)
	case "transmissionnumspeeds":
		stmt, err = db.Prepare(vcdb_GetTransmissionNumSpeedsForVehicle)
	case "transmissiontype":
		stmt, err = db.Prepare(vcdb_GetTransmissionTypeForVehicle)
	case "valves":
		stmt, err = db.Prepare(vcdb_GetValvesForVehicle)
	case "wheelbase":
		stmt, err = db.Prepare(vcdb_GetWheelBaseForVehicle)
	default:
	}
	if err != nil || stmt == nil {
		return opts, err
	}
	defer stmt.Close()

	res, err := stmt.Query(id)
	if err != nil || res == nil {
		return opts, err
	}

	for res.Next() {
		var id int
		var val string
		if err = res.Scan(&id, &val); err == nil {
			opts = append(opts, val)
		}
	}
	defer res.Close()

	return opts, nil
}

// getCurtOptions will return the configuration options
// that fit the provided vehicle and the provided type from CurtDev.
func (c *Configuration) getCurtOptions(v Vehicle) ([]string, error) {
	var opts []string
	var err error

	db, err := sql.Open("mysql", database.ConnectionString())
	if err != nil {
		return opts, err
	}
	defer db.Close()

	stmt, err := db.Prepare(getAllOptionsForType)
	if err != nil {
		return opts, err
	}
	defer stmt.Close()

	res, err := stmt.Query(v.Base.Year, v.Base.Make, v.Base.Model, v.Submodel, c.Key)
	if err != nil || res == nil {
		return opts, err
	}

	hasOther := false
	for res.Next() {
		var val string
		if err = res.Scan(&val); err == nil {
			if strings.ToLower(val) == "other" {
				hasOther = true
			}
			opts = append(opts, val)
		}
	}
	defer res.Close()

	if !hasOther && len(opts) < 2 {
		opts = append(opts, "Other")
	}

	return opts, nil
}
