package apicontext

import (
	"github.com/curt-labs/GoAPI/models/customer"
)

type DataContext struct {
	BrandID      int
	WebsiteID    int
	APIKey       string
	CustomerUser *customer.CustomerUser
}
