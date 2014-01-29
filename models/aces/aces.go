package aces

import (
	"github.com/curt-labs/GoAPI/helpers/database"
	"github.com/curt-labs/GoAPI/helpers/redis"
	"github.com/curt-labs/GoAPI/helpers/xml"
	"strconv"
	"strings"
	"time"
)

type ACESBaseData struct {
	Id                int
	Qty               int
	PartTypeId        int
	ManufacturerLabel string
	PartNumber        int
	SubmodelID        int
	BaseVehicleID     int
	ClassId           int
	ConfigIds         []string
	ConfigNames       []string
	Notes             []string
	PartNotes         []string
}

type ACESVehicleData struct {
	BaseVehicleID int
	Submodel      string
	Config        []ACESConfigData
	Notes         []string
}

type ACESConfigData struct {
	Id    int
	Value string
}

func GetUniqueACESPartNumbers() (ids []int, err error) {

	qry, err := database.GetStatement("UniqueACESPartNumbers")
	if database.MysqlError(err) {
		return
	}

	rows, _, err := qry.Exec()
	if database.MysqlError(err) {
		return
	}

	for _, row := range rows {
		ids = append(ids, row.Int(0))
	}
	return
}

func GetACESPartData() (resp string, err error) {

	redis_key := "reports:aces:base:data"
	data_bytes, err := redis.RedisClient.Get(redis_key)
	// if err == nil {
	// 	//return string(data_bytes), nil
	// }

	qry, err := database.GetStatement("BigFatACESQuery")
	if database.MysqlError(err) {
		return
	}

	rows, res, err := qry.Exec()
	if database.MysqlError(err) {
		return
	}

	vID := res.Map("ID")
	typeID := res.Map("ACESPartTypeID")
	sDesc := res.Map("shortDesc")
	partID := res.Map("PartNumber")
	baseID := res.Map("AAIABaseVehicleID")
	subID := res.Map("AAIASubmodelID")
	classId := res.Map("classID")
	configIds := res.Map("configIDs")
	typeIds := res.Map("configNames")
	notes := res.Map("notes")
	part_notes := res.Map("part_notes")

	//types, err := GetAcesTypes()

	// Create Header
	doc := xml_helper.E("ACES",
		xml_helper.A("xmlns:xsi", "http://www.w3.org/2001/XMLSchema-instance"),
		xml_helper.A("xmlns:xsd", "http://www.w3.org/2001/XMLSchema"),
		xml_helper.A("version", "1.0"))
	header := xml_helper.E("Header",
		xml_helper.E("Company", xml_helper.T("CURT Manufacturing")),
		xml_helper.E("SenderName", xml_helper.T("Nichole Scott")),
		xml_helper.E("SenderPhone", xml_helper.T("715-838-4160")),
		xml_helper.E("TransferDate", xml_helper.T(time.Now().Format("2006-01-06"))),
		xml_helper.E("MfrCode", xml_helper.T("BKDK")),
		xml_helper.E("DocumentTitle", xml_helper.T("Trailer Hitches")),
		xml_helper.E("EffectiveDate", xml_helper.T(time.Now().Format("2006-01-06"))),
		xml_helper.E("SubmissionType", xml_helper.T("FULL")),
		xml_helper.E("VcdbVersionDate", xml_helper.T("2013-07-26")))

	doc.Add(header)

	for _, row := range rows {
		data := ACESBaseData{
			Id:                row.Int(vID),
			Qty:               1,
			PartTypeId:        row.Int(typeID),
			ManufacturerLabel: row.Str(sDesc),
			PartNumber:        row.Int(partID),
			BaseVehicleID:     row.Int(baseID),
			SubmodelID:        row.Int(subID),
			ClassId:           row.Int(classId),
			ConfigIds:         strings.Split(row.Str(configIds), ","),
			ConfigNames:       strings.Split(row.Str(typeIds), ","),
			Notes:             strings.Split(row.Str(notes), ","),
			PartNotes:         strings.Split(row.Str(part_notes), ","),
		}

		if data.BaseVehicleID > 0 {
			partEl := xml_helper.E("App",
				xml_helper.A("action", "A"),
				xml_helper.A("id", strconv.Itoa(data.Id)))

			partEl.Add(xml_helper.E("BaseVehicle", xml_helper.A("id", strconv.Itoa(data.BaseVehicleID))))
			if data.SubmodelID > 0 {
				partEl.Add(xml_helper.E("SubModel", xml_helper.A("id", strconv.Itoa(data.SubmodelID))))
			}
			for i := 0; i < len(data.ConfigIds); i++ {
				config := data.ConfigNames[i]
				if config != "" {
					noteEl := xml_helper.E(config, xml_helper.A("id", data.ConfigIds[i]))
					partEl.Add(noteEl)
				}
			}
			for _, note := range data.Notes {
				if note != "" {
					partEl.Add(xml_helper.E("Note", xml_helper.T(note)))
				}
			}
			for _, note := range data.PartNotes {
				if note != "" {
					partEl.Add(xml_helper.E("Note", xml_helper.T(note)))
				}
			}

			partEl.Add(xml_helper.E("Qty", xml_helper.T(strconv.Itoa(data.Qty))))
			partEl.Add(xml_helper.E("PartType", xml_helper.A("id", strconv.Itoa(data.PartTypeId))))
			partEl.Add(xml_helper.E("MfrLabel", xml_helper.T(data.ManufacturerLabel)))
			if data.ClassId == 5 {
				partEl.Add(xml_helper.E("Position", xml_helper.A("id", "22")))
			} else if data.ClassId > 0 && data.ClassId != 11 {
				partEl.Add(xml_helper.E("Position", xml_helper.A("id", "30")))
			}
			partEl.Add(xml_helper.E("Part", xml_helper.T(strconv.Itoa(data.PartNumber))))
			doc.Add(partEl)
		}
	}
	footerEl := xml_helper.E("Footer", xml_helper.E("RecordCount", xml_helper.T(strconv.Itoa(len(rows)))))
	doc.Add(footerEl)
	resp = "<?xml version=\"1.0\" encoding=\"utf-16\" standalone=\"yes\"?>\n" + doc.String()

	data_bytes = []byte(resp)
	redis.RedisMaster.Setex(redis_key, 86400, data_bytes)

	return
}

func GetACESBaseSubmodelData(id int) (baseID, subID int, err error) {
	qry, err := database.GetStatement("GetACESBaseSubmodelData")
	if database.MysqlError(err) {
		return
	}

	row, _, err := qry.ExecFirst(id)
	if database.MysqlError(err) {
		return
	}

	baseID = row.Int(0)
	subID = row.Int(0)

	return
}

func ReverseLookup(partId int) (vehicles []ACESVehicleData, err error) {
	return
}

func GetAcesTypes() (types map[string]string, err error) {
	types = make(map[string]string, 0)
	qry, err := database.GetStatement("GetAcesTypes")
	if database.MysqlError(err) {
		return
	}

	rows, res, err := qry.Exec()
	if database.MysqlError(err) {
		return
	}

	ID := res.Map("ID")
	name := res.Map("name")

	for _, row := range rows {
		types[row.Str(ID)] = row.Str(name)
	}

	return
}
