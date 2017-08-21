package acesFile

type Aces struct {
	Version string `xml:"version,attr"`
	Apps    []app  `xml:"App"`
}

type app struct {
	Action       string          `xml:"action,attr"`
	ID           string          `xml:"id,attr"`
	BaseVehicle  *baseVehicleXml `xml:"BaseVehicle,omitempty"`
	SubModel     *submodelXml    `xml:"SubModel,omitempty"`
	Notes        []string        `xml:"Note,omitempty"`
	Qty          int             `xml:"Qty,omitempty"`
	PartType     *partType       `xml:"PartType,omitempty"`
	MfrLabel     string          `xml:"MfrLabel,omitempty"`
	Position     *position       `xml:"Position,omitempty"`
	Part         string          `xml:"Part,omitempty"`
	Years        *yearsXml       `xml:"Years,omitempty"`
	Make         *makeXml        `xml:"Make,omitempty"`
	Model        *modelXml       `xml:"Model,omitempty"`
	BodyType     *bodyType       `xml:"BodyType,omitempty"`
	BodyNumDoors *bodyNumDoors   `xml:"BodyNumDoors,omitempty"`
	BedLength    []bedLength     `xml:"BedLength,omitempty"`
}

type yearsXml struct {
	To   string `xml:"to,attr"`
	From string `xml:"from,attr"`
}

type bedLength struct {
	ID string `xml:"id,attr,omitempty"`
}

type bodyNumDoors struct {
	ID string `xml:"id,attr,omitempty"`
}

type bodyType struct {
	ID string `xml:"id,attr,omitempty"`
}

type modelXml struct {
	ID string `xml:"id,attr,omitempty"`
}

type makeXml struct {
	ID string `xml:"id,attr,omitempty"`
}

type baseVehicleXml struct {
	ID string `xml:"id,attr,omitempty"`
}

type submodelXml struct {
	ID string `xml:"id,attr,omitempty"`
}

type partType struct {
	ID string `xml:"id,attr,omitempty"`
}

type position struct {
	ID string `xml:"id,attr,omitempty"`
}
