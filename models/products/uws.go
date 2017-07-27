package products

import (
	"encoding/json"
	"fmt"
	"log"
	"sort"
	"strconv"
	"strings"

	"github.com/curt-labs/API/helpers/database"
	"github.com/curt-labs/API/helpers/redis"

	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

// UWSLookupContext Holds required configuration settings and resources.
type UWSLookupContext struct {
	Session  *mgo.Session
	Statuses []int
}

// UWSVehicleBase The lowest level vehicle properties used to query
// a UWSCategoryVehicle.
type UWSVehicleBase struct {
	Year  int    `json:"year" xml:"year,attr"`
	Make  string `json:"make" xml:"make,attr"`
	Model string `json:"model" xml:"model,attr"`
}

// UWSCategoryVehicle Represents the requested `Base` vehicle and all matching
// LookupCategory types that are fitments of the `Base`.
type UWSCategoryVehicle struct {
	Base       CategoryVehicleBaseUWS `json:"base_vehicle" xml:"base_vehicle"`
	Years      []int                  `json:"available_years,omitempty" xml:"available_years,omitempty"`
	Makes      []string               `json:"available_makes,omitempty" xml:"available_makes,omitempty"`
	Models     []string               `json:"available_models,omitempty" xml:"available_models,omitempty"`
	Categories []UWSLookupCategory    `json:"lookup_category" xml:"lookup_category"`
	Products   []Part                 `json:"products" xml:"products"`
}

// UWSLookupCategory Represents a specific category of `StyleOption` fitments.
type UWSLookupCategory struct {
	Category Category            `json:"category" xml:"category"`
	Fitments []*UWSFitment       `bson:"fitments" json:"fitments" xml:"fitments"`
	Products []UWSFitmentMapping `bson:"products" json:"products" xml:"products"`
}

type UWSFitment struct {
	Title   string   `json:"title" xml:"title"`
	Options []string `json:"options" xml:"options"`
}

type ByUWSCategoryTitle []UWSLookupCategory

// sort functions for ByUWSCategoryTitle
func (a ByUWSCategoryTitle) Len() int           { return len(a) }
func (a ByUWSCategoryTitle) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a ByUWSCategoryTitle) Less(i, j int) bool { return a[i].Category.Title < a[j].Category.Title }

// FitmentMapping defines the product and any associated attributes that
// relevant to the fitment on a specific application (install time, drilling, etc)
type UWSFitmentMapping struct {
	Attributes []UWSFitmentAttribute `json:"fitment_attributes" xml:"fitment_attributes"`
	Number     string                `json:"product_identifier" xml:"product_identifier"`
}

// FitmentAttribute A name value for a note of a fitment application.
type UWSFitmentAttribute struct {
	Key   string `json:"key" xml:"key"`
	Value string `json:"value" xml:"value"`
}

// Query Returns a `CategoryVehicle` that holds matching information for the
// queried `CategoryVehicleBase` attributes.
func UWSQuery(ctx *UWSLookupContext, args ...string) (*UWSCategoryVehicle, error) {

	var redisKey string
	var category string
	var vehicle UWSCategoryVehicle

	for i, arg := range args {
		if i == 0 {
			redisKey = fmt.Sprintf("UWS:%s", arg)
		} else {
			redisKey = fmt.Sprintf("%s:%s", redisKey, arg)
		}
	}
	data, err := redis.Get(redisKey)
	if err == nil && len(data) > 0 {
		err = json.Unmarshal(data, &vehicle)
		if err == nil {
			return &vehicle, nil
		}
	}
	log.Println("missed cache:", redisKey)
	switch len(args) {
	case 1:
		vehicle.Base.Year, err = strconv.Atoi(args[0])
		if err != nil {
			return &vehicle, err
		}
	case 2:
		vehicle.Base.Year, err = strconv.Atoi(args[0])
		if err != nil {
			return &vehicle, err
		}
		vehicle.Base.Make = args[1]
	case 3:
		vehicle.Base.Year, err = strconv.Atoi(args[0])
		if err != nil {
			return &vehicle, err
		}
		vehicle.Base.Make = args[1]
		vehicle.Base.Model = args[2]
	case 4:
		vehicle.Base.Year, err = strconv.Atoi(args[0])
		if err != nil {
			return &vehicle, err
		}
		vehicle.Base.Make = args[1]
		vehicle.Base.Model = args[2]
		category = args[3]
	}

	if vehicle.Base.Year == 0 {
		vehicle.Years, err = getUWSYears(ctx)
	} else if vehicle.Base.Year != 0 && vehicle.Base.Make == "" {
		vehicle.Makes, err = getUWSMakes(ctx, vehicle.Base.Year)
	} else if vehicle.Base.Year != 0 && vehicle.Base.Make != "" && vehicle.Base.Model == "" {
		vehicle.Models, err = getUWSModels(ctx, vehicle.Base.Year, vehicle.Base.Make)
	} else if vehicle.Base.Year != 0 && vehicle.Base.Make != "" && vehicle.Base.Model != "" {
		vehicle.Products, vehicle.Categories, err = getUWSStyles(ctx, vehicle.Base.Year, vehicle.Base.Make, vehicle.Base.Model, category)
	}

	redis.Setex(redisKey, vehicle, 60*60*24)

	return &vehicle, err
}

func getUWSYears(ctx *UWSLookupContext) ([]int, error) {
	if ctx == nil {
		return nil, fmt.Errorf("missing context")
	} else if ctx.Session == nil {
		return nil, fmt.Errorf("invalid mongodb connection")
	}

	c := ctx.Session.DB(database.ProductMongoDatabase).C(database.ProductCollectionName)
	qry := bson.M{
		"status": bson.M{
			"$in": ctx.Statuses,
		},
		"brand.id": 6,
	}
	var res []int
	err := c.Find(qry).Select(bson.M{
		"uws_applications.year": 1,
		"_id": -1,
	}).Distinct("uws_applications.year", &res)
	if err != nil {
		return nil, err
	}
	sort.Sort(sort.Reverse(sort.IntSlice(res)))
	// sort.Sort(sort.Reverse(sort.StringSlice(res)))

	return res, err
}

func getUWSMakes(ctx *UWSLookupContext, year int) ([]string, error) {
	if ctx == nil {
		return nil, fmt.Errorf("missing context")
	} else if ctx.Session == nil {
		return nil, fmt.Errorf("invalid mongodb connection")
	}

	c := ctx.Session.DB(database.ProductMongoDatabase).C(database.ProductCollectionName)

	type Apps struct {
		Apps []UWSApplication `bson:"uws_applications"`
	}

	var apps []Apps
	qry := bson.M{
		"uws_applications": bson.M{
			"$elemMatch": bson.M{
				"year": year,
			},
		},
		"status": bson.M{
			"$in": ctx.Statuses,
		},
		"brand.id": 6,
	}

	err := c.Find(qry).Select(bson.M{"uws_applications.make": 1, "uws_applications.year": 1, "_id": 0}).All(&apps)
	if err != nil {
		return nil, err
	}
	var makes []string
	existing := make(map[string]string, 0)
	for _, app := range apps {
		for _, a := range app.Apps {
			a.Make = strings.Title(a.Make)
			if _, ok := existing[a.Make]; !ok {
				if a.Year == year {
					makes = append(makes, a.Make)
					existing[a.Make] = a.Make
				}
			}
		}
	}
	sort.Strings(makes)

	return makes, err
}

func getUWSModels(ctx *UWSLookupContext, year int, vehicleMake string) ([]string, error) {
	if ctx == nil {
		return nil, fmt.Errorf("missing context")
	} else if ctx.Session == nil {
		return nil, fmt.Errorf("invalid mongodb connection")
	}

	c := ctx.Session.DB(database.ProductMongoDatabase).C(database.ProductCollectionName)

	type Apps struct {
		Apps []UWSApplication `bson:"uws_applications"`
	}

	var apps []Apps
	err := c.Find(bson.M{
		"uws_applications": bson.M{
			"$elemMatch": bson.M{
				"year": year,
				"make": bson.RegEx{
					Pattern: "^" + vehicleMake + "$",
					Options: "i",
				},
			},
		},
		"status": bson.M{
			"$in": ctx.Statuses,
		},
		"brand.id": 6,
	}).Select(bson.M{"uws_applications": 1, "_id": 0}).All(&apps)
	if err != nil {
		return nil, err
	}

	var models []string

	existing := make(map[string]string, 0)
	for _, app := range apps {
		for _, a := range app.Apps {
			// Some parts support multi-year and different makes, so we have to filter the year and make back out
			if strings.EqualFold(strconv.Itoa(a.Year), strconv.Itoa(year)) && strings.EqualFold(a.Make, vehicleMake) {
				a.Model = strings.Title(a.Model)
				if _, ok := existing[a.Model]; !ok {
					models = append(models, a.Model)
					existing[a.Model] = a.Model
				}
			}
		}
	}

	sort.Strings(models)

	return models, err
}

func getUWSStyles(ctx *UWSLookupContext, year int, vehicleMake, model, category string) ([]Part, []UWSLookupCategory, error) {
	if ctx == nil {
		return nil, nil, fmt.Errorf("missing context")
	} else if ctx.Session == nil {
		return nil, nil, fmt.Errorf("invalid mongodb connection")
	}

	c := ctx.Session.DB(database.ProductMongoDatabase).C(database.ProductCollectionName)

	type Apps struct {
		Apps    []UWSApplication `bson:"vehicle_applications"`
		PartNum string           `bson:"part_number"`
	}

	var parts []Part
	qry := bson.M{
		"uws_applications": bson.M{
			"$elemMatch": bson.M{
				"year": year,
				"make": bson.RegEx{
					Pattern: "^" + vehicleMake + "$",
					Options: "i",
				},
				"model": bson.RegEx{
					Pattern: "^" + model + "$",
					Options: "i",
				},
			},
		},
		"status": bson.M{
			"$in": []int{700, 800, 810, 815, 850, 870, 888, 900, 910, 950},
		},
		"brand.id": 6,
	}
	if category != "" {
		qry["categories.title"] = category
	}
	err := c.Find(qry).All(&parts)
	if err != nil || len(parts) == 0 {
		return nil, nil, err
	}
	cleanedParts, cats := generateUWSCategoryStyles(parts, year, vehicleMake, model)
	sort.Sort(ByUWSCategoryTitle(cats))
	return cleanedParts, cats, nil
}

//
func generateUWSCategoryStyles(parts []Part, year int, vehicleMake, model string) ([]Part, []UWSLookupCategory) {
	lc := make(map[string]UWSLookupCategory, 0)
	y := year
	ma := strings.ToLower(vehicleMake)
	mod := strings.ToLower(model)

	var cleanParts []Part
	for _, p := range parts {
		if len(p.Categories) == 0 {
			continue
		}

		for _, va := range p.UWSVehicles {
			if va.Year != y || strings.ToLower(va.Make) != ma || strings.ToLower(va.Model) != mod {
				continue
			}

			lc = mapUWSPartToCategoryFitments(p, lc)
		}

		p.Categories = nil

		cleanParts = append(cleanParts, p)
	}

	var cats []UWSLookupCategory
	for _, l := range lc {
		cats = append(cats, l)
	}

	return cleanParts, cats
}

func mapUWSPartToCategoryFitments(p Part, lookupCats map[string]UWSLookupCategory) map[string]UWSLookupCategory {
	for _, cat := range p.Categories {
		lc, ok := lookupCats[cat.Identifier.String()]
		if !ok {
			cat.PartIDs = nil
			cat.Children = nil
			cat.ProductListing = nil
			lc = UWSLookupCategory{
				Category: cat,
			}
		}

		lc.ProcessPart(p)
		lookupCats[cat.Identifier.String()] = lc
	}
	return lookupCats
}

// Process Part creates the individual record for the part for the lc.Products and lc.Fitments
func (lc *UWSLookupCategory) ProcessPart(p Part) {
	// map the Part and fitments to the category to create a UWSLookUpCategory
	var newP UWSFitmentMapping
	newP.Number = p.PartNumber
	lc.Products = append(lc.Products, newP)
}
