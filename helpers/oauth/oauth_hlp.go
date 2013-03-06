package oauth_helper

import (
	"code.google.com/p/goauth2/oauth"
	"encoding/gob"
	"fmt"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
)

func TokenCacheFile(config *oauth.Config) string {
	return filepath.Join(osUserCacheDir(), url.QueryEscape(
		fmt.Sprintf("go-api-%s-%s-%s", config.ClientId, config.ClientSecret, config.Scope)))
}

func TokenFromFile(file string) (*oauth.Token, error) {
	if !*cacheToken {
		return nil, errors.New("--cachetoken is false")
	}
	f, err := os.Open(file)
	if err != nil {
		return nil, err
	}
	t := new(oauth.Token)
	err = gob.NewDecoder(f).Decode(t)
	return t, err
}

func TokenFromWeb(config *oauth.Config) *oauth.Token {
	ch := make(chan string)
	randState := fmt.Sprintf("st%d", time.Now())
	ts := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		if req.URL.Path == "/favicon.ico" {
			http.Error(rw, "", 404)
			return
		}
		if req.FormValue("state") != randState {
			log.Printf("State doesn't match: req = %#v", req)
			http.Error(rw, "", 500)
			return
		}
		if code := req.FormValue("code"); code != "" {
			fmt.Fprintf(rw, "<h1>Success</h1>Authorized.")
			rw.(http.Flusher).Flush()
			ch <- code
			return
		}
		log.Printf("no code")
		http.Error(rw, "", 500)
	}))
	defer ts.Close()

	config.RedirectURL = ts.URL
	authUrl := config.AuthCodeURL(randState)
	go openUrl(authUrl)
	log.Printf("Authorize this app at: %s", authUrl)
	code := <-ch
	log.Printf("Got code: %s", code)

	t := &oauth.Transport{
		Config:    config,
		Transport: condDebugTransport(http.DefaultTransport),
	}
	_, err := t.Exchange(code)
	if err != nil {
		log.Fatalf("Token exchange error: %v", err)
	}
	return t.Token
}

func openUrl(url string) {
	try := []string{"xdg-open", "google-chrome", "open"}
	for _, bin := range try {
		err := exec.Command(bin, url).Run()
		if err == nil {
			return
		}
	}
	log.Printf("Error opening URL in browser.")
}

func saveToken(file string, token *oauth.Token) {
	f, err := os.Create(file)
	if err != nil {
		log.Printf("Warning: failed to cache oauth token: %v", err)
		return
	}
	defer f.Close()
	gob.NewEncoder(f).Encode(token)
}

func condDebugTransport(rt http.RoundTripper) http.RoundTripper {
	if *debug {
		return &logTransport{rt}
	}
	return rt
}

func GetOAuthClient(config *oauth.Config) *http.Client {
	cacheFile := TokenCacheFile(config)
	token, err := TokenFromFile(cacheFile)
	if err != nil {
		token = TokenFromWeb(config)
		saveToken(cacheFile, token)
	} else {
		log.Printf("Using cached token %#v from %q", token, cacheFile)
	}

	t := &oauth.Transport{
		Token:     token,
		Config:    config,
		Transport: condDebugTransport(http.DefaultTransport),
	}
	return t.Client()
}
