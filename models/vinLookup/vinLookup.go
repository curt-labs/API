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
	"log"
	"net/http"
	"strconv"
)

//vehicle
type Vehicle struct {
	ID             int
	YearID         int
	MakeID         int
	AAIAMakeID     int
	ModelID        int
	AAIAModelID    int
	BaseVehicleID  int
	SubmodelID     int
	AAIASubmodelID int
	AcesID         int
	Year           int
	Make           string
	Model          string
	Submodel       string
	RegionID       int
	BodyTypeID     int
	Bumper         interface{}
	BedLength      interface{}
	DriveType      interface{}
}

type VehicleConfiguration struct {
	Vehicle Vehicle
	TypeID  int
	ValueID int
	Type    string
	Value   string
	Part    products.Part
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
	getVehiclesPreConfig = `SELECT vv.ID, vca.VehicleConfigID, vca.AttributeID, ca.value
								FROM vcdb_Vehicle AS vv
								LEFT JOIN BaseVehicle AS bv ON bv.ID = vv.BaseVehicleID
								LEFT JOIN Submodel AS sm ON sm.ID = vv.SubmodelID
								LEFT JOIN VehicleConfigAttribute AS vca ON vca.VehicleConfigID = vv.ConfigID
								LEFT JOIN ConfigAttribute AS ca ON ca.ID = vca.AttributeID
								WHERE bv.AAIABaseVehicleID = ?
								AND sm.AAIASubmodelID = ?`

	// getVehicleConfigIDs = `SELECT vca.AttributeID, vca.* FROM VehicleConfigAttribute AS vca
	// 						LEFT JOIN ConfigAttribute AS ca ON ca.ID = vca.AttributeID
	// 						WHERE ca.ConfigAttributeTypeID = 2
	// 						AND ca.vcdbID = 17` //list of Configs linked to 34
	// getAttributeID = `SELECT ID FROM ConfigAttribute WHERE vcdbID = 17 AND ConfigAttributeTypeID = 2` //34

	getVehicles = `SELECT vv.ID
						FROM vcdb_Vehicle AS vv
						LEFT JOIN BaseVehicle AS bv ON bv.ID = vv.BaseVehicleID
						LEFT JOIN Submodel AS sm ON sm.ID = vv.SubmodelID
						LEFT JOIN VehicleConfigAttribute AS vca ON vca.VehicleConfigID = vv.ConfigID
						LEFT JOIN ConfigAttribute AS ca ON ca.ID = vca.AttributeID
						WHERE bv.AAIABaseVehicleID = ?
						AND sm.AAIASubmodelID = ?
						AND (vv.configID IN (SELECT vca.VehicleConfigID FROM VehicleConfigAttribute AS vca 
						LEFT JOIN ConfigAttribute AS ca ON ca.ID = vca.AttributeID
						WHERE ca.ConfigAttributeTypeID = ? 
						AND ca.vcdbID = ?) 
						OR vv.configID = 0 OR vv.configID IS NULL)`

	getPartID = `SELECT PartNumber FROM vcdb_VehiclePart WHERE VehicleID = ?`
)

func CreateQuery(vin string) (output []byte, err error) {
	var e Envelope
	e.SoapEnv = "http://schemas.xmlsoap.org/soap/envelope/"
	e.Web = "http://webservice.vindecoder.polk.com/"
	e.Body.DecodeVin.Vin = vin
	e.Body.DecodeVin.RequestedFields = "ACES_BASE_VEHICLE,ACES_MAKE_ID,ACES_MDL_ID,ACES_SUB_MDL_ID,ACES_YEAR_ID,ACES_BODY_TYPE,MAK_NM,MDL_DESC,TRIM_DESC,ACES_REGION_ID,ACES_VEHICLE_ID"

	output, err = xml.MarshalIndent(e, " ", "\t")
	if err != nil {
		return output, err
	}
	return output, err
}

func Lookup(vin string) (vcs []VehicleConfiguration, err error) {
	var v VehicleConfiguration
	data := []byte(database.VintelligencePass())
	password := base64.StdEncoding.EncodeToString(data)

	b, err := CreateQuery(vin)
	if err != nil {
		return vcs, err
	}
	buffer := bytes.NewReader(b)
	client := http.Client{}
	req, err := http.NewRequest("POST", "https://vintelligence3.polk.com/vindecoder/VinDecoderService", buffer)
	if err != nil {
		return vcs, err
	}
	req.Header.Add("Authorization", "Basic "+password)
	req.Header.Add("Content-Type", "text/xml;charset=utf-8")
	req.Header.Add("Host", "\"api.curtmfg.com\"")

	resp, err := client.Do(req)
	if err != nil {
		return vcs, err
	}

	if resp.StatusCode != 200 {
		err = errors.New(resp.Status)
		return vcs, err
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return vcs, err
	}
	log.Print(string(body))

	var x XMLResponse
	err = xml.Unmarshal(body, &x)
	if err != nil {
		return vcs, err
	}

	for _, field := range x.Body.DecodeVinResponse.VinResponse.Fields {
		switch field.Key {
		case "ACES_BASE_VEHICLE":
			v.Vehicle.BaseVehicleID, err = strconv.Atoi(field.Value)
		case "ACES_MAKE_ID":
			v.Vehicle.AAIAMakeID, err = strconv.Atoi(field.Value)
		case "ACES_MDL_ID":
			v.Vehicle.AAIAModelID, err = strconv.Atoi(field.Value)
		case "ACES_SUB_MDL_ID":
			v.Vehicle.AAIASubmodelID, err = strconv.Atoi(field.Value)
		case "ACES_YEAR_ID":
			v.Vehicle.YearID, err = strconv.Atoi(field.Value)
			v.Vehicle.Year, err = strconv.Atoi(field.Value)
		case "ACES_BODY_TYPE":
			v.TypeID = 2 //Body Type
			v.ValueID, err = strconv.Atoi(field.Value)
			// v.Configurations = append(v.Configurations, vcBody)
			v.Vehicle.BodyTypeID, err = strconv.Atoi(field.Value)
		case "MAK_NM":
			v.Vehicle.Make = field.Value
		case "MDL_DESC":
			v.Vehicle.Model = field.Value
			// case "TRIM_DESC":
			// 	v.BaseVehicleID, err = strconv.Atoi(field.Value)
		case "ACES_REGION_ID":
			v.Vehicle.RegionID, err = strconv.Atoi(field.Value)
		case "ACES_VEHICLE_ID":
			v.Vehicle.AcesID, err = strconv.Atoi(field.Value)
		case "BUMPER":
			v.Vehicle.Bumper = field.Value
		case "BED_LENGTH":
			v.Vehicle.BedLength = field.Value
		case "ACES_DRIVE_TYPE":
			v.Vehicle.DriveType = field.Value

		}

	}

	vcs, err = v.partLookup()
	log.Print(vcs)

	return vcs, err
}

//get curt vehicle, then parts
func (v *VehicleConfiguration) partLookup() (vcs []VehicleConfiguration, err error) {
	// get vehicle ID
	var vs []VehicleConfiguration

	db, err := sql.Open("mysql", database.ConnectionString())
	if err != nil {
		return vcs, err
	}
	defer db.Close()
	stmt, err := db.Prepare(getVehiclesPreConfig)
	if err != nil {
		return vcs, err
	}
	defer stmt.Close()

	res, err := stmt.Query(v.Vehicle.BaseVehicleID, v.Vehicle.AAIASubmodelID)
	// var tempVehicle Vehicle
	var tempCon VehicleConfiguration
	var valID *int
	var val *string
	var configID *int
	for res.Next() {
		tempCon = *v
		err = res.Scan(&tempCon.Vehicle.ID, &configID, &valID, &val)
		if err != nil {
			return vcs, err
		}
		if valID != nil {
			tempCon.ValueID = *valID
		}
		if val != nil {
			tempCon.Value = *val
		}
		vs = append(vs, tempCon)
	}
	//not enough configs coming from polk - returns a variety of configs from a single vehicle

	for _, vc := range vs {
		var p products.Part
		db, err := sql.Open("mysql", database.ConnectionString())
		if err != nil {
			return vcs, err
		}
		defer db.Close()
		stmt, err := db.Prepare(getPartID)
		if err != nil {
			return vcs, err
		}
		defer stmt.Close()
		err = stmt.QueryRow(vc.Vehicle.ID).Scan(&p.ID)
		err = p.Basics()
		vc.Part = p
		vcs = append(vcs, vc)
	}
	return vcs, err
}
