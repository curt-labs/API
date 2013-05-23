package main

import (
	"./controllers/category"
	"./controllers/customer"
	"./controllers/dealers"
	"./controllers/part"
	"./controllers/vehicle"
	"./controllers/videos"
	"./helpers/auth"
	"./helpers/database"
	"./helpers/plate"
	"errors"
	"fmt"
	"log"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"
)

func run_test_request(t *testing.T, server *plate.Server, method, url_str string, payload url.Values) (*httptest.ResponseRecorder, http.Request) {

	url_obj, err := url.Parse(url_str)
	if err != nil {
		t.Fatal(err)
	}

	r := http.Request{
		Method: method,
		URL:    url_obj,
	}

	if payload != nil {
		r.URL.RawQuery = payload.Encode()
	}

	if strings.ToUpper(method) == "POST" {
		r.Form = payload
	}

	recorder := httptest.NewRecorder()

	server.ServeHTTP(recorder, &r)

	return recorder, r
}

func code_is(t *testing.T, r *httptest.ResponseRecorder, expected_code int) error {
	if r.Code != expected_code {
		return errors.New(fmt.Sprintf("Code %d expected, got: %d", expected_code, r.Code))
	}
	return nil
}

func content_type_is_json(t *testing.T, r *httptest.ResponseRecorder) error {
	ct := r.HeaderMap.Get("Content-Type")
	if ct != "application/json" {
		return errors.New(fmt.Sprintf("Content type 'application/json' expected, got: %s", ct))
	}
	return nil
}

func body_is(t *testing.T, r *httptest.ResponseRecorder, expected_body string) error {
	body := r.Body.String()
	if body != expected_body {
		return errors.New(fmt.Sprintf("Body '%s' expected, got: '%s'", expected_body, body))
	}
	return nil
}

type ErrorMessage struct {
	StatusCode int
	Error      string
	Route      *url.URL
}

func checkError(req http.Request, rec *httptest.ResponseRecorder, err error, t *testing.T) {
	if err != nil {
		t.Errorf("\nError: %s \nRoute: %s \n\n", err.Error(), req.URL.String())
	} else {
		t.Log("PASS: " + req.URL.String())
	}
}

func TestHandler(t *testing.T) {

	err := database.PrepareAll()
	if err != nil {
		log.Fatal(err)
	}

	server := plate.NewServer("doughboy")
	server.Logging = true

	server.AddFilter(auth.AuthHandler)
	server.AddFilter(CorsHandler)
	server.Get("/", func(w http.ResponseWriter, r *http.Request) {
		http.Redirect(w, r, "http://labs.curtmfg.com/", http.StatusFound)
	}).NoFilter()

	server.Get("/.status", func(w http.ResponseWriter, r *http.Request) {
		server.StatusService.GetStatus(w, r)
	}).NoFilter()

	server.Get("/vehicle", vehicle_ctlr.Year)
	server.Get("/vehicle/:year", vehicle_ctlr.Make)
	server.Get("/vehicle/:year/:make", vehicle_ctlr.Model)
	server.Get("/vehicle/:year/:make/:model", vehicle_ctlr.Submodel)
	server.Get("/vehicle/:year/:make/:model/connector", vehicle_ctlr.Connector)
	server.Get("/vehicle/:year/:make/:model/:submodel", vehicle_ctlr.Config)
	server.Get("/vehicle/:year/:make/:model/:submodel/connector", vehicle_ctlr.Connector)
	server.Get("/vehicle/:year/:make/:model/:submodel/:config(.+)/connector", vehicle_ctlr.Connector)
	server.Get("/vehicle/:year/:make/:model/:submodel/:config(.+)", vehicle_ctlr.Config)

	server.Get("/category", category_ctlr.Parents)
	server.Get("/category/:id", category_ctlr.GetCategory)
	server.Get("/category/:id/subs", category_ctlr.SubCategories)
	server.Get("/category/:id/parts", category_ctlr.GetParts)
	server.Get("/category/:id/parts/:page/:count", category_ctlr.GetParts)

	server.Get("/part/:part/vehicles", part_ctlr.Vehicles)
	server.Get("/part/:part/attributes", part_ctlr.Attributes)
	server.Get("/part/:part/reviews", part_ctlr.Reviews)
	server.Get("/part/:part/categories", part_ctlr.Categories)
	server.Get("/part/:part/content", part_ctlr.GetContent)
	server.Get("/part/:part/images", part_ctlr.Images)
	server.Get("/part/:part((.*?)\\.(PDF|pdf)$)", part_ctlr.InstallSheet).NoFilter() // Resolves: /part/11000.pdf
	server.Get("/part/:part/packages", part_ctlr.Packaging)
	server.Get("/part/:part/pricing", part_ctlr.Prices)
	server.Get("/part/:part/related", part_ctlr.GetRelated)
	server.Get("/part/:part/videos", part_ctlr.Videos)
	server.Get("/part/:part/:year/:make/:model", part_ctlr.GetWithVehicle)
	server.Get("/part/:part/:year/:make/:model/:submodel", part_ctlr.GetWithVehicle)
	server.Get("/part/:part/:year/:make/:model/:submodel/:config(.+)", part_ctlr.GetWithVehicle)
	server.Get("/part/:part", part_ctlr.Get)

	server.Post("/customer/auth", customer_ctlr.UserAuthentication).NoFilter()
	server.Get("/customer/auth", customer_ctlr.KeyedUserAuthentication).NoFilter()

	server.Post("/customer/locations", customer_ctlr.GetLocations)
	server.Post("/customer/users", customer_ctlr.GetUsers) // Requires a user to be marked as sudo

	/**
	 * Video
	 */
	server.Get("/videos", videos_ctlr.DistinctVideos).NoFilter()

	/**** INTERNAL USE ONLY ****/
	server.Get("/dealers/etailer", dealers_ctlr.Etailers).NoFilter()
	server.Get("/dealers/etailer/platinum", dealers_ctlr.PlatinumEtailers).NoFilter()
	server.Get("/dealers/local", dealers_ctlr.LocalDealers).NoFilter()
	server.Get("/dealers/local/regions", dealers_ctlr.LocalRegions).NoFilter()
	server.Get("/dealers/local/tiers", dealers_ctlr.LocalDealerTiers).NoFilter()
	server.Get("/dealers/local/types", dealers_ctlr.LocalDealerTypes).NoFilter()
	server.Get("/dealers/search", dealers_ctlr.SearchLocations).NoFilter()
	server.Get("/dealers/search/type", dealers_ctlr.SearchLocationsByType).NoFilter()
	server.Get("/dealers/search/geo", dealers_ctlr.SearchLocationsByLatLng).NoFilter()
	server.Get("/dealers/search/geo/:latitude/:longitude", dealers_ctlr.SearchLocationsByLatLng).NoFilter()

	qs := url.Values{}
	qs.Add("key", "8aee0620-412e-47fc-900a-947820ea1c1d")

	recorder, req := run_test_request(t, server, "GET", "http://localhost:8080/vehicle", nil)
	err = code_is(t, recorder, 401)
	checkError(req, recorder, err, t)

	recorder, req = run_test_request(t, server, "GET", "http://localhost:8080/vehicle", qs)
	err = code_is(t, recorder, 200)
	checkError(req, recorder, err, t)
	err = content_type_is_json(t, recorder)
	checkError(req, recorder, err, t)

	recorder, req = run_test_request(t, server, "GET", "http://localhost:8080/vehicle/2012", qs)
	err = code_is(t, recorder, 200)
	checkError(req, recorder, err, t)
	err = content_type_is_json(t, recorder)
	checkError(req, recorder, err, t)

	recorder, req = run_test_request(t, server, "GET", "http://localhost:8080/vehicle/2012/Chevrolet", qs)
	err = code_is(t, recorder, 200)
	checkError(req, recorder, err, t)
	err = content_type_is_json(t, recorder)
	checkError(req, recorder, err, t)

	recorder, req = run_test_request(t, server, "GET", "http://localhost:8080/vehicle/2012/Chevrolet/Silverado 1500", qs)
	err = code_is(t, recorder, 200)
	checkError(req, recorder, err, t)
	err = content_type_is_json(t, recorder)
	checkError(req, recorder, err, t)

	recorder, req = run_test_request(t, server, "GET", "http://localhost:8080/vehicle/2012/Chevrolet/Silverado 1500/LT", qs)
	err = code_is(t, recorder, 200)
	checkError(req, recorder, err, t)
	err = content_type_is_json(t, recorder)
	checkError(req, recorder, err, t)

	recorder, req = run_test_request(t, server, "GET", "http://localhost:8080/vehicle/2012/Chevrolet/Silverado 1500/LT/with factory tow package", qs)
	err = code_is(t, recorder, 200)
	checkError(req, recorder, err, t)
	err = content_type_is_json(t, recorder)
	checkError(req, recorder, err, t)

	recorder, req = run_test_request(t, server, "GET", "http://localhost:8080/part/12289", qs)
	err = code_is(t, recorder, 200)
	checkError(req, recorder, err, t)
	err = content_type_is_json(t, recorder)
	checkError(req, recorder, err, t)

	recorder, req = run_test_request(t, server, "GET", "http://localhost:8080/part/12289/vehicles", qs)
	err = code_is(t, recorder, 200)
	checkError(req, recorder, err, t)
	err = content_type_is_json(t, recorder)
	checkError(req, recorder, err, t)

	recorder, req = run_test_request(t, server, "GET", "http://localhost:8080/part/12289/attributes", qs)
	err = code_is(t, recorder, 200)
	checkError(req, recorder, err, t)
	err = content_type_is_json(t, recorder)
	checkError(req, recorder, err, t)

	recorder, req = run_test_request(t, server, "GET", "http://localhost:8080/part/12289/reviews", qs)
	err = code_is(t, recorder, 200)
	checkError(req, recorder, err, t)
	err = content_type_is_json(t, recorder)
	checkError(req, recorder, err, t)

	recorder, req = run_test_request(t, server, "GET", "http://localhost:8080/part/12289/categories", qs)
	err = code_is(t, recorder, 200)
	checkError(req, recorder, err, t)
	err = content_type_is_json(t, recorder)
	checkError(req, recorder, err, t)

	recorder, req = run_test_request(t, server, "GET", "http://localhost:8080/part/12289/content", qs)
	err = code_is(t, recorder, 200)
	checkError(req, recorder, err, t)
	err = content_type_is_json(t, recorder)
	checkError(req, recorder, err, t)

	recorder, req = run_test_request(t, server, "GET", "http://localhost:8080/part/12289/images", qs)
	err = code_is(t, recorder, 200)
	checkError(req, recorder, err, t)
	err = content_type_is_json(t, recorder)
	checkError(req, recorder, err, t)

	recorder, req = run_test_request(t, server, "GET", "http://localhost:8080/part/12289.pdf", nil)
	err = code_is(t, recorder, 200)
	checkError(req, recorder, err, t)

	recorder, req = run_test_request(t, server, "GET", "http://localhost:8080/part/12289/packages", qs)
	err = code_is(t, recorder, 200)
	checkError(req, recorder, err, t)
	err = content_type_is_json(t, recorder)
	checkError(req, recorder, err, t)

	recorder, req = run_test_request(t, server, "GET", "http://localhost:8080/part/12289/pricing", qs)
	err = code_is(t, recorder, 200)
	checkError(req, recorder, err, t)
	err = content_type_is_json(t, recorder)
	checkError(req, recorder, err, t)

	recorder, req = run_test_request(t, server, "GET", "http://localhost:8080/part/12289/2004/Dodge/Grand Caravan/Anniversary Edition/Stow 'n Go", qs)
	err = code_is(t, recorder, 200)
	checkError(req, recorder, err, t)
	err = content_type_is_json(t, recorder)
	checkError(req, recorder, err, t)

	recorder, req = run_test_request(t, server, "GET", "http://localhost:8080/category", qs)
	err = code_is(t, recorder, 200)
	checkError(req, recorder, err, t)
	err = content_type_is_json(t, recorder)
	checkError(req, recorder, err, t)

	recorder, req = run_test_request(t, server, "GET", "http://localhost:8080/category/3", qs)
	err = code_is(t, recorder, 200)
	checkError(req, recorder, err, t)
	err = content_type_is_json(t, recorder)
	checkError(req, recorder, err, t)

	//log.Println(qs)
	recorder, req = run_test_request(t, server, "GET", "http://localhost:8080/category/Trailer Hitches?key=8aee0620-412e-47fc-900a-947820ea1c1d", nil)
	err = code_is(t, recorder, 200)
	checkError(req, recorder, err, t)
	err = content_type_is_json(t, recorder)
	checkError(req, recorder, err, t)

	recorder, req = run_test_request(t, server, "GET", "http://localhost:8080/category/Trailer Hitches/parts", qs)
	err = code_is(t, recorder, 200)
	checkError(req, recorder, err, t)
	err = content_type_is_json(t, recorder)
	checkError(req, recorder, err, t)

	recorder, req = run_test_request(t, server, "GET", "http://localhost:8080/category/Trailer Hitches/subs", qs)
	err = code_is(t, recorder, 200)
	checkError(req, recorder, err, t)
	err = content_type_is_json(t, recorder)
	checkError(req, recorder, err, t)

	recorder, req = run_test_request(t, server, "GET", "http://localhost:8080/category/Class III Trailer Hitches/parts/2/20", qs)
	err = code_is(t, recorder, 200)
	checkError(req, recorder, err, t)
	err = content_type_is_json(t, recorder)
	checkError(req, recorder, err, t)

	recorder, req = run_test_request(t, server, "GET", "http://localhost:8080/dealers/etailer", qs)
	err = code_is(t, recorder, 200)
	checkError(req, recorder, err, t)
	err = content_type_is_json(t, recorder)
	checkError(req, recorder, err, t)

	recorder, req = run_test_request(t, server, "GET", "http://localhost:8080/dealers/etailer/platinum", qs)
	err = code_is(t, recorder, 200)
	checkError(req, recorder, err, t)
	err = content_type_is_json(t, recorder)
	checkError(req, recorder, err, t)

	recorder, req = run_test_request(t, server, "GET", "http://localhost:8080/dealers/local?latlng=44.279362,-93.166165,46.21285,-88.063015&center=45.254334,-90.61459", qs)
	err = code_is(t, recorder, 200)
	checkError(req, recorder, err, t)
	err = content_type_is_json(t, recorder)
	checkError(req, recorder, err, t)

	recorder, req = run_test_request(t, server, "GET", "http://localhost:8080/dealers/local/regions", qs)
	err = code_is(t, recorder, 200)
	checkError(req, recorder, err, t)
	err = content_type_is_json(t, recorder)
	checkError(req, recorder, err, t)

	recorder, req = run_test_request(t, server, "GET", "http://localhost:8080/dealers/local/tiers", qs)
	err = code_is(t, recorder, 200)
	checkError(req, recorder, err, t)
	err = content_type_is_json(t, recorder)
	checkError(req, recorder, err, t)

	recorder, req = run_test_request(t, server, "GET", "http://localhost:8080/dealers/local/types", qs)
	err = code_is(t, recorder, 200)
	checkError(req, recorder, err, t)
	err = content_type_is_json(t, recorder)
	checkError(req, recorder, err, t)

	recorder, req = run_test_request(t, server, "GET", "http://localhost:8080/dealers/search/jc auto", qs)
	err = code_is(t, recorder, 200)
	checkError(req, recorder, err, t)
	err = content_type_is_json(t, recorder)
	checkError(req, recorder, err, t)

	recorder, req = run_test_request(t, server, "GET", "http://localhost:8080/dealers/search/type/jc auto", qs)
	err = code_is(t, recorder, 200)
	checkError(req, recorder, err, t)
	err = content_type_is_json(t, recorder)
	checkError(req, recorder, err, t)

	recorder, req = run_test_request(t, server, "GET", "http://localhost:8080/dealers/search/geo/43.853282/-95.571675", qs)
	err = code_is(t, recorder, 200)
	checkError(req, recorder, err, t)
	err = content_type_is_json(t, recorder)
	checkError(req, recorder, err, t)

	recorder, req = run_test_request(t, server, "GET", "http://localhost:8080/videos", nil)
	err = code_is(t, recorder, 200)
	checkError(req, recorder, err, t)
	err = content_type_is_json(t, recorder)
	checkError(req, recorder, err, t)

	// This test is failing because for some reason the encrypted password for the test user
	// did not properly carry over the password

	// authForm := url.Values{}
	// authForm.Add("email", "test@curtmfg.com")
	// authForm.Add("password", "")
	// recorder, req = run_test_request(t, server, "POST", "http://localhost:8080/customer/auth", authForm)
	// err = code_is(t, recorder, 200)
	// err = content_type_is_json(t, recorder)

	// authForm := url.Values{}
	// authForm.Add("key", "c8bd5d89-8d16-11e2-801f-00155d47bb0a")
	// recorder, req = run_test_request(t, server, "GET", "http://localhost:8080/customer/auth", authForm)
	// err = code_is(t, recorder, 200)
	// checkError(req, recorder, err, t)
	// err = content_type_is_json(t, recorder)
	// checkError(req, recorder, err, t)

}
