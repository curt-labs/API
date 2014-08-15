package segmenter

import (
	"github.com/curt-labs/GoAPI/helpers/slack"
	"github.com/go-martini/martini"
	"github.com/segmentio/analytics-go"
	"net/http"
	"time"
)

func Log() martini.Handler {
	return func(res http.ResponseWriter, r *http.Request, c martini.Context) {
		start := time.Now()
		c.Next()
		go func(req *http.Request) {
			client := analytics.New("oactr73lbg")

			key := r.Header.Get("key")
			if key == "" {
				vals := r.URL.Query()
				key = vals.Get("key")
			}
			if key == "" {
				key = r.FormValue("key")
			}

			err := client.Track(map[string]interface{}{
				"event":       r.URL.String(),
				"userId":      key,
				"method":      r.Method,
				"header":      r.Header,
				"form":        r.Form,
				"requestTime": time.Since(start),
			})
			if err != nil {
				m := slack.Message{
					Channel:  "debugging",
					Username: "GoAPI",
					Text:     err.Error(),
				}
				m.Send()
			}
		}(r)

	}
}
