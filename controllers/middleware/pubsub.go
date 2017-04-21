package middleware

import (
	"context"
	"encoding/base64"
	"encoding/binary"
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"os"
	"time"

	uuid "github.com/satori/go.uuid"
	"golang.org/x/oauth2/google"
	"golang.org/x/oauth2/jwt"
	"google.golang.org/api/option"

	"cloud.google.com/go/pubsub"
)

// 4/13/17
/******************************************************
* NOTE: This can and should be broken out
* later. It is nearly identical to the analytics code
* in v2mock. Also I'm adding better comments in this one
*******************************************************/

var (
	pubsubContext context.Context
	topicName     = "api_v3.1_analytics"
	ProjectID     = "curt-groups"
	ClientKey     = os.Getenv("CLIENT_KEY")
	OAuthEmail    = os.Getenv("OAUTH_EMAIL")
)

// Header A key-value store for structuring header information. We need
// this since json.Marshal can't handle map[string][]string "net/http Header".
type Header struct {
	Key   string   `bson:"key" json:"key" xml:"key"`
	Value []string `bson:"value" json:"value" xml:"value"`
}

// RequestMetrics Holds data surrounding the incoming request.
type RequestMetrics struct {
	IP          string    `bson:"ip" json:"ip" xml:"ip"`
	ContentType string    `bson:"content_type" json:"content_type" xml:"content_type"`
	Body        []byte    `bson:"body" json:"body" xml:"body"`
	URI         string    `bson:"uri" json:"uri" xml:"uri"`
	Title       string    `bson:"title" json:"title" xml:"title"`
	Method      string    `bson:"method" json:"method" xml:"method"`
	Headers     []Header  `bson:"headers" json:"headers" xml:"headers"`
	Timestamp   time.Time `bson:"timestamp" json:"timestamp" xml:"timestamp"`
}

// ResponseMetrics Holds data surround the outgoing response.
type ResponseMetrics struct {
	ContentType string    `bson:"content_type" json:"content_type" xml:"content_type"`
	StatusCode  int       `bson:"status_code" json:"status_code" xml:"status_code"`
	Headers     []Header  `bson:"headers" json:"headers" xml:"headers"`
	Timestamp   time.Time `bson:"timestamp" json:"timestamp" xml:"timestamp"`
}

// Metrics Holds relevant information that will help report out request
// analytics.
type Metrics struct {
	Application      string          `json:"application" bson:"application"`
	Machine          string          `bson:"machine" json:"machine" xml:"machine"`
	Request          RequestMetrics  `bson:"request_metrics" json:"request_metrics" xml:"request_metrics"`
	Response         ResponseMetrics `bson:"response_metrics" json:"response_metrics" xml:"response_metrics"`
	Latency          int64           `bson:"latency" json:"latency" xml:"latency"`
	Body             []byte          `bson:"body" json:"body" xml:"body"`
	User             string          `bson:"user" json:"user" xml:"user"`
	AnalyticsAccount string          `bson:"analytics_account" json:"analytics_account" xml:"analytics_account"`
}

// ToPubSub is a function meant to break down the request sent to us
// and take that information and send to Google Analytics (or more accurately,
// send it to curt-consumer which sends it to analytics for us)
func ToPubSub(w http.ResponseWriter, r *http.Request, startTime time.Time) error {
	//body, _ := ioutil.ReadAll(r.Body)
	//Reading the body here causes it to be unaccessible in the actual request processing

	var reqHeaders []Header
	for k, v := range r.Header {
		reqHeaders = append(reqHeaders, Header{
			Key:   k,
			Value: v,
		})
	}

	var respHeaders []Header
	for k, v := range w.Header() {
		respHeaders = append(respHeaders, Header{
			Key:   k,
			Value: v,
		})
	}

	//This is attempting to generate a string which will be unique for
	//every user. It adds together their user-agent, their reported IP,
	//and the X-Forwarded-For and X-Real-IP headers if they exist
	user := r.UserAgent() + r.RemoteAddr + r.Header.Get("X-Forwarded-For") +
		r.Header.Get("X-Real-IP")

	reqMetrics := RequestMetrics{
		IP:          r.RemoteAddr,
		ContentType: r.Header.Get("Content-Type"),
		//Body:        body,
		URI:       r.URL.String(),
		Title:     r.URL.Path,
		Method:    r.Method,
		Headers:   reqHeaders,
		Timestamp: startTime,
	}

	respMetrics := ResponseMetrics{
		ContentType: w.Header().Get("Content-Type"),
		Headers:     respHeaders,
		Timestamp:   time.Now(),
	}

	data := Metrics{
		AnalyticsAccount: os.Getenv("GOAPI_GA_ACCOUNT"),
		Application:      "apiv2.2",
		Request:          reqMetrics,
		Response:         respMetrics,
		Latency:          time.Since(startTime).Nanoseconds(),
	}

	//This generates a v4 UUID that can be send to GA to uniquely identify the user
	userUUID, err := uuid.FromString(genUUID(base64.StdEncoding.EncodeToString([]byte(user))))
	if err != nil {
		return err
	}
	data.User = userUUID.String()

	//Here we are putting all of this data into a message to send to pubsub
	var msg pubsub.Message
	msg.Data, err = json.Marshal(&data)
	if err != nil {
		return err
	}

	//Creating the client and topic so we can send
	client, err := createClient()
	if err != nil {
		return err
	}

	defer client.Close()

	topic := client.Topic(topicName)

	var exists bool
	exists, err = topic.Exists(pubsubContext)
	if err != nil {
		log.Println(err)
		return err
	}
	if !exists {
		topic, err = client.CreateTopic(pubsubContext, topicName)
		if err != nil {
			return err
		}
	}

	defer topic.Stop()

	//Finally, publishing our message
	res := topic.Publish(pubsubContext, &msg)
	_, err = res.Get(pubsubContext)
	return err
}

//createClient creates a pubsub client for single time use, which should be
//closed by the function who calls this function
func createClient() (*pubsub.Client, error) {
	pubsubContext = context.Background()
	conf := &jwt.Config{
		Email:      OAuthEmail,
		PrivateKey: []byte(ClientKey),
		Scopes: []string{
			pubsub.ScopePubSub,
			pubsub.ScopeCloudPlatform,
		},
		TokenURL: google.JWTTokenURL,
	}

	ts := conf.TokenSource(pubsubContext)

	var err error
	client, err := pubsub.NewClient(pubsubContext, ProjectID, option.WithTokenSource(ts))

	return client, err
}

//genUUID generates a v4 UUID to give to Google Analytics based off of the
//user agent and IP address
//xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx Format
func genUUID(userData string) string {
	num := binary.BigEndian.Uint64([]byte(userData))
	s1 := rand.NewSource(int64(num))
	r1 := rand.New(s1)

	uuidArr := make([]byte, 16)
	r1.Read(uuidArr)

	//UUID v4 has specific bits in the byte array that must be set in order
	//to be a v4 UUID. Google Analytics looks for these bits and won't accept
	//UUIDs without them
	uuidArr[8] = uuidArr[8]&^0xc0 | 0x80
	uuidArr[6] = uuidArr[6]&^0xf0 | 0x40

	uuid := fmt.Sprintf("%x-%x-%x-%x-%x", uuidArr[0:4], uuidArr[4:6], uuidArr[6:8], uuidArr[8:10], uuidArr[10:])

	return uuid
}
