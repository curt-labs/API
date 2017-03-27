package products

import (
	"encoding/json"
	"fmt"
	"log"
	"sort"
	"strings"

	"gopkg.in/mgo.v2/bson"

	"github.com/curt-labs/API/helpers/apicontext"
	"github.com/curt-labs/API/helpers/database"
	"github.com/curt-labs/API/helpers/redis"
)

var (
	statuses = []int{700, 800, 810, 815, 850, 870, 888, 900, 910, 950}
)

const (
	CURT_LOOKUP_KEY     = "curtlookup:"
	CL_YEARS_KEY        = CURT_LOOKUP_KEY + "years"
	CL_YEARS_PARTS_KEY  = CL_YEARS_KEY + ":parts"
	CL_MAKES_KEY        = ":makes"
	CL_MAKES_PARTS_KEY  = CL_MAKES_KEY + ":parts"
	CL_MODELS_KEY       = ":models"
	CL_MODELS_PARTS_KEY = CL_MODELS_KEY + ":parts"
)

type CurtVehicle struct {
	Year            string      `json:"year,omitempty" xml:"year, omitempty"`
	Make            string      `json:"make,omitempty" xml:"make, omitempty"`
	Model           string      `json:"model,omitempty" xml:"model, omitempty"`
	Style           string      `json:"style,omitempty" xml:"style, omitempty"`
	Parts           []BasicPart `json:"parts,omitempty" xml:"parts, omitempty"`
	PartIdentifiers []int       `json:"parts_ids" xml:"-"`
}

type CurtLookup struct {
	Years  []string `json:"available_years,omitempty" xml:"available_years, omitempty"`
	Makes  []string `json:"available_makes,omitempty" xml:"available_makes, omitempty"`
	Models []string `json:"available_models,omitempty" xml:"available_models, omitempty"`
	Styles []string `json:"available_styles,omitempty" xml:"available_styles, omitempty"`
	Parts  []Part   `json:"parts,omitempty" xml:"parts, omitempty"`
	CurtVehicle
}

func (c *CurtLookup) GetYears(heavyduty bool) error {
	// See if we have the years in redis
	yearsData, err := redis.Get(CL_YEARS_KEY)
	if err == nil {
		// We also need the "years parts"
		yearPartsdata, err := redis.Get(CL_YEARS_PARTS_KEY)
		if err == nil {
			// Unmarshall the data
			err = json.Unmarshal(yearsData, &c.Years)
			err = json.Unmarshal(yearPartsdata, &c.PartIdentifiers)
			//  If we have data
			if (len(c.Years) > 0) && (len(c.PartIdentifiers) > 0) {
				// Exit using the data we found in redis
				return nil
			}
		}
	}

	log.Println("cl.GetYears - missed cache ", CL_YEARS_KEY, " or ", CL_YEARS_PARTS_KEY)

	err = database.Init()
	if err != nil {
		return err
	}
	session := database.ProductMongoSession.Copy()
	defer session.Close()

	col := session.DB(database.ProductDatabase).C(database.ProductCollectionName)

	qry := bson.M{
		"status": bson.M{
			"$in": statuses,
		},
		"vehicle_applications.0": bson.M{
			"$exists": true,
		},
		"brand.id": 1,
	}

	type YearResp struct {
		Apps []VehicleApplication `bson:"vehicle_applications"`
		ID   int                  `bson:"id"`
	}
	var resp []YearResp
	err = col.Find(qry).Select(bson.M{
		"vehicle_applications.year": 1,
		"id":  1,
		"_id": -1,
	}).All(&resp)
	if err != nil {
		return err
	}

	var years []string

	existing := make(map[string]string, 0)
	existingIDS := make(map[int]int, 0)
	for _, app := range resp {
		if _, ok := existingIDS[app.ID]; !ok {
			c.PartIdentifiers = append(c.PartIdentifiers, app.ID)
			existingIDS[app.ID] = app.ID
		}
		for _, a := range app.Apps {
			if _, ok := existing[a.Year]; !ok {
				years = append(years, a.Year)
				existing[a.Year] = a.Year
			}
		}
	}
	c.Years = years

	sort.Sort(sort.Reverse(sort.StringSlice(c.Years)))

	// Set the keys to expire after a week
	go redis.Setex(CL_YEARS_KEY, c.Years, 604800)
	go redis.Setex(CL_YEARS_PARTS_KEY, c.PartIdentifiers, 604800)

	return nil
}

func (c *CurtLookup) GetMakes(heavyduty bool) error {
	// Build the redis keys
	redisMakesKey := fmt.Sprintf("%s:%s%s", CL_YEARS_KEY, c.Year, CL_MAKES_KEY)
	redisMakesPartsKey := fmt.Sprintf("%s:%s%s", CL_YEARS_KEY, c.Year, CL_MAKES_PARTS_KEY)

	// Pull from redis
	makesData, err := redis.Get(redisMakesKey)
	if err == nil {
		// We also need the "makes parts"
		makesPartsData, err := redis.Get(redisMakesPartsKey)
		if err == nil {
			// Unmarshall the data
			err = json.Unmarshal(makesData, &c.Makes)
			err = json.Unmarshal(makesPartsData, &c.PartIdentifiers)
			//  If we have data
			if (len(c.Makes) > 0) && (len(c.PartIdentifiers) > 0) {
				// Exit using the data we found in redis
				return nil
			}
		}
	}

	log.Println("cl.GetMakes - missed cache ", redisMakesKey, " or ", redisMakesPartsKey)

	err = database.Init()
	if err != nil {
		return err
	}
	session := database.ProductMongoSession.Copy()
	defer session.Close()

	col := session.DB(database.ProductDatabase).C(database.ProductCollectionName)

	qry := bson.M{
		"status": bson.M{
			"$in": statuses,
		},
		"vehicle_applications": bson.M{
			"$elemMatch": bson.M{
				"year": c.Year,
			},
		},
		"vehicle_applications.0": bson.M{
			"$exists": true,
		},
		"brand.id": 1,
	}

	type YearResp struct {
		Apps []VehicleApplication `bson:"vehicle_applications"`
		ID   int                  `bson:"id"`
	}
	var resp []YearResp
	err = col.Find(qry).Select(bson.M{
		"vehicle_applications": 1,
		"id":  1,
		"_id": -1,
	}).All(&resp)
	if err != nil {
		return err
	}

	var makes []string

	existing := make(map[string]string, 0)
	existingIDS := make(map[int]int, 0)
	for _, app := range resp {
		if _, ok := existingIDS[app.ID]; !ok {
			c.PartIdentifiers = append(c.PartIdentifiers, app.ID)
			existingIDS[app.ID] = app.ID
		}
		for _, a := range app.Apps {
			if a.Year != c.Year {
				continue
			}
			a.Make = strings.Title(a.Make)
			if _, ok := existing[a.Make]; !ok {
				makes = append(makes, a.Make)
				existing[a.Make] = a.Make
			}
		}
	}
	c.Makes = makes

	sort.Strings(c.Makes)

	// Set the keys to expire after a week
	go redis.Setex(redisMakesKey, c.Makes, 604800)
	go redis.Setex(redisMakesPartsKey, c.PartIdentifiers, 604800)

	return nil
}

func (c *CurtLookup) GetModels(heavyduty bool) error {
	// Build the redis keys
	redisModelsKey := fmt.Sprintf("%s:%s%s:%s%s",
		CL_YEARS_KEY, c.Year, CL_MAKES_KEY, c.Make, CL_MODELS_KEY)
	redisModelsPartsKey := fmt.Sprintf("%s:%s%s:%s%s",
		CL_YEARS_KEY, c.Year, CL_MAKES_KEY, c.Make, CL_MODELS_PARTS_KEY)

	// Pull from redis
	modelsData, err := redis.Get(redisModelsKey)
	if err == nil {
		// We also need the "models parts"
		modelsPartsData, err := redis.Get(redisModelsPartsKey)
		if err == nil {
			// Unmarshall the data
			err = json.Unmarshal(modelsData, &c.Models)
			err = json.Unmarshal(modelsPartsData, &c.PartIdentifiers)
			//  If we have data
			if (len(c.Models) > 0) && (len(c.PartIdentifiers) > 0) {
				// Exit using the data we found in redis
				return nil
			}
		}
	}

	log.Println("cl.GetModels - missed cache ", redisModelsKey, " or ", redisModelsPartsKey)

	err = database.Init()
	if err != nil {
		return err
	}
	session := database.ProductMongoSession.Copy()
	defer session.Close()

	col := session.DB(database.ProductDatabase).C(database.ProductCollectionName)

	qry := bson.M{
		"status": bson.M{
			"$in": statuses,
		},
		"vehicle_applications": bson.M{
			"$elemMatch": bson.M{
				"year": c.Year,
				"make": bson.RegEx{
					Pattern: "^" + c.Make + "$",
					Options: "i",
				},
			},
		},
		"vehicle_applications.0": bson.M{
			"$exists": true,
		},
		"brand.id": 1,
	}

	type YearResp struct {
		Apps []VehicleApplication `bson:"vehicle_applications"`
		ID   int                  `bson:"id"`
	}
	var resp []YearResp
	err = col.Find(qry).Select(bson.M{
		"vehicle_applications": 1,
		"id":  1,
		"_id": -1,
	}).All(&resp)
	if err != nil {
		return err
	}

	var models []string

	existing := make(map[string]string, 0)
	existingIDS := make(map[int]int, 0)
	for _, app := range resp {
		if _, ok := existingIDS[app.ID]; !ok {
			c.PartIdentifiers = append(c.PartIdentifiers, app.ID)
			existingIDS[app.ID] = app.ID
		}
		for _, a := range app.Apps {
			if a.Year != c.Year || strings.ToUpper(a.Make) != strings.ToUpper(c.Make) {
				continue
			}
			a.Model = strings.Title(a.Model)
			if _, ok := existing[a.Model]; !ok {
				models = append(models, a.Model)
				existing[a.Model] = a.Model
			}
		}
	}
	c.Models = models

	sort.Strings(c.Models)

	// Set the keys to expire after a week
	go redis.Setex(redisModelsKey, c.Models, 604800)
	go redis.Setex(redisModelsPartsKey, c.PartIdentifiers, 604800)

	return nil
}

func (c *CurtLookup) GetStyles(heavyduty bool) error {

	err := database.Init()
	if err != nil {
		return err
	}
	session := database.ProductMongoSession.Copy()
	defer session.Close()

	col := session.DB(database.ProductDatabase).C(database.ProductCollectionName)

	qry := bson.M{
		"status": bson.M{
			"$in": statuses,
		},
		"vehicle_applications": bson.M{
			"$elemMatch": bson.M{
				"year": c.Year,
				"make": bson.RegEx{
					Pattern: "^" + c.Make + "$",
					Options: "i",
				},
				"model": bson.RegEx{
					Pattern: "^" + c.Model + "$",
					Options: "i",
				},
			},
		},
		"vehicle_applications.0": bson.M{
			"$exists": true,
		},
		"brand.id": 1,
	}

	type YearResp struct {
		Apps []VehicleApplication `bson:"vehicle_applications"`
		ID   int                  `bson:"id"`
	}
	var resp []YearResp
	err = col.Find(qry).Select(bson.M{
		"vehicle_applications": 1,
		"id":  1,
		"_id": -1,
	}).All(&resp)
	if err != nil {
		return err
	}

	var styles []string

	existing := make(map[string]string, 0)
	existingIDS := make(map[int]int, 0)
	for _, app := range resp {
		if _, ok := existingIDS[app.ID]; !ok {
			c.PartIdentifiers = append(c.PartIdentifiers, app.ID)
			existingIDS[app.ID] = app.ID
		}
		for _, a := range app.Apps {
			if a.Year != c.Year || strings.ToUpper(a.Make) != strings.ToUpper(c.Make) || strings.ToUpper(a.Model) != strings.ToUpper(c.Model) {
				continue
			}

			a.Style = strings.Title(a.Style)
			if _, ok := existing[a.Style]; !ok {
				styles = append(styles, a.Style)
				existing[a.Style] = a.Style
			}
		}
	}
	c.Styles = styles

	sort.Strings(c.Styles)

	return nil
}

func (c *CurtLookup) GetParts(dtx *apicontext.DataContext, heavyduty bool) error {
	err := database.Init()
	if err != nil {
		return err
	}
	session := database.ProductMongoSession.Copy()
	defer session.Close()

	col := session.DB(database.ProductDatabase).C(database.ProductCollectionName)

	appQuery := bson.M{
		"$elemMatch": bson.M{
			"year": c.Year,
			"make": bson.RegEx{
				Pattern: "^" + c.Make + "$",
				Options: "i",
			},
			"model": bson.RegEx{
				Pattern: "^" + c.Model + "$",
			},
		},
	}

	if c.Style != "" {
		appQuery = bson.M{
			"$elemMatch": bson.M{
				"year": c.Year,
				"make": bson.RegEx{
					Pattern: "^" + c.Make + "$",
					Options: "i",
				},
				"model": bson.RegEx{
					Pattern: "^" + c.Model + "$",
				},
				"style": bson.RegEx{
					Pattern: "^" + c.Style + "$",
				},
			},
		}
	}

	qry := bson.M{
		"status": bson.M{
			"$in": statuses,
		},
		"vehicle_applications": appQuery,
		"vehicle_applications.0": bson.M{
			"$exists": true,
		},
		"brand.id": 1,
	}

	err = col.Find(qry).All(&c.Parts)

	return err
}
