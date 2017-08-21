package acesFile

type Aces struct {
	Version string `xml:"version,attr"`
	Apps    []app  `xml:"App"`
}

type app struct {
	Action       string         `xml:"action,attr"`
	ID           string         `xml:"id,attr"`
	BaseVehicle  baseVehicleXml `xml:"BaseVehicle"`
	SubModel     submodelXml    `xml:"SubModel"`
	Notes        []string       `xml:"Note"`
	Qty          int            `xml:"Qty"`
	PartType     partType       `xml:"PartType"`
	MfrLabel     string         `xml:"MfrLabel"`
	Position     position       `xml:"Position"`
	Part         string         `xml:"Part"`
	Years        yearsXml       `xml:"Years"`
	Make         makeXml        `xml:"Make"`
	Model        modelXml       `xml:"Model"`
	BodyType     bodyType       `xml:"BodyType"`
	BodyNumDoors bodyNumDoors   `xml:"BodyNumDoors"`
	BedLength    []bedLength    `xml:"BedLength"`
}

type yearsXml struct {
	To   string `xml:"to,attr"`
	From string `xml:"from,attr"`
}

type bedLength struct {
	ID string `xml:"id,attr"`
}

type bodyNumDoors struct {
	ID string `xml:"id,attr"`
}

type bodyType struct {
	ID string `xml:"id,attr"`
}

type modelXml struct {
	ID string `xml:"id,attr"`
}

type makeXml struct {
	ID string `xml:"id,attr"`
}

type baseVehicleXml struct {
	ID string `xml:"id,attr"`
}

type submodelXml struct {
	ID string `xml:"id,attr"`
}

type partType struct {
	ID string `xml:"id,attr"`
}

type position struct {
	ID string `xml:"id,attr"`
}
