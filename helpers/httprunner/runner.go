package httprunner

import (
	"bytes"
	"encoding/json"
	"github.com/curt-labs/API/controllers/middleware"
	"github.com/curt-labs/API/helpers/encoding"
	"github.com/go-martini/martini"
	"github.com/martini-contrib/render"
	"github.com/rakyll/pb"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"time"
)

type result struct {
	err           error
	statusCode    int
	duration      time.Duration
	contentLength int64
}

type ReqOpts struct {
	Method           string
	URL              string
	ParameterizedURL string
	Header           http.Header
	Username         string
	Password         string
	Handler          martini.Handler
	Middleware       []martini.Handler
	Body             string
	// OriginalHost represents the original host name user is provided.
	// Request host is an resolved IP. TLS/SSL handshakes may require
	// the original server name, keep it to initate the TLS client.
	OriginalHost string
}

type BenchmarkOptions struct {
	Method             string
	Route              string
	ParameterizedRoute string
	Header             http.Header
	Username           string
	Password           string
	Handler            martini.Handler
	Middleware         []martini.Handler
	QueryString        *url.Values
	JsonBody           interface{}
	FormBody           url.Values
	Output             string
	Runs               int
	ConcurrentUsers    int
}

// Creates a req object from req options
func (r *ReqOpts) GenerateRequest() *http.Request {
	var req *http.Request
	if r.Body != "" && strings.ToUpper(r.Method) != "GET" {
		req, _ = http.NewRequest(r.Method, r.URL, bytes.NewBufferString(r.Body))
	} else if r.Body != "" {
		req, _ = http.NewRequest(r.Method, r.URL+"?"+r.Body, nil)
	} else {
		req, _ = http.NewRequest(r.Method, r.URL, nil)
	}
	req.Header = r.Header

	// update the Host value in the Request - this is used as the host header in any subsequent request
	req.Host = r.OriginalHost

	if r.Username != "" && r.Password != "" {
		req.SetBasicAuth(r.Username, r.Password)
	}
	return req
}

type Runner struct {
	// Req represents the options of the request to be made.
	// TODO(jbd): Make it work with an http.Request instead.
	Req *ReqOpts

	// N is the total number of requests to make.
	N int

	// C is the concurrency level, the number of concurrent workers to run.
	C int

	// Timeout in seconds.
	Timeout int

	// Qps is the rate limit.
	Qps int

	// AllowInsecure is an option to allow insecure TLS/SSL certificates.
	AllowInsecure bool

	// DisableCompression is an option to disable compression in response
	DisableCompression bool

	// DisableKeepAlives is an option to prevents re-use of TCP connections between different HTTP requests
	DisableKeepAlives bool

	// Output represents the output type. If "csv" is provided, the
	// output will be dumped as a csv stream.
	Output string

	// ProxyAddr is the address of HTTP proxy server in the format on "host:port".
	// Optional.
	ProxyAddr *url.URL

	bar     *pb.ProgressBar
	results chan *result
}

func newPb(size int) (bar *pb.ProgressBar) {
	bar = pb.New(size)
	bar.Format("Bom !")
	bar.Start()
	return
}

func Request(method, route string, body *url.Values, handler martini.Handler) httptest.ResponseRecorder {
	if body != nil && strings.ToUpper(method) != "GET" {
		return Req(handler, method, "", route, nil, body)
	} else if body != nil {
		return Req(handler, method, "", route, body)
	}
	return Req(handler, method, "", route)
}

func ParameterizedRequest(method, prepared_route string, route string, qs *url.Values, body *url.Values, handler martini.Handler) httptest.ResponseRecorder {
	return Req(handler, method, prepared_route, route, qs, body)
}

func ParameterizedJsonRequest(method, prepared_route string, route string, qs *url.Values, iface interface{}, handler martini.Handler) httptest.ResponseRecorder {
	headers := map[string]interface{}{
		"Content-Type": "application/json",
	}
	return Req(handler, method, prepared_route, route, qs, nil, iface, headers)
}

func JsonRequest(method, route string, qs *url.Values, iface interface{}, handler martini.Handler) httptest.ResponseRecorder {
	headers := map[string]interface{}{
		"Content-Type": "application/json",
	}
	return Req(handler, method, "", route, qs, nil, iface, headers)
}

func Req(handler martini.Handler, method, prepared_route, route string, args ...interface{}) httptest.ResponseRecorder {

	// args[0] - *url.Values
	// args[1] - *url.Values
	// args[2] - interface{} req.Body
	// args[3] - map[string]interface{} Headers

	var response httptest.ResponseRecorder
	if prepared_route == "" {
		prepared_route = route
	}

	m := martini.New()
	r := martini.NewRouter()
	switch strings.ToUpper(method) {
	case "GET":
		r.Get(prepared_route, handler)
	case "POST":
		r.Post(prepared_route, handler)
	case "PUT":
		r.Put(prepared_route, handler)
	case "PATCH":
		r.Patch(prepared_route, handler)
	case "DELETE":
		r.Delete(prepared_route, handler)
	case "HEAD":
		r.Head(prepared_route, handler)
	default:
		r.Any(prepared_route, handler)
	}

	m.Use(render.Renderer())
	m.Use(encoding.MapEncoder)
	m.Use(middleware.Meddler())
	m.Action(r.Handle)

	if len(args) > 0 && args[0] != nil {
		qs := args[0].(*url.Values)
		route = route + "?" + qs.Encode()
	}

	var request *http.Request
	if len(args) > 1 && args[1] != nil && args[1].(*url.Values) != nil {
		request, _ = http.NewRequest(method, route, bytes.NewBuffer([]byte(args[1].(*url.Values).Encode())))
	} else if len(args) > 1 && args[1] == nil {
		js, err := json.Marshal(args[2])
		if err != nil {
			return response
		}
		request, _ = http.NewRequest(method, route, bytes.NewBuffer(js))
	} else {
		request, _ = http.NewRequest(method, route, nil)
	}

	if len(args) == 4 {
		headers := args[3].(map[string]interface{})
		for key, val := range headers {
			request.Header.Set(key, val.(string))
		}
	}

	resp := httptest.NewRecorder()
	if resp == nil {
		return response
	}
	response = *resp
	m.ServeHTTP(&response, request)

	return response
}

func (opts *BenchmarkOptions) RequestBenchmark() {

	var body string
	if opts.JsonBody != nil {
		js, err := json.Marshal(opts.JsonBody)
		if err != nil {
			return
		}
		body = string(js)
	} else if opts.QueryString != nil {
		opts.Route = opts.Route + "?" + opts.QueryString.Encode()
	} else if opts.FormBody != nil {
		body = opts.FormBody.Encode()
	}

	if opts.ConcurrentUsers == 0 {
		opts.ConcurrentUsers = 1
	}

	runner := &Runner{
		Req: &ReqOpts{
			Body:             body,
			Handler:          opts.Handler,
			URL:              opts.Route,
			ParameterizedURL: opts.ParameterizedRoute,
			Method:           opts.Method,
			Header:           opts.Header,
			Username:         opts.Username,
			Password:         opts.Password,
			Middleware:       opts.Middleware,
		},
		N:      opts.Runs,
		C:      opts.ConcurrentUsers,
		Output: opts.Output,
	}
	runner.Run()
}
