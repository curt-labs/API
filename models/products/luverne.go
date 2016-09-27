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
	Base CategoryVehicleBase `json:"base" xml:"BaseVehicle"`

	Years      []string                `json:"availableYears,omitempty" xml:"availableYears,omitempty"`
	Makes      []string                `json:"availableMakes,omitempty" xml:"availableMakes,omitempty"`
	Models     []string                `json:"availableModels,omitempty" xml:"availableModels,omitempty"`
	Categories []LuverneLookupCategory `json:"lookup_category" xml:"StyleOptions"`
	Products   []Part                  `json:"products" xml:"products"`
}

// LuverneLookupCategory Represents a specific category of `StyleOption` fitments.
type LuverneLookupCategory struct {
	Category   Category                `json:"category" xml:"category"`
	Bodies     []string                `bson:"availableBodies" json:"availableBodies" xml:"availableBodies"`
	Boxes      []string                `bson:"availableBoxes" json:"availableBoxes" xml:"availableBoxes"`
	Cabs       []string                `bson:"availableCabs" json:"availableCabs" xml:"availableCabs"`
	FuelTypes  []string                `bson:"availableFuelTypes" json:"availableFuelTypes" xml:"favailableFuelTypes"`
	WheelTypes []string                `bson:"availableWheelTypes" json:"availableWheelTypes" xml:"availableWheelTypes"`
	Fitments   []LuverneFitment        `bson:"fitments" json:"fitments" xml:"fitments"`
	Products   []LuverneFitmentMapping `bson:"products" json:"products" xml:"products"`
}

type LuverneFitment struct {
	Body      string   `bson:"body" json:"body" xml:"body"`
	BoxLength string   `bson:"box" json:"box" xml:"box"`
	CabLength string   `bson:"cab" json:"cab" xml:"cab"`
	Fuel      string   `bson:"fuel" json:"fuel" xml:"fuel"`
	Wheel     string   `bson:"wheel" json:"wheel" xml:"wheel"`
	SKUs      []string `bson:"skus" json:"skus" xml:"skus"`
}

// // LuverneStyleOption Matches a slice of `Part` that have equal fitments to the matched
// // `Style`.
// type LuverneStyleOption struct {
// 	Style          string           `json:"style" xml:"Style"`
// 	FitmentNumbers []FitmentMapping `json:"fitments" xml:"fitments"`
// }
//
// // FitmentMapping defines the product and any associated attributes that
// // relevant to the fitment on a specific application (install time, drilling, etc)
type LuverneFitmentMapping struct {
	Attributes []LuverneFitmentAttribute `json:"attributes" xml:"Attributes"`
	Number     string                    `json:"product_identifier" xml:"product_identifier"`
}

//
// // FitmentAttribute A name value for a note of a fitment application.
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
	var cleanedParts []Part
	var cats []LuverneLookupCategory
	// cleanedParts, cats := generateCategoryStyles(parts, year, vehicleMake, model)
	// sort.Sort(ByCategoryTitle(cats))
	return cleanedParts, cats, nil
}

//
// func generateCategoryStyles(parts []Part, year, vehicleMake, model string) ([]Part, []LookupCategory) {
// 	lc := make(map[string]LookupCategory, 0)
// 	y := year
// 	ma := strings.ToLower(vehicleMake)
// 	mod := strings.ToLower(model)
// 	var cleanParts []Part
// 	for _, p := range parts {
// 		if len(p.Categories) == 0 {
// 			continue
// 		}
//
// 		for _, va := range p.LuverneVehicles {
// 			if va.Year != y || strings.ToLower(va.Make) != ma || strings.ToLower(va.Model) != mod {
// 				continue
// 			}
//
// 			lc = mapPartToCategoryFitments(p, lc, va.Body, va.BoxLength, va.CabLength, va.FuelType, va.WheelType)
// 		}
//
// 		p.Categories = nil
//
// 		cleanParts = append(cleanParts, p)
// 	}
//
// 	var cats []LookupCategory
// 	for _, l := range lc {
// 		cats = append(cats, l)
// 	}
//
// 	return cleanParts, cats
// }
//
// // AddPart Creates a record of the provided part under the referenced style.
// func (lc *LuverneLookupCategory) AddPart(body, box, cab, fuel, wheel string, p Part) {
// 	if strings.TrimSpace(body) == "" && strings.TrimSpace(box) == "" && strings.TrimSpace(cab) == "" && strings.TrimSpace(fuel) == "" && strings.TrimSpace(wheel) == "" {
// 		style = AllPlaceholder
// 	}
//
// 	for i, options := range lc.StyleOptions {
// 		if strings.TrimSpace(options.Style) == "" {
// 			options.Style = AllPlaceholder
// 			lc.StyleOptions[i].Style = AllPlaceholder
// 		}
// 		if strings.Compare(
// 			strings.ToLower(options.Style),
// 			strings.ToLower(style),
// 		) == 0 {
// 			lc.StyleOptions[i].FitmentNumbers = append(lc.StyleOptions[i].FitmentNumbers,
// 				FitmentMapping{
// 					Number:     p.PartNumber,
// 					Attributes: []FitmentAttribute{},
// 				},
// 			)
// 			return
// 		}
// 	}
//
// 	lc.StyleOptions = append(lc.StyleOptions, StyleOption{
// 		Style: style,
// 		FitmentNumbers: []FitmentMapping{
// 			FitmentMapping{
// 				Number:     p.PartNumber,
// 				Attributes: []FitmentAttribute{},
// 			},
// 		},
// 	})
// }

func mapPartToCategoryFitments(p Part, lookupCats map[string]LuverneLookupCategory, body, box, cab, fuel, wheel string) map[string]LuverneLookupCategory {
	childCat, err := getChildCategory(p.Categories)
	if err != nil || childCat.Identifier.String() == "" {
		return lookupCats
	}

	lc, ok := lookupCats[childCat.Identifier.String()]
	if !ok {
		childCat.PartIDs = nil
		childCat.Children = nil
		childCat.ProductListing = nil
		lc = LuverneLookupCategory{
			Category: childCat,
		}
	}

	// we're going to clear out the category information here, since
	// the products are already being grouped into their respective
	// categories at the higher level. (saves on da bits)
	p.Categories = nil

	// add the part to the appropriate style
	// lc.AddPart(body, box, cab, fuel, wheel, p)

	lookupCats[childCat.Identifier.String()] = lc

	return lookupCats
}
