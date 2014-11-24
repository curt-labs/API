package httprunner

import (
	"bytes"
	"github.com/go-martini/martini"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/rakyll/pb"
)

type result struct {
	err           error
	statusCode    int
	duration      time.Duration
	contentLength int64
}

type ReqOpts struct {
	Method     string
	URL        string
	Header     http.Header
	Username   string
	Password   string
	Handler    martini.Handler
	Middleware []martini.Handler
	Body       *url.Values
	// OriginalHost represents the original host name user is provided.
	// Request host is an resolved IP. TLS/SSL handshakes may require
	// the original server name, keep it to initate the TLS client.
	OriginalHost string
}

// Creates a req object from req options
func (r *ReqOpts) Request() *http.Request {
	var req *http.Request
	if r.Body != nil && strings.ToUpper(r.Method) != "GET" {
		req, _ = http.NewRequest(r.Method, r.URL, bytes.NewBufferString(r.Body.Encode()))
	} else if r.Body != nil {
		req, _ = http.NewRequest(r.Method, r.URL+"?"+r.Body.Encode(), nil)
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
