package vinLookup

import (
	"bytes"
	"database/sql"
	"encoding/base64"
	"encoding/xml"
	"errors"
	"github.com/curt-labs/GoAPI/helpers/database"
	"github.com/curt-labs/GoAPI/models/products"
	"io/ioutil"
	// "log"
	"net/http"
	"strconv"
)

type AcesVehicle struct {
	AcesID            int
	AAIABaseVehicleID int
	AAIAMakeID        int
	AAIAModelID       int
	AAIAYearID        int
	AAIASubmodelID    int
	AAIARegionID      int
}
type CurtVehicle struct {
	ID            int
	BaseVehicle   BaseVehicle
	Submodel      Submodel
	Configuration VehicleConfiguration
	Parts         []products.Part
}

type BaseVehicle struct {
	ID        int
	ModelID   int
	MakeID    int
	YearID    int
	ModelName string
	MakeName  string
}

type Submodel struct {
	ID   int
	Name string
}

type VehicleConfiguration struct {
	// Vehicle Vehicle
	TypeID  int
	ValueID int
	Type    string
	Value   string
	// Part    products.Part
}

//reponse
type XMLResponse struct {
	XMLName xml.Name
	Body    Body
}
type Body struct {
	XMLName           xml.Name
	DecodeVinResponse DecodeVinResponse `xml:"decodeVinResponse"`
}

type DecodeVinResponse struct {
	XMLName     xml.Name    `xml:"decodeVinResponse"`
	VinResponse VinResponse `xml:"VinResponse"`
}

type VinResponse struct {
	Vin          string  `xml:"vin"`
	ReturnCode   string  `xml:"returnCode"`
	CorrectedVin string  `xml:"correctedVin"`
	ErrorBytes   string  `xml:"errorBytes"`
	Fields       []Field `xml:"fields"`
}
type Field struct {
	Key   string `xml:"name,attr"`
	Value string `xml:",chardata"`
}

//request
type Envelope struct {
	XMLName xml.Name       `xml:"soapenv:Envelope"`
	SoapEnv string         `xml:"xmlns:soapenv,attr"`
	Web     string         `xml:"xmlns:web,attr"`
	Header  EnvelopeHeader `xml:"soapenv:Header"`
	Body    EnvelopeBody   `xml:"soapenv:Body"`
}

type EnvelopeHeader struct {
}
type EnvelopeBody struct {
	DecodeVin DecodeVin `xml:"web:decodeVin"`
}
type DecodeVin struct {
	Vin             string `xml:"VinRequest>vin"`
	RequestedFields string `xml:"RequestedFields"`
}

var (
	// getVehiclesPreConfig = `SELECT vv.ID, vca.VehicleConfigID, vca.AttributeID, ca.value
	// 							FROM vcdb_Vehicle AS vv
	// 							LEFT JOIN BaseVehicle AS bv ON bv.ID = vv.BaseVehicleID
	// 							LEFT JOIN Submodel AS sm ON sm.ID = vv.SubmodelID
	// 							LEFT JOIN VehicleConfigAttribute AS vca ON vca.VehicleConfigID = vv.ConfigID
	// 							LEFT JOIN ConfigAttribute AS ca ON ca.ID = vca.AttributeID
	// 							WHERE bv.AAIABaseVehicleID = ?
	// 							AND sm.AAIASubmodelID = ?`
	getCurtVehiclesPreConfig = `SELECT vv.ID, vmd.ID,vmd.ModelName, vmk.ID, vmk.MakeName, vyr.YearID, sm.ID, sm.SubmodelName, cat.name, cat.ID, ca.value, ca.ID
								FROM vcdb_Vehicle AS vv
								LEFT JOIN BaseVehicle AS bv ON bv.ID = vv.BaseVehicleID
								LEFT JOIN vcdb_Model AS vmd ON vmd.ID = bv.ModelID
								LEFT JOIN vcdb_Make AS vmk ON vmk.ID = bv.MakeID 
								LEFT JOIN vcdb_year AS vyr ON vyr.YearID = bv.YearID
								LEFT JOIN Submodel AS sm ON sm.ID = vv.SubmodelID
								LEFT JOIN VehicleConfigAttribute AS vca ON vca.VehicleConfigID = vv.ConfigID
								LEFT JOIN ConfigAttribute AS ca ON ca.ID = vca.AttributeID
								LEFT JOIN ConfigAttributeType AS cat ON cat.ID = ca.ConfigAttributeTypeID
								WHERE bv.AAIABaseVehicleID = ?
								AND (sm.AAIASubmodelID = ?  OR sm.AAIASubmodelID IS NULL)`

	// getVehicleConfigIDs = `SELECT vca.AttributeID, vca.* FROM VehicleConfigAttribute AS vca
	// 						LEFT JOIN ConfigAttribute AS ca ON ca.ID = vca.AttributeID
	// 						WHERE ca.ConfigAttributeTypeID = 2
	// 						AND ca.vcdbID = 17` //list of Configs linked to 34
	// getAttributeID = `SELECT ID FROM ConfigAttribute WHERE vcdbID = 17 AND ConfigAttributeTypeID = 2` //34

	// getVehicles = `SELECT vv.ID
	// 					FROM vcdb_Vehicle AS vv
	// 					LEFT JOIN BaseVehicle AS bv ON bv.ID = vv.BaseVehicleID
	// 					LEFT JOIN Submodel AS sm ON sm.ID = vv.SubmodelID
	// 					LEFT JOIN VehicleConfigAttribute AS vca ON vca.VehicleConfigID = vv.ConfigID
	// 					LEFT JOIN ConfigAttribute AS ca ON ca.ID = vca.AttributeID
	// 					WHERE bv.AAIABaseVehicleID = ?
	// 					AND sm.AAIASubmodelID = ?
	// 					AND (vv.configID IN (SELECT vca.VehicleConfigID FROM VehicleConfigAttribute AS vca
	// 					LEFT JOIN ConfigAttribute AS ca ON ca.ID = vca.AttributeID
	// 					WHERE ca.ConfigAttributeTypeID = ?
	// 					AND ca.vcdbID = ?)
	// 					OR vv.configID = 0 OR vv.configID IS NULL)`

	getPartID = `SELECT PartNumber FROM vcdb_VehiclePart WHERE VehicleID = ?`
)

func VinPartLookup(vin string) (vs []CurtVehicle, err error) {
	//get ACES vehicles
	av, err := getAcesVehicle(vin)
	if err != nil {
		return vs, err
	}

	//get CURT vehicle
	curtVehicles, err := av.getCurtVehicles()

	//get parts
	var p products.Part
	for _, v := range curtVehicles {
		//get part id
		db, err := sql.Open("mysql", database.ConnectionString())
		if err != nil {
			return vs, err
		}
		defer db.Close()
		stmt, err := db.Prepare(getPartID)
		if err != nil {
			return vs, err
		}
		defer stmt.Close()
		res, err := stmt.Query(v.ID)
		for res.Next() {
			err = res.Scan(&p.ID)
			if err != nil {
				return vs, err
			}
			//get part -- adds some weight
			err = p.FromDatabase()
			if err != nil {
				return vs, err
			}

			//append to vehicle.parts
			v.Parts = append(v.Parts, p)

		}
		//omit null vehicles (Base, pre-config vehicles with no associated parts)
		if v.Parts != nil {
			vs = append(vs, v)
		}
	}

	// for _, vehicle := range vs {
	// 	var l products.Lookup
	// 	l.Vehicle = vehicle
	// 	partChan := make(chan []products.Part)
	// 	go l.LoadParts(partChan)
	// 	ps = <-partChan
	// }

	return vs, err
}

func GetVehicleConfigs(vin string) (curtVehicles []CurtVehicle, err error) {
	//get ACES vehicles
	av, err := getAcesVehicle(vin)
	if err != nil {
		return curtVehicles, err
	}
	//get CURT vehicle
	curtVehicles, err = av.getCurtVehicles()
	return curtVehicles, err
}

func (v *CurtVehicle) GetPartsFromVehicleConfig() (ps []products.Part, err error) {
	//get parts
	var p products.Part
	//get part id
	db, err := sql.Open("mysql", database.ConnectionString())
	if err != nil {
		return ps, err
	}
	defer db.Close()
	stmt, err := db.Prepare(getPartID)
	if err != nil {
		return ps, err
	}
	defer stmt.Close()
	res, err := stmt.Query(v.ID)
	for res.Next() {
		err = res.Scan(&p.ID)
		if err != nil {
			return ps, err
		}
		//get part -- adds some weight
		err = p.FromDatabase()
		if err != nil {
			return ps, err
		}

		// //append to vehicle.parts
		// v.Parts = append(v.Parts, p)
		ps = append(ps, p)
	}
	//omit null vehicles (Base, pre-config vehicles with no associated parts)
	// if v.Parts != nil {
	// 	vs = append(vs, v)
	// }

	return ps, err
}

func query(vin string) (output []byte, err error) {
	var e Envelope
	e.SoapEnv = "http://schemas.xmlsoap.org/soap/envelope/"
	e.Web = "http://webservice.vindecoder.polk.com/"
	e.Body.DecodeVin.Vin = vin
	e.Body.DecodeVin.RequestedFields = `ACES_BASE_VEHICLE,ACES_MAKE_ID,ACES_MDL_ID,ACES_SUB_MDL_ID,ACES_YEAR_ID,ACES_REGION_ID,ACES_VEHICLE_ID`

	output, err = xml.MarshalIndent(e, " ", "\t")
	if err != nil {
		return output, err
	}
	return output, err
}

func getAcesVehicle(vin string) (av AcesVehicle, err error) {
	data := []byte(database.VintelligencePass())
	password := base64.StdEncoding.EncodeToString(data)

	b, err := query(vin)
	if err != nil {
		return av, err
	}
	buffer := bytes.NewReader(b)
	client := http.Client{}
	req, err := http.NewRequest("POST", "https://vintelligence3.polk.com/vindecoder/VinDecoderService", buffer)
	if err != nil {
		return av, err
	}
	req.Header.Add("Authorization", "Basic "+password)
	req.Header.Add("Content-Type", "text/xml;charset=utf-8")
	req.Header.Add("Host", "\"api.curtmfg.com\"")

	resp, err := client.Do(req)
	if err != nil {
		return av, err
	}

	if resp.StatusCode != 200 {
		err = errors.New(resp.Status)
		return av, err
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return av, err
	}

	var x XMLResponse
	err = xml.Unmarshal(body, &x)
	if err != nil {
		return av, err
	}

	for _, field := range x.Body.DecodeVinResponse.VinResponse.Fields {
		switch field.Key {
		case "ACES_BASE_VEHICLE":
			av.AAIABaseVehicleID, err = strconv.Atoi(field.Value)
		case "ACES_MAKE_ID":
			av.AAIAMakeID, err = strconv.Atoi(field.Value)
		case "ACES_MDL_ID":
			av.AAIAModelID, err = strconv.Atoi(field.Value)
		case "ACES_SUB_MDL_ID":
			av.AAIASubmodelID, err = strconv.Atoi(field.Value)
		case "ACES_YEAR_ID":
			av.AAIAYearID, err = strconv.Atoi(field.Value)
		case "ACES_REGION_ID":
			av.AAIARegionID, err = strconv.Atoi(field.Value)
		case "ACES_VEHICLE_ID":
			av.AcesID, err = strconv.Atoi(field.Value)
		}
	}
	return av, err
}

func (av *AcesVehicle) getCurtVehicles() (cvs []CurtVehicle, err error) { //get CURT vehicles
	db, err := sql.Open("mysql", database.ConnectionString())
	if err != nil {
		return cvs, err
	}
	defer db.Close()
	stmt, err := db.Prepare(getCurtVehiclesPreConfig)
	if err != nil {
		return cvs, err
	}
	defer stmt.Close()
	res, err := stmt.Query(av.AAIABaseVehicleID, av.AAIASubmodelID)

	var sub, configKey, configValue *string
	var subID, configKeyID, configValueID *int
	var cv CurtVehicle
	for res.Next() {
		err = res.Scan(
			&cv.ID,
			&cv.BaseVehicle.ModelID,
			&cv.BaseVehicle.ModelName,
			&cv.BaseVehicle.MakeID,
			&cv.BaseVehicle.MakeName,
			&cv.BaseVehicle.YearID,
			&subID,
			&sub,
			&configKey,
			&configKeyID,
			&configValue,
			&configValueID,
		)
		if subID != nil {
			cv.Submodel.ID = *subID
		}
		if sub != nil {
			cv.Submodel.Name = *sub
		}
		if configKey != nil {
			cv.Configuration.Type = *configKey
		}
		if configValue != nil {
			cv.Configuration.Value = *configValue
		}
		if configKeyID != nil {
			cv.Configuration.TypeID = *configKeyID
		}
		if configValueID != nil {
			cv.Configuration.ValueID = *configValueID
		}
		cvs = append(cvs, cv)
	}
	return cvs, err
}
