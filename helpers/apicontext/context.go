package apicontext

import (
	"database/sql"
	"github.com/curt-labs/GoAPI/helpers/database"
	"strconv"
	"strings"
)

type DataContext struct {
	BrandID    int
	WebsiteID  int
	APIKey     string
	CustomerID int
	UserID     string
	Globals    map[string]interface{}
}

var (
	apiToBrandStmt = `select brandID from ApiKeyToBrand as aktb 
		join ApiKey as ak on ak.id = aktb.keyID
		where ak.api_key = ?`
)

func (dtx *DataContext) GetBrandsFromKey() ([]int, error) {
	var err error
	var b int
	var brands []int
	db, err := sql.Open("mysql", database.ConnectionString())
	if err != nil {
		return brands, err
	}
	defer db.Close()

	stmt, err := db.Prepare(apiToBrandStmt)
	if err != nil {
		return brands, err
	}
	defer stmt.Close()
	res, err := stmt.Query(dtx.APIKey)
	if err != nil {
		return brands, err
	}
	for res.Next() {
		err = res.Scan(&b)
		if err != nil {
			return brands, err
		}
		brands = append(brands, b)
	}
	return brands, err
}

func GetBrandsString(apiKey string, brandId int) (string, error) {
	var err error
	var brands string
	var brandInts []int
	var brandStrs []string
	var brandIdApproved bool = false

	//get brandIds from apiKey
	db, err := sql.Open("mysql", database.ConnectionString())
	if err != nil {
		return brands, err
	}
	defer db.Close()

	stmt, err := db.Prepare(apiToBrandStmt)
	if err != nil {
		return brands, err
	}
	defer stmt.Close()
	res, err := stmt.Query(apiKey)
	if err != nil {
		return brands, err
	}
	var b int
	for res.Next() {
		err = res.Scan(&b)
		if err != nil {
			return brands, err
		}
		brandInts = append(brandInts, b)
	}
	for _, bId := range brandInts {
		if bId == brandId {
			brandIdApproved = true
		}
		brandStrs = append(brandStrs, strconv.Itoa(bId))
	}
	if brandId > 0 && brandIdApproved == true {
		brands = strconv.Itoa(brandId)
		return brands, err
	}
	if brandId > 0 && brandIdApproved == false {
		return "", err
	}
	brands = strings.Join(brandStrs, ",")

	return brands, err
}
