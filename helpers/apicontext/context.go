package apicontext

import (
	"github.com/curt-labs/GoAPI/helpers/database"
	"github.com/curt-labs/GoAPI/models/apiKeyType"
	"github.com/curt-labs/GoAPI/models/customer"
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

func (dtx *DataContext) Mock() error {
	var c customer.Customer
	var cu customer.CustomerUser
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

	dtx.WebsiteID = 1
	dtx.BrandID = 1
	dtx.APIKey = apiKey
	dtx.Customer = &c
	dtx.CustomerUser = &cu
	return err
}

func (dtx *DataContext) DeMock() error {
	var err error
	err = dtx.Customer.Delete()
	err = dtx.CustomerUser.Delete()
	var pub, pri, auth apiKeyType.ApiKeyType
	if database.EmptyDb != nil {
		err = pub.Delete()
		err = pri.Delete()
		err = auth.Delete()
	}
	return err
}
