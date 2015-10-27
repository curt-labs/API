package products

type Package struct {
	ID                 int         `json:"id,omitempty" xml:"id,omitempty"`
	PartID             int         `json:"partId,omitempty" xml:"partId,omitempty"`
	Height             float64     `json:"height,omitempty" xml:"height,omitempty"`
	Width              float64     `json:"width,omitempty" xml:"width,omitempty"`
	Length             float64     `json:"length,omitempty" xml:"length,omitempty"`
	Weight             float64     `json:"weight,omitempty" xml:"weight,omitempty"`
	DimensionUnit      string      `json:"dimensionUnit,omitempty" xml:"dimensionUnit,omitempty"`
	DimensionUnitLabel string      `json:"dimensionUnitLabel,omitempty" xml:"dimensionUnitLabel,omitempty"`
	WeightUnit         string      `json:"weightUnit,omitempty" xml:"weightUnit,omitempty"`
	WeightUnitLabel    string      `json:"weightUnitLabel,omitempty" xml:"weightUnitLabel,omitempty"`
	PackageUnit        string      `json:"packageUnit,omitempty" xml:"packageUnit,omitempty"`
	PackageUnitLabel   string      `json:"packageUnitLabel,omitempty" xml:"packageUnitLabel,omitempty"`
	Quantity           int         `json:"quantity,omitempty" xml:"quantity,omitempty"`
	PackageType        PackageType `json:"packageType,omitempty" xml:"packageType,omitempty"`
}

type PackageType struct {
	ID   int    `json:"id,omitempty" xml:"id,omitempty"`
	Name string `json:"name,omitempty" xml:"name,omitempty"`
}
