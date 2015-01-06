package apicontext

import (
	"github.com/curt-labs/GoAPI/helpers/database"
	"github.com/curt-labs/GoAPI/models/apiKeyType"
	"github.com/curt-labs/GoAPI/models/customer"
	"github.com/curt-labs/GoAPI/models/site"

	"database/sql"
	"strings"
)

type DataContext struct {
	BrandID      int
	WebsiteID    int
	APIKey       string
	Customer     *customer.Customer
	CustomerUser *customer.CustomerUser
	Globals      map[string]interface{}
}

var (
	apiToBrandStmt = `select brandID from ApiKeyToBrand as aktb 
		join ApiKey as ak on ak.id = aktb.keyID
		where ak.api_key = ?`
)

func (dtx *DataContext) Mock() error {
	var c customer.Customer
	var cu customer.CustomerUser
	var w site.Website
	c.Name = "test cust"
	var pub, pri, auth apiKeyType.ApiKeyType
	if database.EmptyDb != nil {
		//setup apiKeyTypes
		pub.Type = "Public"
		pri.Type = "Private"
		auth.Type = "Authentication"
		pub.Create()
		pri.Create()
		auth.Create()
	}
	c.Create()
	cu.CustomerID = c.Id
	cu.Name = "test cust content user"
	cu.Email = "pretend@test.com"
	cu.Password = "test"
	cu.Sudo = true
	cu.Create()
	var err error
	var apiKey string
	for _, key := range cu.Keys {
		if strings.ToLower(key.Type) == "public" {
			apiKey = key.Key
		}
	}
	err = w.Create()
	w.BrandIDs = append(w.BrandIDs, 1)
	err = w.JoinToBrand()
	if err != nil {
		return err
	}

	dtx.WebsiteID = w.ID
	dtx.BrandID = 1
	dtx.APIKey = apiKey
	dtx.Customer = &c
	dtx.CustomerUser = &cu
	return err
}

func (dtx *DataContext) DeMock() error {
	var err error

	var w site.Website

	w.ID = dtx.WebsiteID
	w.Get()

	err = w.Delete()

	err = dtx.Customer.Delete()

	err = dtx.CustomerUser.Delete()

	var pub, pri, auth apiKeyType.ApiKeyType
	if database.EmptyDb != nil {
		for _, key := range dtx.CustomerUser.Keys {
			if strings.ToLower(key.Type) == "public" {
				pub.ID = key.TypeId
			}
			if strings.ToLower(key.Type) == "private" {
				pri.ID = key.TypeId
			}
			if strings.ToLower(key.Type) == "authentication" {
				auth.ID = key.TypeId
			}
		}
		err = pub.Delete()
		err = pri.Delete()
		err = auth.Delete()

	}

	return err
}

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
