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
		fitmentMap, err = v.getIconMediaFitmentMap(collectionToPartsMap)
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
func (v *NoSqlVehicle) GetIconMediaVehicle() (envisionAPI.Vehicle, []string, error) {
	var iconMediaVehicle envisionAPI.Vehicle
	iconUser := os.Getenv("ICON_USER")
	iconPass := os.Getenv("ICON_PASS")
	iconDomain := os.Getenv("ICON_DOMAIN")
	if iconDomain == "" || iconPass == "" || iconUser == "" {
		return iconMediaVehicle, []string{}, errors.New("Missing iCon Credentials")
	}
	conf, err := envisionAPI.NewConfig(iconUser, iconPass, iconDomain)
	if err != nil {
		return iconMediaVehicle, []string{}, err
	}

	resp, err := envisionAPI.GetVehicleByYearMakeModel(*conf, v.Year, v.Make, v.Model)
	if err != nil {
		return iconMediaVehicle, []string{}, err
	}

	// VEHICLE PREFERENCE find vehicle with the MOST parts
	var mostParts []string
	var vehicleWithMostParts envisionAPI.Vehicle
	for _, iconVehicle := range resp.Vehicles {
		vehicleID, err := strconv.Atoi(iconVehicle.ID)
		if err != nil {
			return iconMediaVehicle, []string{}, err
		}

		partNumbers, err := getPartsAttachedToVehicleImages(vehicleID)
		if err != nil {
			return iconMediaVehicle, partNumbers, err
		}
		if len(partNumbers) > len(mostParts) {
			mostParts = partNumbers
			vehicleWithMostParts = iconVehicle
		}
	}
	return vehicleWithMostParts, mostParts, err
}

// returns an array of partNumbers for the Vehicle v
func getPartsAttachedToVehicleImages(vehicleID int) ([]string, error) {
	var partNumbers []string
	iconUser := os.Getenv("ICON_USER")
	iconPass := os.Getenv("ICON_PASS")
	iconDomain := os.Getenv("ICON_DOMAIN")
	if iconDomain == "" || iconPass == "" || iconUser == "" {
		return partNumbers, errors.New("Missing iCon Credentials")
	}
	conf, err := envisionAPI.NewConfig(iconUser, iconPass, iconDomain)
	if err != nil {
		return partNumbers, err
	}

	vehicleProductResponse, err := envisionAPI.GetVehicleProducts(*conf, vehicleID)
	if err != nil {
		return partNumbers, err
	}
	for _, partNumber := range vehicleProductResponse.Numbers {
		partNumbers = append(partNumbers, partNumber.Number)
	}
	return partNumbers, err
}

// returns an array of fitments (iconMedia object PartNumber, Fitment) for the Vehicle v
func (v *NoSqlVehicle) getIconMediaFitmentMap(partsMap map[string][]Part) (map[string]bool, error) {
	var parts []Part
	fitments := make(map[string]bool)
	for _, p := range partsMap {
		parts = append(parts, p...)
	}

	_, partNumbers, err := v.GetIconMediaVehicle() //Currently vehicle with the MOST MAPPED PARTS returned by iCon API; styles don't match our styles
	if err != nil {
		return fitments, err
	}

	for _, p := range parts {
		if inArray(partNumbers, p.PartNumber) {
			fitments[p.PartNumber] = true
		} else {
			fitments[p.PartNumber] = false
		}
	}
	return fitments, err
}

func inArray(arr []string, a string) bool {
	for _, ar := range arr {
		if ar == a {
			return true
		}
	}
	return false
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
