package apicontext

type DataContext struct {
	BrandID     int
	WebsiteID   int
	APIKey      string
	CustomerID  int
	UserID      string
	Globals     map[string]interface{}
	BrandArray  []int
	BrandString string
}

var AllBrandsArray = []int{1, 3, 4, 5, 6}

// var (
// 	apiToBrandStmt = `select brandID from ApiKeyToBrand as aktb
// 		join ApiKey as ak on ak.id = aktb.keyID
// 		where ak.api_key = ?`
// )

// func (dtx *DataContext) GetBrandsFromKey() ([]int, error) {
// 	var err error
// 	var b int
// 	var brands []int
// 	err = database.Init()
// 	if err != nil {
// 		return brands, err
// 	}
//
// 	stmt, err := database.DB.Prepare(apiToBrandStmt)
// 	if err != nil {
// 		return brands, err
// 	}
// 	defer stmt.Close()
// 	res, err := stmt.Query(dtx.APIKey)
// 	if err != nil {
// 		return brands, err
// 	}
// 	for res.Next() {
// 		err = res.Scan(&b)
// 		if err != nil {
// 			return brands, err
// 		}
// 		brands = append(brands, b)
// 	}
// 	return brands, err
// }

// func (dtx *DataContext) GetBrandsArrayAndString(apiKey string, brandId int) error {
// 	var err error
// 	var brandInts []int
// 	var brandStringArray []string
// 	var brandIdApproved bool = false
//
// 	//get brandIds from apiKey
// 	err = database.Init()
// 	if err != nil {
// 		return err
// 	}
//
// 	stmt, err := database.DB.Prepare(apiToBrandStmt)
// 	if err != nil {
// 		return err
// 	}
// 	defer stmt.Close()
// 	res, err := stmt.Query(apiKey)
// 	if err != nil {
// 		return err
// 	}
// 	var b int
// 	for res.Next() {
// 		err = res.Scan(&b)
// 		if err != nil {
// 			return err
// 		}
// 		brandInts = append(brandInts, b)
// 		brandStringArray = append(brandStringArray, strconv.Itoa(b))
// 	}
// 	if brandId > 0 {
// 		for _, bId := range brandInts {
// 			if bId == brandId {
// 				brandIdApproved = true
// 				dtx.BrandArray = []int{brandId}
// 				dtx.BrandString = "brands:" + strconv.Itoa(brandId)
// 				return err
// 			}
// 		}
// 	}
// 	if brandId > 0 && brandIdApproved == false {
// 		dtx.BrandArray = []int{}
// 		dtx.BrandString = ""
// 		err = errors.New("That brand is not associated with this API Key.")
// 		return err
// 	}
// 	dtx.BrandString = "brands:" + strings.Join(brandStringArray, ",")
// 	dtx.BrandArray = brandInts
// 	return err
// }
