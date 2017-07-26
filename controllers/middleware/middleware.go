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
)

// SIDEEFFECTS
// res.Header().Add("Access-Control-Allow-Origin", "*")
// res.Header().Add("Cache-Control", "max-age=86400")
// Log request
func Meddler() martini.Handler {
	return func(res http.ResponseWriter, r *http.Request, c martini.Context) {
		res.Header().Add("Access-Control-Allow-Origin", "*")
		res.Header().Add("Cache-Control", "max-age=86400")
		if strings.ToLower(r.Method) == "options" {
			return
		}

		// Skip authentication for favicon
		if strings.Contains(r.URL.String(), "favicon") {
			res.Write([]byte(""))
			return
		}

		// Flag paths that do not r)equire auth
		// FIXME favicon could probably go here too
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
const API_KEY_PARAM = "key"
const BRAND_ID_PARAM = "brandID"
const WEBSITE_ID_PARAM = "websiteID"
const ERR_MISSING_KEY = "No API Key Supplied."

func getKey(r *http.Request) (apiKey string, err error) {
	qs := r.URL.Query()

	apiKey = qs.Get(API_KEY_PARAM)

	if apiKey == "" {
		apiKey = r.FormValue(API_KEY_PARAM)
	}

	if apiKey == "" {
		apiKey = r.Header.Get(API_KEY_PARAM)
	}

	if apiKey == "" {
		err = errors.New(ERR_MISSING_KEY)
	}

	return apiKey, err
}

func getId(param string, r *http.Request) (int, error){
	qs := r.URL.Query()
	id := qs.Get(param)

	if id == "" {
		id = r.FormValue(param)
	}

	if id == "" {
		id = r.Header.Get(param)
	}

	return strconv.Atoi(id)
}

func processDataContext(r *http.Request, c martini.Context) (*apicontext.DataContext, error) {
	qs := r.URL.Query()
	website := qs.Get("websiteID")

	apiKey, err := getKey(r)
	if err != nil {
		return nil, err
	}

	//gets customer user from api key
	user, err := getCustomerID(apiKey)
	if err != nil || user.Id == "" {
		return nil, errors.New("No User for this API Key.")
	}
	// go user.LogApiRequest(r)

	// TODO some duplicate code here
	var brandID int
	if id, err := getId(BRAND_ID_PARAM, r); err == nil {
		brandID = id
	}

	var websiteID int
	if id, err := getId(WEBSITE_ID_PARAM, r); err == nil {
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
	err = dtx.GetBrandsArrayAndString(apiKey, brandID)
	if err != nil {
		return nil, err
	}

	return dtx, nil
}

func getCustomerID(apiKey string) (*customer.CustomerUser, error) {
	err := database.Init()
	if err != nil {
		return nil, err
	}
	session := database.ProductMongoSession.Copy()
	defer session.Close()

	query := bson.M{"users.keys.key": apiKey}

	var resp = struct {
		Users []customer.CustomerUser `bson:"users"`
	}{}
	err = session.DB(database.ProductDatabase).C(database.CustomerCollectionName).Find(query).Select(bson.M{"users.$": 1, "_id": 0}).One(&resp)
	if len(resp.Users) == 0 {
		return nil, fmt.Errorf("failed to find user for that API key")
	}
	return &resp.Users[0], err
}

//logRequest is simply the launcher for the analytics function ToPubSub
//Here we filter a little bit, making sure not to log any healthchecks
func logRequest(w http.ResponseWriter, r *http.Request, reqTime time.Time) {
	if strings.Contains(r.URL.Path, "checkup") || strings.Contains(r.URL.Path, "status") {
		return
	}

	ToPubSub(w, r, reqTime)
}
