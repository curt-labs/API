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
	TypeID      int
	ValueID     int
	Type        string
	Value       string
	AcesValueID int
}

type ConfigurationBits struct {
	WheelBase                        interface{} //WHL_BAS_SHRST_INCHS
	BodyType                         interface{} //ACES_BODY_TYPE
	DriveType                        interface{} //ACES_DRIVE_ID
	NumberOfDoors                    interface{} //DOOR_CNT
	FuelType                         interface{}
	Engine                           interface{} //ACES_LITERS + ACES_CYLINDERS--not quite
	Aspiration                       interface{} //ACES_ASP_ID
	BedLength                        interface{} //TRK_BED_LEN_CD
	BedType                          interface{}
	BrakeABS                         interface{}
	BrakeSystem                      interface{}
	CylinderHeadType                 interface{}
	EngineDesignation                interface{}
	EngineManufacturer               interface{}
	EngineVersion                    interface{}
	EngineVin                        interface{} //ACES_ENG_VIN_ID
	FrontBrakeType                   interface{}
	FrontSpringType                  interface{}
	FuelDeliverySubType              interface{}
	FuelDeliveryType                 interface{} //ACES_FUEL
	FuelSystemControlType            interface{}
	FuelSystemDesign                 interface{} //ACES_FUEL
	IgnitionSystemDesign             interface{}
	ManufacturerBodyCode             interface{}
	PowerOutput                      interface{}
	RearBrakeType                    interface{}
	RearSpringType                   interface{}
	SteeringSystem                   interface{}
	SteeringType                     interface{}
	TransmissionElectronicControlled interface{}
	Transmission                     interface{} //TRANS_CD
	TransmissionControlType          interface{}
	TransmissionManufacturerCode     interface{}
	TransmissionNumberOfSpeeds       interface{} //TRANS_OPT1_SPEED_CD
	TransmissionType                 interface{}
	ValvesPerEngine                  interface{}
	Region                           interface{} //ACES_REGION_ID
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
	getCurtVehiclesPreConfig = `SELECT vv.ID, vmd.ID,vmd.ModelName, vmk.ID, vmk.MakeName, vyr.YearID, sm.ID, sm.SubmodelName, cat.name, cat.ID, ca.value, ca.ID, ca.vcdbID
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
								AND (sm.AAIASubmodelID = ?  OR sm.AAIASubmodelID IS NULL) `

	getPartID = `SELECT PartNumber FROM vcdb_VehiclePart WHERE VehicleID = ?`
)

const (
	soapRequestedFields = `ACES_BASE_VEHICLE,ACES_MAKE_ID,ACES_MDL_ID,ACES_SUB_MDL_ID,ACES_YEAR_ID,ACES_REGION_ID,ACES_VEHICLE_ID,
		ACES_FUEL,ACES_FUEL_DELIVERY,ACES_ENG_VIN_ID,ACES_ASP_ID,ACES_DRIVE_ID,ACES_BODY_TYPE,ACES_REGION_ID,ACES_LITERS,ACES_CC_DISPLACEMENT,ACES_CI_DISPLACEMENT,
		ACES_CYLINDERS,ACES_RESERVED,DOOR_CNT,BODY_STYLE_DESC,WHL_BAS_SHRST_INCHS,TRK_BED_LEN_DESC,TRANS_CD,TRK_BED_LEN_CD,ENG_FUEL_DESC`
)

func VinPartLookup(vin string) (vs []CurtVehicle, err error) {
	//get ACES vehicles
	av, configMap, err := getAcesVehicle(vin)
	if err != nil {
		return vs, err
	}

	//get CURT vehicle
	curtVehicles, err := av.getCurtVehicles(configMap)

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

	return vs, err
}

func GetVehicleConfigs(vin string) (curtVehicles []CurtVehicle, err error) {
	//get ACES vehicles
	av, configMap, err := getAcesVehicle(vin)
	if err != nil {
		return curtVehicles, err
	}
	//get CURT vehicle
	curtVehicles, err = av.getCurtVehicles(configMap)
	return curtVehicles, err
}

//already have vehicleID (vcdb_vehicle.ID)? get parts
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

		ps = append(ps, p)
	}
	return ps, err
}

func query(vin string) (output []byte, err error) {
	var e Envelope
	e.SoapEnv = "http://schemas.xmlsoap.org/soap/envelope/"
	e.Web = "http://webservice.vindecoder.polk.com/"
	e.Body.DecodeVin.Vin = vin
	e.Body.DecodeVin.RequestedFields = soapRequestedFields

	output, err = xml.MarshalIndent(e, " ", "\t")
	if err != nil {
		return output, err
	}
	return output, err
}

func getAcesVehicle(vin string) (av AcesVehicle, configMap map[int]interface{}, err error) {
	data := []byte(database.VintelligencePass())
	password := base64.StdEncoding.EncodeToString(data)

	b, err := query(vin)
	if err != nil {
		return av, configMap, err
	}
	buffer := bytes.NewReader(b)
	client := http.Client{}
	req, err := http.NewRequest("POST", "https://vintelligence3.polk.com/vindecoder/VinDecoderService", buffer)
	if err != nil {
		return av, configMap, err
	}
	req.Header.Add("Authorization", "Basic "+password)
	req.Header.Add("Content-Type", "text/xml;charset=utf-8")
	req.Header.Add("Host", "\"api.curtmfg.com\"")

	resp, err := client.Do(req)
	if err != nil {
		return av, configMap, err
	}

	if resp.StatusCode != 200 {
		err = errors.New(resp.Status)
		return av, configMap, err
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return av, configMap, err
	}
	// log.Print(string(body))

	var x XMLResponse
	err = xml.Unmarshal(body, &x)
	if err != nil {
		return av, configMap, err
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
	//check out them configs
	configMap, err = av.checkConfigs(x.Body.DecodeVinResponse.VinResponse.Fields)

	return av, configMap, err
}

//creates a map of config options from the SOAP request to check against curt vehicles
func (av *AcesVehicle) checkConfigs(responseFields []Field) (configMap map[int]interface{}, err error) {
	//map of configAttributeType AcesID to configAttribute Aces ID
	configMap = make(map[int]interface{})
	for _, field := range responseFields {
		switch field.Key {
		case "WHL_BAS_SHRST_INCHS":
			configMap[1], err = strconv.Atoi(field.Value)
		case "ACES_BODY_TYPE":
			configMap[2], err = strconv.Atoi(field.Value)
		case "ACES_DRIVE_ID":
			configMap[3], err = strconv.Atoi(field.Value)
		case "DOOR_CNT":
			configMap[4], err = strconv.Atoi(field.Value)
		case "ACES_ASP_ID":
			configMap[8], err = strconv.Atoi(field.Value)
		case "ACES_ENG_VIN_ID":
			configMap[16], err = strconv.Atoi(field.Value)
		case "ACES_FUEL":
			configMap[20], err = strconv.Atoi(field.Value)
		case "TRANS_CD":
			configMap[34] = field.Value
		case "TRANS_OPT1_SPEED_CD":
			configMap[38], err = strconv.Atoi(field.Value)

			if err != nil {
				return configMap, err
			}
		}

	}
	// log.Print("ALL CBS", configMap)
	return configMap, err
}

func (av *AcesVehicle) getCurtVehicles(configMap map[int]interface{}) (cvs []CurtVehicle, err error) { //get CURT vehicles

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
	if err != nil {
		return cvs, err
	}

	var sub, configKey, configValue *string
	var subID, configKeyID, configValueID, acesConfigValID *int
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
			&acesConfigValID,
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
		if acesConfigValID != nil {
			cv.Configuration.AcesValueID = *acesConfigValID
		}

		cvs = append(cvs, cv)
		//Pop off configs that have non-conforming values
		// log.Print("CONFIG TYP", cv.Configuration.TypeID)
		if name, ok := configMap[cv.Configuration.TypeID]; ok {
			// log.Print("CONFIG VAL ", name, " ACES CONVFIG VAL ", cv.Configuration.AcesValueID)
			if cv.Configuration.AcesValueID != name {
				cvs = cvs[:len(cvs)-1]
			}
		}

	}
	return cvs, err
}
