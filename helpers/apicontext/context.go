package apicontext

import (
	"database/sql"
	"errors"
	"github.com/curt-labs/API/helpers/database"
	"strconv"
	"strings"
)

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

var (
	apiToBrandStmt = `select ID from Brand`
)

// @deprecated - API keys are no longer tied to specific brands.
// This function now returns all brands reguardless of what brands
// A user might have tied a key to.
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

func (dtx *DataContext) GetBrandsArrayAndString(apiKey string, brandId int) error {
	var err error
	var brandInts []int
	var brandStringArray []string
	var brandIdApproved bool = false

	//get brandIds from apiKey
	db, err := sql.Open("mysql", database.ConnectionString())
	if err != nil {
		return err
	}
	defer db.Close()

	stmt, err := db.Prepare(apiToBrandStmt)
	if err != nil {
		return err
	}
	defer stmt.Close()
	res, err := stmt.Query(apiKey)
	if err != nil {
		return err
	}
	var b int
	for res.Next() {
		err = res.Scan(&b)
		if err != nil {
			return err
		}
		brandInts = append(brandInts, b)
		brandStringArray = append(brandStringArray, strconv.Itoa(b))
	}
	if brandId > 0 {
		for _, bId := range brandInts {
			if bId == brandId {
				brandIdApproved = true
				dtx.BrandArray = []int{brandId}
				dtx.BrandString = "brands:" + strconv.Itoa(brandId)
				return err
			}
		}
	}
	if brandId > 0 && brandIdApproved == false {
		dtx.BrandArray = []int{}
		dtx.BrandString = ""
		err = errors.New("That brand is not associated with this API Key.")
		return err
	}
	dtx.BrandString = "brands:" + strings.Join(brandStringArray, ",")
	dtx.BrandArray = brandInts
	return err
}
