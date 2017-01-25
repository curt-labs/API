package products

import (
	"sort"
	"strings"

	"gopkg.in/mgo.v2/bson"

	"github.com/curt-labs/API/helpers/apicontext"
	"github.com/curt-labs/API/helpers/database"
)

var (
	GetYearsStmt = `select distinct year, '' from Vehicle v
					join Year as y on v.yearID = y.yearID
					join VehiclePart vp on v.vehicleID = vp.vehicleID
					join Part as p on vp.partID = p.partID
					where p.status in (700, 800, 810, 815, 850, 870, 888, 900, 910, 950)
					order by year desc`
	GetMakesStmt = `select distinct ma.make, group_concat(p.partID) from Make as ma
						join Vehicle as v on ma.makeID = v.makeID
						join Year as y on v.yearID = y.yearID
						join VehiclePart as vp on v.vehicleID = vp.vehicleID
						join Part as p on vp.partID = p.partID
						where y.year = ? && p.status in (700, 800, 810, 815, 850, 870, 888, 900, 910, 950) && p.classID > ?
						group by ma.make
						order by ma.make`
	GetModelsStmt = `select distinct mo.model, group_concat(p.partID) from Model as mo
						join Vehicle as v on mo.modelID = v.modelID
						join Year as y on v.yearID = y.yearID
						join Make as ma on v.makeID = ma.makeID
						join VehiclePart as vp on v.vehicleID = vp.vehicleID
						join Part as p on vp.partID = p.partID
						where y.year = ? && ma.make = ? && p.status in (700, 800, 810, 815, 850, 870, 888, 900, 910, 950) && p.classID > ?
						group by mo.model
						order by mo.model`
	GetStylesStmt = `select distinct s.style, group_concat(p.partID) from Style as s
						join Vehicle as v on s.styleID = v.styleID
						join Year as y on v.yearID = y.yearID
						join Make as ma on v.makeID = ma.makeID
						join Model as mo on v.modelID = mo.modelID
						join VehiclePart as vp on v.vehicleID = vp.vehicleID
						join Part as p on vp.partID = p.partID
						where y.year = ? && ma.make = ? && mo.model = ? && p.status in (700, 800, 810, 815, 850, 870, 888, 900, 910, 950) && p.classID > ?
						group by s.style
						order by s.style`
	GetPartNumbersStmt = `select distinct p.partID from Style as s
							join Vehicle as v on s.styleID = v.styleID
							join Year as y on v.yearID = y.yearID
							join Make as ma on v.makeID = ma.makeID
							join Model as mo on v.modelID = mo.modelID
							join VehiclePart as vp on v.vehicleID = vp.vehicleID
							join Part as p on vp.partID = p.partID
							where y.year = ? && ma.make = ? && mo.model = ? && s.style = ? && p.status in (700, 800, 810, 815, 850, 870, 888, 900, 910, 950) && p.classID > ?
							order by p.partID`
	GetPartNumbersWithoutStyleStmt = `select distinct p.partID from Style as s
							join Vehicle as v on s.styleID = v.styleID
							join Year as y on v.yearID = y.yearID
							join Make as ma on v.makeID = ma.makeID
							join Model as mo on v.modelID = mo.modelID
							join VehiclePart as vp on v.vehicleID = vp.vehicleID
							join Part as p on vp.partID = p.partID
							where y.year = ? && ma.make = ? && mo.model = ? && p.status in (700, 800, 810, 815, 850, 870, 888, 900, 910, 950) && p.classID > ?
							order by p.partID`

	statuses = []int{700, 800, 810, 815, 850, 870, 888, 900, 910, 950}
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

	return nil
}

func (c *CurtLookup) GetMakes(heavyduty bool) error {

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

	return nil
}

func (c *CurtLookup) GetModels(heavyduty bool) error {

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
				"style": bson.RegEx{
					Pattern: "^" + c.Style + "$",
					Options: "i",
				},
			},
		},
		"vehicle_applications.0": bson.M{
			"$exists": true,
		},
		"brand.id": 1,
	}

	err = col.Find(qry).Select(bson.M{
		"vehicle_applications": 1,
		"id":  1,
		"_id": -1,
	}).All(&c.Parts)

	return err
}
