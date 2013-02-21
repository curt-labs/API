package models

type Package struct {
	Height, Width, Length, Quantity   int
	Weight                            float64
	DimensionUnit, DimensionUnitLabel string
	WeightUnit, WeightUnitLabel       string
	PackageUnti, PackageUnitLabel     string
}
