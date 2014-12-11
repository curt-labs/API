package testHelper

import (
	"github.com/curt-labs/GoAPI/helpers/apicontext"
	"github.com/curt-labs/GoAPI/models/apiKeyType"
	"github.com/curt-labs/GoAPI/models/customer"
)

type DataContextDELETETHIS struct {
	BrandID      int
	WebsiteID    int
	APIKey       string
	CustomerUser *customer.CustomerUser
	Globals      map[string]interface{}
}
