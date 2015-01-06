package apicontextmock

import (
	"github.com/curt-labs/GoAPI/helpers/apicontext"
	"github.com/curt-labs/GoAPI/helpers/database"
	"github.com/curt-labs/GoAPI/models/apiKeyType"
	"github.com/curt-labs/GoAPI/models/customer"
	"github.com/curt-labs/GoAPI/models/site"

	// "database/sql"
	"strings"
)

func Mock() (*apicontext.DataContext, error) {
	var dtx apicontext.DataContext
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
		return &dtx, err
	}

	dtx.WebsiteID = w.ID
	dtx.BrandID = 1
	dtx.APIKey = apiKey
	dtx.CustomerID = c.Id
	dtx.UserID = cu.Id
	return &dtx, err
}

func DeMock(dtx *apicontext.DataContext) error {
	var err error
	var cust customer.Customer
	var user customer.CustomerUser

	cust.Id = dtx.CustomerID
	user.Id = dtx.UserID

	var w site.Website

	w.ID = dtx.WebsiteID
	w.Get()

	err = w.Delete()

	err = cust.Delete()

	err = user.Delete()

	var pub, pri, auth apiKeyType.ApiKeyType
	if database.EmptyDb != nil {
		for _, key := range user.Keys {
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
