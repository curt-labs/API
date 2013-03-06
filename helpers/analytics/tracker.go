package analytics

const (
	version = "1.0"
)

type Tracker struct {
	Config     *Config
	AccountId  string
	DomainName string
}

// func NewTracker(acctId string, domain string, conf *Config) (Tracker, error) {
// 	track := Tracker{
// 		Config:     conf,
// 		AccountId:  acctId,
// 		DomainName: domain,
// 	}

// 	return track, nil
// }

// func (t *Tracker) TrackPageView(page *Page, sess *Session, visitor *Visitor) error {
// 	req := PageViewRequest{page, sess, visitor, t, conf}
// 	req.Fire()

// }
