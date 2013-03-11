package main

import (
	"./controllers/category"
	"./controllers/part"
	"./controllers/vehicle"
	"./helpers/auth"
	"./plate"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"
)

func run_test_request(t *testing.T, server *plate.Server, method, url_str string, payload url.Values) *httptest.ResponseRecorder {

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

	recorder := httptest.NewRecorder()

	server.ServeHTTP(recorder, &r)

	return recorder
}

func code_is(t *testing.T, r *httptest.ResponseRecorder, expected_code int) {
	if r.Code != expected_code {
		t.Errorf("Code %d expected, got: %d", expected_code, r.Code)
	}
}

func content_type_is_json(t *testing.T, r *httptest.ResponseRecorder) {
	ct := r.HeaderMap.Get("Content-Type")
	if ct != "application/json" {
		t.Errorf("Content type 'application/json' expected, got: %s", ct)
	}
}

func body_is(t *testing.T, r *httptest.ResponseRecorder, expected_body string) {
	body := r.Body.String()
	if body != expected_body {
		t.Errorf("Body '%s' expected, got: '%s'", expected_body, body)
	}
}

func TestHandler(t *testing.T) {

	server := plate.NewServer("doughboy")

	server.AddFilter(auth.AuthHandler)
	server.AddFilter(CorsHandler)
	server.Get("/", func(w http.ResponseWriter, r *http.Request) {
		http.Redirect(w, r, "http://labs.curtmfg.com/", http.StatusFound)
	}).NoFilter()

	server.Get("/vehicle", vehicle_ctlr.Year)
	server.Get("/vehicle/:year", vehicle_ctlr.Make)
	server.Get("/vehicle/:year/:make", vehicle_ctlr.Model)
	server.Get("/vehicle/:year/:make/:model", vehicle_ctlr.Submodel)
	server.Get("/vehicle/:year/:make/:model/:submodel", vehicle_ctlr.Config)
	server.Get("/vehicle/:year/:make/:model/:submodel/:config(.+)", vehicle_ctlr.Config)

	server.Get("/category", category_ctlr.Parents)
	server.Get("/category/:id", category_ctlr.GetCategory)
	server.Get("/category/:id/subs", category_ctlr.SubCategories)
	server.Get("/category/:id/parts", category_ctlr.GetParts)
	server.Get("/category/:id/parts/:page/:count", category_ctlr.GetParts)

	server.Get("/part/:part/vehicles", part_ctlr.Vehicles)
	server.Get("/part/:part", part_ctlr.Get)

	qs := url.Values{}
	qs.Add("key", "8aee0620-412e-47fc-900a-947820ea1c1d")

	recorder := run_test_request(t, server, "GET", "http://localhost:8080/vehicle", nil)
	code_is(t, recorder, 401)

	recorder = run_test_request(t, server, "GET", "http://localhost:8080/vehicle", qs)
	code_is(t, recorder, 200)
	content_type_is_json(t, recorder)

	recorder = run_test_request(t, server, "GET", "http://localhost:8080/vehicle/2012", qs)
	code_is(t, recorder, 200)
	content_type_is_json(t, recorder)

	recorder = run_test_request(t, server, "GET", "http://localhost:8080/vehicle/2012/Chevrolet", qs)
	code_is(t, recorder, 200)
	content_type_is_json(t, recorder)

	recorder = run_test_request(t, server, "GET", "http://localhost:8080/vehicle/2012/Chevrolet/Silverado 1500", qs)
	code_is(t, recorder, 200)
	content_type_is_json(t, recorder)

	recorder = run_test_request(t, server, "GET", "http://localhost:8080/vehicle/2012/Chevrolet/Silverado 1500/LT", qs)
	code_is(t, recorder, 200)
	content_type_is_json(t, recorder)

	recorder = run_test_request(t, server, "GET", "http://localhost:8080/vehicle/2012/Chevrolet/Silverado 1500/LT/with factory tow package", qs)
	code_is(t, recorder, 200)
	content_type_is_json(t, recorder)

	recorder = run_test_request(t, server, "GET", "http://localhost:8080/category", qs)
	code_is(t, recorder, 200)
	content_type_is_json(t, recorder)

	recorder = run_test_request(t, server, "GET", "http://localhost:8080/category/3", qs)
	code_is(t, recorder, 200)
	content_type_is_json(t, recorder)

	recorder = run_test_request(t, server, "GET", "http://localhost:8080/category/Hitches", qs)
	code_is(t, recorder, 200)
	content_type_is_json(t, recorder)

	recorder = run_test_request(t, server, "GET", "http://localhost:8080/category/Hitches/parts", qs)
	code_is(t, recorder, 200)
	content_type_is_json(t, recorder)

	recorder = run_test_request(t, server, "GET", "http://localhost:8080/category/Hitches/subs", qs)
	code_is(t, recorder, 200)
	content_type_is_json(t, recorder)

	recorder = run_test_request(t, server, "GET", "http://localhost:8080/category/Class III Trailer Hitches/parts/2/20", qs)
	code_is(t, recorder, 200)
	content_type_is_json(t, recorder)
}
