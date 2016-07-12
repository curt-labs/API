package products

import (
	"errors"
	"github.com/aries-auto/envision-api"
	"github.com/curt-labs/API/helpers/database"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
	"os"
	"sort"
	"strconv"
	"strings"
)

type CatStylePart struct {
	Name   string  `json:"name"`   //collection name
	Styles []Style `json:"styles"` //slice of styles
}

type Style struct {
	Name  string `json:"name"`  //style name
	Parts []Part `json:"parts"` //slice of parts
}

// CategoryStyleParts queries mongo and returns []CatStylePart
// get IconMediaLayer map (part_number:layerID)
// get parts from product_data, mapped (category name:[]Part)
// get Fitment map (part_number:bool)
// for each map key (category), break part array into style structs w/ arrays of parts matching that style
func CategoryStyleParts(v NoSqlVehicle, brands []int, sess *mgo.Session, envision bool) ([]CatStylePart, error) {
	var csps []CatStylePart

	layerChan := make(chan error)
	layerMap := make(map[string]string)
	var err error

	go func() {
		if envision {
			layerMap, err = getIconMediaLayers()
			if err != nil {
				layerChan <- err
				return
			}
		}
		layerChan <- nil
	}()

	collectionToPartsMap, err := v.mapCollectionToParts(sess)
	if err != nil {
		return csps, err
	}

	err = <-layerChan
	if err != nil {
		return csps, err
	}

	var fitmentMap map[string]bool
	if envision {
		fitments, err := v.getIconMediaFitments(collectionToPartsMap)
		if err != nil {
			return csps, err
		}
		fitmentMap, err = mapFitments(fitments)
		if err != nil {
			return csps, err
		}
	}

	styleChan := make(chan []Style, len(collectionToPartsMap))

	for col, parts := range collectionToPartsMap {
		go func() {
			styleChan <- v.parseCollectionStyles(parts, layerMap, fitmentMap)
		}()

		csp := CatStylePart{
			Name:   col,
			Styles: <-styleChan,
		}
		if len(csp.Styles) == 0 {
			continue
		}
		csps = append(csps, csp)
	}
	sort.Sort(ByName(csps))
	return csps, nil
}

// parseCollectionStyles takes an array of parts, iconMedia layer map, iconMedia fitment map
// and returns an []Style for a Vehicle
func (v *NoSqlVehicle) parseCollectionStyles(parts []Part, layerMap map[string]string, fitmentMap map[string]bool) []Style {
	styleMap := make(map[string][]Part)
	for _, part := range parts {
		part.Layer = layerMap[part.PartNumber]
		part.MappedToVehicle = fitmentMap[part.PartNumber]
		for _, pv := range part.Vehicles {
			if strings.ToLower(strings.TrimSpace(pv.Year)) == strings.ToLower(strings.TrimSpace(v.Year)) && strings.ToLower(strings.TrimSpace(pv.Make)) == strings.ToLower(strings.TrimSpace(v.Make)) && strings.ToLower(strings.TrimSpace(pv.Model)) == strings.ToLower(strings.TrimSpace(v.Model)) {
				styleMap[pv.Style] = append(styleMap[pv.Style], part)
			}
		}
	}

	var styles []Style
	for styleName, styleParts := range styleMap {
		styles = append(styles, Style{Name: styleName, Parts: styleParts})
	}
	return styles
}

// mapCollectionToParts queries mongo product database for Vehicle v and
// returns a map of categoryname: []Part for that vehicle
func (v *NoSqlVehicle) mapCollectionToParts(sess *mgo.Session) (map[string][]Part, error) {
	categoryToPartMap := make(map[string][]Part)
	query := bson.M{
		"vehicle_applications.year":  strings.ToLower(v.Year),
		"vehicle_applications.make":  strings.ToLower(v.Make),
		"vehicle_applications.model": strings.ToLower(v.Model),
		"brand.id":                   3,
	}
	var result []Part
	err := sess.DB(database.ProductDatabase).C(database.ProductCollectionName).Find(query).All(&result)
	if err != nil {
		return categoryToPartMap, err
	}
	for _, part := range result {
		for _, cat := range part.Categories {
			if !inCategoryMap(categoryToPartMap[cat.Title], part) {
				categoryToPartMap[cat.Title] = append(categoryToPartMap[cat.Title], part)
			}
		}
	}
	return categoryToPartMap, nil
}

func inCategoryMap(parts []Part, part Part) bool {
	for _, p := range parts {
		if part.PartNumber == p.PartNumber {
			return true
		}
	}
	return false
}

// getIconMediaLayers returns a map of product_number:iconMediaLayerID
func getIconMediaLayers() (map[string]string, error) {
	layerMap := make(map[string]string)

	iconUser := os.Getenv("ICON_USER")
	iconPass := os.Getenv("ICON_PASS")
	iconDomain := os.Getenv("ICON_DOMAIN")
	if iconDomain == "" || iconPass == "" || iconUser == "" {
		return layerMap, errors.New("Missing iCon Credentials")
	}
	conf, err := envisionAPI.NewConfig(iconUser, iconPass, iconDomain)
	if err != nil {
		return layerMap, err
	}

	resp, err := envisionAPI.GetLayers(*conf, "", "")
	if err != nil {
		return layerMap, err
	}
	for _, layer := range resp.Layers {
		layerMap[layer.ProductNumber] = layer.LayerID
	}
	return layerMap, err
}

// GetIconMediaVehicle returns an iconMedia Vehicle for the provided Vehicle v
// looks for a match on style "base"; otherwise, it returns alphabetically the first vehicle style
func (v *NoSqlVehicle) GetIconMediaVehicle() (envisionAPI.Vehicle, error) {
	var iconMediaVehicle envisionAPI.Vehicle
	iconUser := os.Getenv("ICON_USER")
	iconPass := os.Getenv("ICON_PASS")
	iconDomain := os.Getenv("ICON_DOMAIN")
	if iconDomain == "" || iconPass == "" || iconUser == "" {
		return iconMediaVehicle, errors.New("Missing iCon Credentials")
	}
	conf, err := envisionAPI.NewConfig(iconUser, iconPass, iconDomain)
	if err != nil {
		return iconMediaVehicle, err
	}

	resp, err := envisionAPI.GetVehicleByYearMakeModel(*conf, v.Year, v.Make, v.Model)
	if err != nil {
		return iconMediaVehicle, err
	}
	sort.Sort(ByBody(resp.Vehicles))

	iconMediaVehicle = resp.Vehicles[0] //FIRST vehicle
	for _, veh := range resp.Vehicles {
		if strings.ToLower(veh.BodyType) == "base" {
			iconMediaVehicle = veh
			break
		}
	}
	return iconMediaVehicle, err
}

// returns an array of fitments (iconMedia object PartNumber, Fitment) for the Vehicle v
func (v *NoSqlVehicle) getIconMediaFitments(partsMap map[string][]Part) ([]envisionAPI.Fitment, error) {
	var parts []Part
	for _, p := range partsMap {
		parts = append(parts, p...)
	}
	var fitments []envisionAPI.Fitment
	iconUser := os.Getenv("ICON_USER")
	iconPass := os.Getenv("ICON_PASS")
	iconDomain := os.Getenv("ICON_DOMAIN")
	if iconDomain == "" || iconPass == "" || iconUser == "" {
		return fitments, errors.New("Missing iCon Credentials")
	}
	conf, err := envisionAPI.NewConfig(iconUser, iconPass, iconDomain)
	if err != nil {
		return fitments, err
	}

	iconMediaVehicle, err := v.GetIconMediaVehicle() //Currently uses first vehicle (or base) returned by iCon API; styles don't match our styles
	if err != nil {
		return fitments, err
	}
	vehicleID, err := strconv.Atoi(iconMediaVehicle.ID)
	if err != nil {
		return fitments, err
	}

	var partNumbers []string
	for _, p := range parts {
		partNumbers = append(partNumbers, p.PartNumber)
	}
	resp, err := envisionAPI.MatchFitment(*conf, vehicleID, partNumbers...)
	if err != nil {
		if err.Error() == "invalid character ',' looking for beginning of value" { // No matches for vehicle; envision returns invalid json
			return nil, nil
		}
		return fitments, err
	}
	return resp.Fitments, err
}

// mapFitments accepts an array of fitments and returns a map product_number:bool(fitment)
func mapFitments(fitments []envisionAPI.Fitment) (map[string]bool, error) {
	fitmentMap := make(map[string]bool)
	var err error
	for _, fitment := range fitments {
		fitmentMap[fitment.Number], err = strconv.ParseBool(fitment.Mapped)
		if err != nil {
			return fitmentMap, err
		}

	}
	return fitmentMap, nil
}

// Sort Utils
type ByBody []envisionAPI.Vehicle

func (a ByBody) Len() int           { return len(a) }
func (a ByBody) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a ByBody) Less(i, j int) bool { return a[i].BodyType < a[j].BodyType }

type ByName []CatStylePart

func (a ByName) Len() int           { return len(a) }
func (a ByName) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a ByName) Less(i, j int) bool { return a[i].Name < a[j].Name }
