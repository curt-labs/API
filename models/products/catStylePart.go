package products

import (
	"gopkg.in/mgo.v2"
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

type ByName []CatStylePart

func (a ByName) Len() int           { return len(a) }
func (a ByName) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a ByName) Less(i, j int) bool { return a[i].Name < a[j].Name }

// CategoryStyleParts queries mongo and returns []CateStylePart
// get parts from several collections
// for each collection, assign to map[collection name] => []Part
// for each map key (category), break part array into style structs w/ arrays of parts matching that style
func CategoryStyleParts(v NoSqlVehicle, brands []int, sess *mgo.Session) ([]CatStylePart, error) {
	var csps []CatStylePart

	cols, err := GetAriesVehicleCollections(sess)
	if err != nil {
		return csps, err
	}

	collectionToPartsMap := make(map[string][]Part)
	queryMap := make(map[string]interface{})
	queryMap["year"] = strings.ToLower(v.Year)
	queryMap["make"] = strings.ToLower(v.Make)
	queryMap["model"] = strings.ToLower(v.Model)

	type PartError struct {
		Parts []Part
		Err   error
	}
	partChan := make(chan PartError, len(cols))

	for _, col := range cols {
		c := sess.DB(AriesDb).C(col)
		go func() {
			var pe PartError
			var ids []int
			err = c.Find(queryMap).Distinct("parts", &ids)
			if err != nil {
				pe.Err = err
				partChan <- pe
				return
			}
			if len(ids) == 0 {
				partChan <- pe
				return
			}

			pe.Parts, err = GetMany(ids, brands, sess)
			if err != nil {
				pe.Err = err
				partChan <- pe
				return
			}

			partChan <- pe
		}()

		parts := <-partChan
		if parts.Err != nil {
			return csps, err
		}
		if len(parts.Parts) > 0 {
			collectionToPartsMap[col] = parts.Parts
		}

	}

	styleChan := make(chan []Style, len(collectionToPartsMap))

	for col, parts := range collectionToPartsMap {
		go func() {
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
			styleChan <- styles
		}()

		csp := CatStylePart{
			Name:   col,
			Styles: <-styleChan,
		}
		csps = append(csps, csp)
	}
	sort.Sort(ByName(csps))
	return csps, nil
}
