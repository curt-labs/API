package products

// AcesVehicle ...
type AcesVehicle struct {
	BaseVehicleID     int                   `bson:"-" json:"-" xml:"-"`
	Base              BaseVehicle           `bson:"base" json:"base" xml:"base"`
	Submodel          string                `bson:"submodel" json:"submodel" xml:"submodel"`
	Attributes        []AcesConfiguration   `bson:"configurations" json:"configurations" xml:"configurations"`
	AttributesByIndex map[int]Configuration `bson:"-" json:"-" xml:"-"`
}

// // BaseVehicle ...
// type BaseVehicle struct {
// 	Year  int    `bson:"year" json:"year" xml:"year"`
// 	Make  string `bson:"make" json:"make" xml:"make"`
// 	Model string `bson:"model" json:"model" xml:"model"`
// }

// Configuration ...
type AcesConfiguration struct {
	Options []ConfigOption `bson:"options" json:"options" xml:"options"`
}

// ConfigOption ...
type ConfigOption struct {
	Key   string `bson:"name" json:"name" xml:"name,attr"`
	Value string `bson:"value" json:"value" xml:"value,attr"`
}
