package middleware

import (
	"fmt"
	"github.com/curt-labs/GoAPI/helpers/slack"
	"github.com/curt-labs/GoAPI/models/cart"
	"github.com/curt-labs/GoAPI/models/customer_new"
	"github.com/go-martini/martini"
	"github.com/segmentio/analytics-go"
	"gopkg.in/mgo.v2/bson"
	"net/http"
	"strings"
	"time"
)

var (
	ExcusedRoutes = []string{"/customer/auth", "/customer/user", "/customer", "/new/customer/auth", "/new/customer/user/register"}
)

func Meddler() martini.Handler {
	return func(res http.ResponseWriter, r *http.Request, c martini.Context) {
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
		if strings.Contains(strings.ToLower(r.URL.String()), "/shopify") {
			if err := mapCart(c, res, r); err != nil {
				http.Error(res, err.Error(), http.StatusInternalServerError)
				return
			}
			excused = true
		}

		if !excused {
			authed := checkAuth(r)
			if !authed {
				http.Error(res, "Unauthorized", http.StatusUnauthorized)
				return
			}
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

func checkAuth(r *http.Request) bool {
	qs := r.URL.Query()
	key := qs.Get("key")
	if key == "" {
		key = r.FormValue("key")
	}
	if key == "" {
		key = r.Header.Get("key")
	}
	if key == "" {
		return false
	}

	user, err := customer_new.GetCustomerUserFromKey(key)
	if err != nil || user.Id == "" {
		return false
	}

	go user.LogApiRequest(r)

	return true
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

	err := client.Track(map[string]interface{}{
		"event":       r.URL.String(),
		"userId":      key,
		"properties":  props,
		"method":      r.Method,
		"header":      r.Header,
		"query":       r.URL.Query().Encode(),
		"referer":     r.Referer(),
		"userAgent":   r.UserAgent(),
		"form":        r.Form,
		"requestTime": int64((reqTime.Nanoseconds() * 1000) * 1000),
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
