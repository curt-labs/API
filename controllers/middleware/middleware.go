package middleware

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/curt-labs/GoAPI/helpers/apicontext"
	"github.com/curt-labs/GoAPI/helpers/error"
	"github.com/curt-labs/GoAPI/helpers/slack"
	"github.com/curt-labs/GoAPI/models/cart"
	"github.com/curt-labs/GoAPI/models/customer"
	"github.com/go-martini/martini"
	"github.com/segmentio/analytics-go"
	"gopkg.in/mgo.v2/bson"
)

var (
	ExcusedRoutes = []string{"/status", "/customer/auth", "/customer/user", "/new/customer/auth", "/customer/user/register", "/customer/user/resetPassword"}
)

func Meddler() martini.Handler {
	return func(res http.ResponseWriter, r *http.Request, c martini.Context) {
		res.Header().Add("Access-Control-Allow-Origin", "*")
		if strings.ToLower(r.Method) == "options" {
			return
		}

		if strings.Contains(r.URL.String(), "favicon") {
			res.Write([]byte(""))
			return
		}
		start := time.Now()

		excused := false
		for _, route := range ExcusedRoutes {
			if strings.Contains(r.URL.String(), route) {
				excused = true
			}
		}

		// check if we need to make a call
		// to the shopping cart middleware
		if strings.Contains(strings.ToLower(r.URL.Path), "/shopify/account") { // account perms
			if strings.ToLower(r.URL.Path) == "/shopify/account/login" {
				shopID := r.URL.Query().Get("shop")
				if bson.IsObjectIdHex(shopID) {
					c.Map(&cart.Shop{
						Id: bson.ObjectIdHex(shopID),
					})
				}
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
		go logRequest(r, time.Since(start))
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
	user, err := customer.GetCustomerUserFromKey(apiKey)
	if err != nil || user.Id == "" {
		return nil, errors.New("No User for this API Key.")
	}
	go user.LogApiRequest(r)

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
	err = dtx.GetBrandsArrayAndString(apiKey, brandID)
	if err != nil {
		return nil, err
	}
	return dtx, nil
}

func logRequest(r *http.Request, reqTime time.Duration) {
	client := analytics.New("oactr73lbg")

	key := r.Header.Get("key")
	if key == "" {
		vals := r.URL.Query()
		key = vals.Get("key")
	}
	if key == "" {
		key = r.FormValue("key")
	}

	vals := r.URL.Query()
	props := make(map[string]interface{}, 0)
	for k, v := range vals {
		props[k] = v
	}

	props["method"] = r.Method
	props["header"] = r.Header
	props["query"] = r.URL.Query().Encode()
	props["referer"] = r.Referer()
	props["userAgent"] = r.UserAgent()
	props["form"] = r.Form
	props["requestTime"] = int64((reqTime.Nanoseconds() * 1000) * 1000)

	err := client.Track(&analytics.Track{
		Event:      r.URL.String(),
		UserId:     key,
		Properties: props,
	})

	if err != nil {
		m := slack.Message{
			Channel:  "debugging",
			Username: "GoAPI",
			Text:     err.Error(),
		}
		m.Send()
	}
}
