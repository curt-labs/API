package middleware

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/curt-labs/API/helpers/apicontext"
	"github.com/curt-labs/API/helpers/database"
	"github.com/curt-labs/API/helpers/error"
	"github.com/curt-labs/API/models/cart"
	"github.com/curt-labs/API/models/customer"
	"github.com/go-martini/martini"
	"gopkg.in/mgo.v2/bson"
)

var (
	ExcusedRoutes = []string{"/status", "/customer/auth", "/customer/user", "/new/customer/auth", "/customer/user/register", "/customer/user/resetPassword", "/cartIntegration/priceTypes", "/cartIntegration", "/cache"}

	FindMagentoUserByKey = `SELECT customer_id FROM curtgroup_api_keys WHERE api_key = ?`

	FindInternalMagentoUserByKey = `SELECT customer_id FROM curtgroup_api_keys WHERE api_key = ? AND is_admin = 1`

	GetKeyType = `SELECT akt.type FROM ApiKey as ak, ApiKeyType as akt WHERE akt.id = ak.type_id AND ak.api_key=?`
)

func Meddler() martini.Handler {
	return func(res http.ResponseWriter, r *http.Request, c martini.Context) {
		res.Header().Add("Access-Control-Allow-Origin", "*")
		res.Header().Add("Cache-Control", "max-age=86400")
		if strings.ToLower(r.Method) == "options" {
			return
		}

		if strings.Contains(r.URL.String(), "favicon") {
			res.Write([]byte(""))
			return
		}

		excused := false
		for _, route := range ExcusedRoutes {
			if strings.Contains(r.URL.String(), route) {
				excused = true
			}
		}

		// check if we need to make a call
		// to the shopping cart middleware
		if strings.Contains(strings.ToLower(r.URL.Path), "/shopify/account") { // account perms
			if strings.ToLower(r.URL.Path) == "/shopify/account/login" || (strings.ToLower(r.Method) == "post" && strings.ToLower(r.URL.Path) == "/shopify/account") {
				shopID := r.URL.Query().Get("shop")
				var crt cart.Shop
				if bson.IsObjectIdHex(shopID) {
					crt.Id = bson.ObjectIdHex(shopID)
				}
				c.Map(&crt)
			} else if err := mapCartAccount(c, res, r); err != nil {
				apierror.GenerateError("", err, res, r)
				return
			}
			excused = true
		} else if strings.Contains(strings.ToLower(r.URL.Path), "/shopify") { // shop perms
			if err := mapCart(c, res, r); err != nil {
				apierror.GenerateError("", err, res, r)
				return
			}
			excused = true
		}

		if !excused {
			dataContext, err := processDataContext(r, c)
			if err != nil {
				apierror.GenerateError("Trouble processing the data context", err, res, r, http.StatusUnauthorized)
				return
			}
			c.Map(dataContext)
		}

		c.Next()

		go logRequest(res, r, time.Now())
	}
}

func mapCart(c martini.Context, res http.ResponseWriter, r *http.Request) error {
	qs := r.URL.Query()
	var shopId string
	if qsId := qs.Get("shop"); qsId != "" {
		shopId = qsId
	} else if formId := r.FormValue("shop"); formId != "" {
		shopId = formId
	} else if headerId := r.Header.Get("shop"); headerId != "" {
		shopId = headerId
	}

	if shopId == "" {
		return fmt.Errorf("error: %s", "you must provide a shop identifier")
	}
	if !bson.IsObjectIdHex(shopId) {
		return fmt.Errorf("error: %s", "invalid shop identifier")
	}
	shop := cart.Shop{
		Id: bson.ObjectIdHex(shopId),
	}

	if shop.Id.Hex() == "" {
		return fmt.Errorf("error: %s", "invalid shop identifier")
	}

	if err := shop.Get(); err != nil {
		return err
	}
	if shop.Id.Hex() == "" {
		return fmt.Errorf("error: %s", "invalid shop identifier")
	}

	c.Map(&shop)
	return nil
}

func mapCartAccount(c martini.Context, res http.ResponseWriter, r *http.Request) error {

	auth := r.Header.Get("Authorization")
	token := strings.Replace(auth, "Bearer ", "", 1)

	cust, err := cart.AuthenticateAccount(token)
	if err != nil {
		return err
	}

	shop := cart.Shop{
		Id: cust.ShopId,
	}

	if shop.Id.Hex() == "" {
		return fmt.Errorf("error: %s", "invalid shop identifier")
	}

	if err := shop.Get(); err != nil {
		return err
	}
	if shop.Id.Hex() == "" {
		return fmt.Errorf("error: %s", "invalid shop identifier")
	}

	c.Map(&shop)
	c.Map(token)
	return nil
}

func processDataContext(r *http.Request, c martini.Context) (*apicontext.DataContext, error) {
	qs := r.URL.Query()
	apiKey := qs.Get("key")
	brand := qs.Get("brandID")
	website := qs.Get("websiteID")

	//handles api key
	if apiKey == "" {
		apiKey = r.FormValue("key")
	}
	if apiKey == "" {
		apiKey = r.Header.Get("key")
	}
	if apiKey == "" {
		return nil, errors.New("No API Key Supplied.")
	}

	//gets customer user from api key
	user, err := getCustomerID(apiKey)
	if err != nil {
		return nil, errors.New("No User for this API Key.")
	}
	// go user.LogApiRequest(r)

	//handles branding
	var brandID int
	if brand == "" {
		brand = r.FormValue("brandID")
	}
	if brand == "" {
		brand = r.Header.Get("brandID")
	}
	if id, err := strconv.Atoi(brand); err == nil {
		brandID = id
	}

	//handles websiteID
	var websiteID int
	if website == "" {
		website = r.FormValue("websiteID")
	}
	if website == "" {
		website = r.Header.Get("websiteID")
	}
	if id, err := strconv.Atoi(website); err == nil {
		websiteID = id
	}

	//load brands in dtx
	//returns our data context...shared amongst controllers
	// var dtx apicontext.DataContext
	dtx := &apicontext.DataContext{
		APIKey:     apiKey,
		BrandID:    brandID,
		WebsiteID:  websiteID,
		UserID:     user.Id, //current authenticated user
		CustomerID: user.CustomerID,
		Globals:    nil,
	}

	return dtx, nil
}

func getCustomerID(apiKey string) (*customer.CustomerUser, error) {
	err := database.Init()
	if err != nil {
		return nil, err
	}

	var customer_id int

	row := database.MagentoDB.QueryRow(FindMagentoUserByKey, apiKey)

	err = row.Scan(&customer_id)
	if err != nil {
		return nil, fmt.Errorf("failed to find user for that API key")
	}

	//Obviously not right but we need to figure out which is which
	var customer_user customer.CustomerUser
	customer_user.CustomerID = customer_id
	customer_user.CustID = customer_id

	return &customer_user, nil
}

func InternalKeyAuthentication(w http.ResponseWriter, req *http.Request) {
	err := database.Init()
	if err != nil {
		http.Error(w, "Key could not be authenticated", http.StatusInternalServerError)
		return
	}

	key := req.URL.Query().Get("key")

	if key == "" {
		http.Error(w, "Access denied", http.StatusUnauthorized)
		return
	}

	row := database.MagentoDB.QueryRow(FindInternalMagentoUserByKey, key)

	var cust_id int

	err = row.Scan(&cust_id)
	if err != nil {
		http.Error(w, "Key could not be authenticated", http.StatusInternalServerError)
		return
	}

	return
}

//logRequest is simply the launcher for the analytics function ToPubSub
//Here we filter a little bit, making sure not to log any healthchecks
func logRequest(w http.ResponseWriter, r *http.Request, reqTime time.Time) {
	if strings.Contains(r.URL.Path, "checkup") || strings.Contains(r.URL.Path, "status") {
		return
	}

	ToPubSub(w, r, reqTime)
}
