package products

type Package struct {
	ID                 int     `json:"id,omitempty" xml:"id,omitempty"`
	PartID             int     `json:"partId,omitempty" xml:"partId,omitempty"`
	Height             float64 `bson:"height" json:"height" xml:"height"`
	Width              float64 `bson:"width" json:"width" xml:"width"`
	Length             float64 `bson:"length" json:"length" xml:"length"`
	Weight             float64 `bson:"weight" json:"weight" xml:"weight"`
	DimensionUnit      string  `bson:"dimensionUnit" json:"dimensionUnit" xml:"dimensionUnit"`
	DimensionUnitLabel string  `bson:"dimensionUnitLabel"json:"dimensionUnitLabel" xml:"dimensionUnitLabel"`
	WeightUnit         string  `bson:"weightUnit" json:"weightUnit" xml:"weightUnit"`
	WeightUnitLabel    string  `bson:"weightUnitLabel" json:"weightUnitLabel" xml:"weightUnitLabel"`
	PackageUnit        string  `bson:"packageUnit" json:"packageUnit" xml:"packageUnit"`
	PackageUnitLabel   string  `bson:"packageUnitLabel" json:"packageUnitLabel" xml:"packageUnitLabel"`
	Quantity           int     `bson:"quantity" json:"quantity" xml:"quantity"`
	PackageType        string  `bson:"name" json:"name" xml:"name"`
	ParcelAllowed      bool    `bson:"parcelAllowed" json:"parcelAllowed" xml:"parcelAllowed"`
}
