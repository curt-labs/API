package products

import (
	"encoding/json"
	"fmt"
	"sort"
	"strings"

	"github.com/curt-labs/API/helpers/database"
	"github.com/curt-labs/API/helpers/redis"

	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

// LuverneLookupContext Holds required configuration settings and resources.
type LuverneLookupContext struct {
	Session  *mgo.Session
	Statuses []int
}

// LuverneVehicleBase The lowest level vehicle properties used to query
// a LuverneCategoryVehicle.
type LuverneVehicleBase struct {
	Year  string `json:"year" xml:"year,attr"`
	Make  string `json:"make" xml:"make,attr"`
	Model string `json:"model" xml:"model,attr"`
}

// LuverneCategoryVehicle Represents the requested `Base` vehicle and all matching
// LookupCategory types that are fitments of the `Base`.
type LuverneCategoryVehicle struct {
	Base       CategoryVehicleBase     `json:"base" xml:"BaseVehicle"`
	Years      []string                `json:"availableYears,omitempty" xml:"availableYears,omitempty"`
	Makes      []string                `json:"availableMakes,omitempty" xml:"availableMakes,omitempty"`
	Models     []string                `json:"availableModels,omitempty" xml:"availableModels,omitempty"`
	Categories []LuverneLookupCategory `json:"lookup_category" xml:"StyleOptions"`
	Products   []Part                  `json:"products" xml:"products"`
}

// LuverneLookupCategory Represents a specific category of `StyleOption` fitments.
type LuverneLookupCategory struct {
	Category Category                `json:"category" xml:"category"`
	Fitments []*LuverneFitment       `bson:"fitments" json:"fitments" xml:"fitments"`
	Products []LuverneFitmentMapping `bson:"products" json:"products" xml:"products"`
}

type LuverneFitment struct {
	Title   string
	Options []string
}

type ByLuverneCategoryTitle []LuverneLookupCategory

// sort functions for ByLuverneCategoryTitle
func (a ByLuverneCategoryTitle) Len() int           { return len(a) }
func (a ByLuverneCategoryTitle) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a ByLuverneCategoryTitle) Less(i, j int) bool { return a[i].Category.Title < a[j].Category.Title }

// FitmentMapping defines the product and any associated attributes that
// relevant to the fitment on a specific application (install time, drilling, etc)
type LuverneFitmentMapping struct {
	Attributes []LuverneFitmentAttribute `json:"fitment_attributes" xml:"fitment_attributes"`
	Number     string                    `json:"product_identifier" xml:"product_identifier"`
}

// FitmentAttribute A name value for a note of a fitment application.
type LuverneFitmentAttribute struct {
	Key   string `json:"key" xml:"Key"`
	Value string `json:"value" xml:"Value"`
}

// Query Returns a `CategoryVehicle` that holds matching information for the
// queried `CategoryVehicleBase` attributes.
func LuverneQuery(ctx *LuverneLookupContext, args ...string) (*LuverneCategoryVehicle, error) {

	var redisKey string
	var category string
	var vehicle LuverneCategoryVehicle

	for i, arg := range args {
		if i == 0 {
			redisKey = fmt.Sprintf("luverne:%s", arg)
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

	switch len(args) {
	case 1:
		vehicle.Base.Year = args[0]
	case 2:
		vehicle.Base.Year = args[0]
		vehicle.Base.Make = args[1]
	case 3:
		vehicle.Base.Year = args[0]
		vehicle.Base.Make = args[1]
		vehicle.Base.Model = args[2]
	case 4:
		vehicle.Base.Year = args[0]
		vehicle.Base.Make = args[1]
		vehicle.Base.Model = args[2]
		category = args[3]
	}

	if vehicle.Base.Year == "" {
		vehicle.Years, err = getLuverneYears(ctx)
	} else if vehicle.Base.Year != "" && vehicle.Base.Make == "" {
		vehicle.Makes, err = getLuverneMakes(ctx, vehicle.Base.Year)
	} else if vehicle.Base.Year != "" && vehicle.Base.Make != "" && vehicle.Base.Model == "" {
		vehicle.Models, err = getLuverneModels(ctx, vehicle.Base.Year, vehicle.Base.Make)
	} else if vehicle.Base.Year != "" && vehicle.Base.Make != "" && vehicle.Base.Model != "" {
		vehicle.Products, vehicle.Categories, err = getLuverneStyles(ctx, vehicle.Base.Year, vehicle.Base.Make, vehicle.Base.Model, category)
	}

	redis.Setex(redisKey, vehicle, 60*60*24)

	return &vehicle, err
}

func getLuverneYears(ctx *LuverneLookupContext) ([]string, error) {
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
		"brand.id": 4,
	}

	var res []string
	err := c.Find(qry).Select(bson.M{
		"luverne_applications.year": 1,
		"_id": -1,
	}).Distinct("luverne_applications.year", &res)

	if err != nil {
		return nil, err
	}

	sort.Sort(sort.Reverse(sort.StringSlice(res)))

	return res, err
}

func getLuverneMakes(ctx *LuverneLookupContext, year string) ([]string, error) {
	if ctx == nil {
		return nil, fmt.Errorf("missing context")
	} else if ctx.Session == nil {
		return nil, fmt.Errorf("invalid mongodb connection")
	}

	c := ctx.Session.DB(database.ProductMongoDatabase).C(database.ProductCollectionName)

	type Apps struct {
		Apps []VehicleApplication `bson:"luverne_applications"`
	}

	var apps []Apps
	qry := bson.M{
		"luverne_applications": bson.M{
			"$elemMatch": bson.M{
				"year": year,
			},
		},
		"status": bson.M{
			"$in": ctx.Statuses,
		},
		"brand.id": 4,
	}
	err := c.Find(qry).Select(bson.M{"luverne_applications.make": 1, "_id": 0}).All(&apps)
	if err != nil {
		return nil, err
	}

	var makes []string

	existing := make(map[string]string, 0)
	for _, app := range apps {
		for _, a := range app.Apps {
			a.Make = strings.Title(a.Make)
			if _, ok := existing[a.Make]; !ok {
				makes = append(makes, a.Make)
				existing[a.Make] = a.Make
			}
		}
	}

	sort.Strings(makes)

	return makes, err
}

func getLuverneModels(ctx *LuverneLookupContext, year, vehicleMake string) ([]string, error) {
	if ctx == nil {
		return nil, fmt.Errorf("missing context")
	} else if ctx.Session == nil {
		return nil, fmt.Errorf("invalid mongodb connection")
	}

	c := ctx.Session.DB(database.ProductMongoDatabase).C(database.ProductCollectionName)

	type Apps struct {
		Apps []VehicleApplication `bson:"luverne_applications"`
	}

	var apps []Apps
	err := c.Find(bson.M{
		"luverne_applications": bson.M{
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
		"brand.id": 4,
	}).Select(bson.M{"luverne_applications.$.model": 1, "_id": 0}).All(&apps)
	if err != nil {
		return nil, err
	}

	var models []string

	existing := make(map[string]string, 0)
	for _, app := range apps {
		for _, a := range app.Apps {
			a.Model = strings.Title(a.Model)
			if _, ok := existing[a.Model]; !ok {
				models = append(models, a.Model)
				existing[a.Model] = a.Model
			}
		}
	}

	sort.Strings(models)

	return models, err
}

func getLuverneStyles(ctx *LuverneLookupContext, year, vehicleMake, model, category string) ([]Part, []LuverneLookupCategory, error) {
	if ctx == nil {
		return nil, nil, fmt.Errorf("missing context")
	} else if ctx.Session == nil {
		return nil, nil, fmt.Errorf("invalid mongodb connection")
	}

	c := ctx.Session.DB(database.ProductMongoDatabase).C(database.ProductCollectionName)

	type Apps struct {
		Apps    []VehicleApplication `bson:"vehicle_applications"`
		PartNum string               `bson:"part_number"`
	}

	var parts []Part
	qry := bson.M{
		"luverne_applications": bson.M{
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
			"$in": []int{800, 900},
		},
		"brand.id": 4,
	}

	if category != "" {
		qry["categories.title"] = category
	}
	err := c.Find(qry).All(&parts)
	if err != nil || len(parts) == 0 {
		return nil, nil, err
	}
	cleanedParts, cats := generateLuverneCategoryStyles(parts, year, vehicleMake, model)
	sort.Sort(ByLuverneCategoryTitle(cats))
	return cleanedParts, cats, nil
}

//
func generateLuverneCategoryStyles(parts []Part, year, vehicleMake, model string) ([]Part, []LuverneLookupCategory) {
	lc := make(map[string]LuverneLookupCategory, 0)
	y := year
	ma := strings.ToLower(vehicleMake)
	mod := strings.ToLower(model)
	var cleanParts []Part
	for _, p := range parts {
		if len(p.Categories) == 0 {
			continue
		}

		for _, va := range p.LuverneVehicles {
			if va.Year != y || strings.ToLower(va.Make) != ma || strings.ToLower(va.Model) != mod {
				continue
			}

			lc = mapPartToCategoryFitments(p, lc, va.Body, va.BoxLength, va.CabLength, va.FuelType, va.WheelType)
		}

		p.Categories = nil

		cleanParts = append(cleanParts, p)
	}

	var cats []LuverneLookupCategory
	for _, l := range lc {
		cats = append(cats, l)
	}

	return cleanParts, cats
}

func mapPartToCategoryFitments(p Part, lookupCats map[string]LuverneLookupCategory, body, box, cab, fuel, wheel string) map[string]LuverneLookupCategory {
	for _, cat := range p.Categories {
		lc, ok := lookupCats[cat.Identifier.String()]
		if !ok {
			cat.PartIDs = nil
			cat.Children = nil
			cat.ProductListing = nil
			lc = LuverneLookupCategory{
				Category: cat,
			}
		}

		lc.ProcessPart(body, box, cab, fuel, wheel, p)
		lc.GenerateFitmentsOfPart(body, box, cab, fuel, wheel, p)
		lookupCats[cat.Identifier.String()] = lc
	}
	return lookupCats
}

// Process Part creates the individual record for the part for the lc.Products and lc.Fitments
func (lc *LuverneLookupCategory) ProcessPart(body, box, cab, fuel, wheel string, p Part) {
	// map the Part and fitments to the category to create a LuverneLookUpCategory
	var newP LuverneFitmentMapping
	newP.Number = p.PartNumber
	if body != "" {
		newP.Attributes = append(newP.Attributes, LuverneFitmentAttribute{"Body", body})
	}
	if box != "" {
		newP.Attributes = append(newP.Attributes, LuverneFitmentAttribute{"Box", box})
	}
	if cab != "" {
		newP.Attributes = append(newP.Attributes, LuverneFitmentAttribute{"Cab", cab})
	}
	if fuel != "" {
		newP.Attributes = append(newP.Attributes, LuverneFitmentAttribute{"Fuel", fuel})
	}
	if wheel != "" {
		newP.Attributes = append(newP.Attributes, LuverneFitmentAttribute{"Wheel", wheel})
	}
	lc.Products = append(lc.Products, newP)
}

func (lc *LuverneLookupCategory) GenerateFitmentsOfPart(body, box, cab, fuel, wheel string, p Part) {
	// if fitments are empty (first time trying to generate fitments)
	if len(lc.Fitments) == 0 {
		if body != "" {
			newFitment := &LuverneFitment{Title: "Body"}
			newFitment.Options = append(newFitment.Options, body)
			lc.Fitments = append(lc.Fitments, newFitment)
		}
		if box != "" {
			newFitment := &LuverneFitment{Title: "Box"}
			newFitment.Options = append(newFitment.Options, box)
			lc.Fitments = append(lc.Fitments, newFitment)
		}
		if cab != "" {
			newFitment := &LuverneFitment{Title: "Cab"}
			newFitment.Options = append(newFitment.Options, cab)
			lc.Fitments = append(lc.Fitments, newFitment)
		}
		if fuel != "" {
			newFitment := &LuverneFitment{Title: "Fuel"}
			newFitment.Options = append(newFitment.Options, fuel)
			lc.Fitments = append(lc.Fitments, newFitment)
		}
		if wheel != "" {
			newFitment := &LuverneFitment{Title: "Wheel"}
			newFitment.Options = append(newFitment.Options, wheel)
			lc.Fitments = append(lc.Fitments, newFitment)
		}
	} else { // beginging fitments have been generated, now add additional non-duplicate fitment options
		for _, fit := range lc.Fitments {
			// BODY - Check for duplicates, if not a duplicate, add it to the fitment options
			if fit.Title == "Body" && body != "" && !CheckDuplicateOptions(fit.Options, body) {
				fit.Options = append(fit.Options, body)
			}
			// Box
			if fit.Title == "Box" && box != "" && !CheckDuplicateOptions(fit.Options, box) {
				fit.Options = append(fit.Options, box)
			}
			// Cab
			if fit.Title == "Cab" && cab != "" && !CheckDuplicateOptions(fit.Options, cab) {
				fit.Options = append(fit.Options, cab)
			}
			// Fuel
			if fit.Title == "Fuel" && fuel != "" && !CheckDuplicateOptions(fit.Options, fuel) {
				fit.Options = append(fit.Options, fuel)
			}
			// Wheel
			if fit.Title == "Wheel" && wheel != "" && !CheckDuplicateOptions(fit.Options, wheel) {
				fit.Options = append(fit.Options, wheel)
			} // end wheel
		} // end each fitment
	} // end if len(fitments) !== 0
}

func CheckDuplicateOptions(options []string, option string) bool {
	for _, opt := range options {
		if opt == option {
			return true
		}
	}
	return false
}
