package products

import (
	"encoding/json"
	"fmt"
	"github.com/curt-labs/API/helpers/database"
	"github.com/curt-labs/API/helpers/redis"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
	"sort"
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
// get parts from product_data, mapped (category name:[]Part)
// for each map key (category), break part array into style structs w/ arrays of parts matching that style
func CategoryStyleParts(v NoSqlVehicle, brandArray []int, sess *mgo.Session) ([]CatStylePart, error) {
	var csps []CatStylePart

	redisKey := fmt.Sprintf("%s:%s:%s", v.Year, v.Make, v.Model)
	data, err := redis.Get(redisKey)
	if err == nil && len(data) > 0 {
		err = json.Unmarshal(data, &csps)
		if err == nil {
			return csps, nil
		}
	}

	collectionToPartsMap, err := v.mapCollectionToParts(sess, brandArray)
	if err != nil {
		return csps, err
	}

	styleChan := make(chan []Style, len(collectionToPartsMap))

	for col, parts := range collectionToPartsMap {
		go func() {
			styleChan <- v.parseCollectionStyles(parts)
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
	redis.Setex(redisKey, csps, 60*60*24)
	return csps, nil
}

// parseCollectionStyles takes an array of parts, iconMedia layer map, iconMedia fitment map
// and returns an []Style for a Vehicle
func (v *NoSqlVehicle) parseCollectionStyles(parts []Part) []Style {
	styleMap := make(map[string][]Part)
	for _, part := range parts {
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
func (v *NoSqlVehicle) mapCollectionToParts(sess *mgo.Session, brands []int) (map[string][]Part, error) {
	categoryToPartMap := make(map[string][]Part)
	query := bson.M{
		"vehicle_applications.year":  strings.ToLower(v.Year),
		"vehicle_applications.make":  strings.ToLower(v.Make),
		"vehicle_applications.model": strings.ToLower(v.Model),
		"brand.id":                   bson.M{"$in": brands},
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

type ByName []CatStylePart

func (a ByName) Len() int           { return len(a) }
func (a ByName) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a ByName) Less(i, j int) bool { return a[i].Name < a[j].Name }
