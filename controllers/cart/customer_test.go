package cart_ctlr

import (
	"github.com/curt-labs/GoAPI/controllers/middleware"
	"github.com/curt-labs/GoAPI/helpers/encoding"
	"github.com/curt-labs/GoAPI/models/cart"
	"github.com/go-martini/martini"
	"github.com/martini-contrib/render"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"net/http"
	"net/http/httptest"
	"regexp"
	"testing"
)

var _ = Describe("Customer", func() {
	// var (
	// 	body []byte
	// 	err  error
	// )

	Context("List All Customers", func() {
		It("returns a 500 status code", func() {
			Request("GET", "/shopify/customers", GetCustomers)
			Expect(response.Code).To(Equal(500))
		})
	})
})

var (
	response *httptest.ResponseRecorder
)

func TestMain(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Customer Suite")
}

func Request(method, route string, handler martini.Handler) {
	m := martini.Classic()
	m.Get(route, handler)
	m.Use(render.Renderer())
	m.Use(MapEncoder)
	m.Use(middleware.Meddler())
	m.Map(&cart.Shop{})

	request, _ := http.NewRequest(method, route, nil)
	response = httptest.NewRecorder()
	m.ServeHTTP(response, request)
}

var rxAccept = regexp.MustCompile(`(?:xml|html|plain|json)\/?$`)

func MapEncoder(c martini.Context, w http.ResponseWriter, r *http.Request) {
	accept := r.Header.Get("Accept")
	if accept == "*/*" {
		accept = r.Header.Get("Content-Type")
	}
	matches := rxAccept.FindStringSubmatch(accept)

	dt := "json"
	if len(matches) == 1 {
		dt = matches[0]
	}
	switch dt {
	case "xml":

		c.MapTo(encoding.XmlEncoder{}, (*encoding.Encoder)(nil))
		w.Header().Set("Content-Type", "application/xml")
	case "plain":
		c.MapTo(encoding.TextEncoder{}, (*encoding.Encoder)(nil))
		w.Header().Set("Content-Type", "text/plain")
	case "html":
		c.MapTo(encoding.TextEncoder{}, (*encoding.Encoder)(nil))
		w.Header().Set("Content-Type", "text/html")
	default:
		c.MapTo(encoding.JsonEncoder{}, (*encoding.Encoder)(nil))
		w.Header().Set("Content-Type", "application/json")
	}
}
