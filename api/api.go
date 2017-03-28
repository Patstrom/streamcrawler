// Package api defines an interface for, and the different implementations of, apis (e.g twitch)
package api

import (
	"encoding/json"
	"net/http"
	"net/url"
	"time"
)

// API is the interface that determines what the actual implementations should
// support
type API interface {
	Streams(game, limit string) Streams
}

// Pair defines a key: value string pair. Mainly used as parameters to the apis
type Pair struct {
	Key, Val string
}

// Streams defines what we expect the Streams() function to return. Each API has to
// convert the actual response into something that fits into this struct. This struct
// should ideally correspond to what we want in the database.
type Streams struct {
	Error        bool
	ErrorMessage string
	Total        float64
	Streams      []Stream
}

// Stream defines information about a stream
type Stream struct {
	Id              float64
	Game            string
	PlatformId      float64
	PreviewImageUrl string
	LogoImageUrl    string
	ExternalUserId  float64
	Viewers         float64
	ChannelId       float64
	UserName        string
	DisplayName     string
	StatusText      string
}

// getJson unmarshals the response from u into target.
// Times out from the GET request after 10 seconds.
func getJson(u *url.URL, auth *Pair, target interface{}) error {
	client := &http.Client{Timeout: 10 * time.Second}
	req := &http.Request{
		Method: "GET",
		URL:    u,
	}
	if len(auth.Key) > 0 && len(auth.Val) > 0 {
		req.Header = http.Header{auth.Key: {auth.Val}}
	}

	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	return json.NewDecoder(resp.Body).Decode(target)
}

// Does a GET request to the URL with the given params and unmarshals json into target
func ApiCall(target interface{}, baseUrl string, auth *Pair, param ...Pair) {
	u, err := url.Parse(baseUrl)
	if err != nil {
		//TODO: Handle error like an adult
	}
	q := u.Query()

	for _, v := range param {
		q.Set(v.Key, v.Val)
	}
	u.RawQuery = q.Encode()

	err = getJson(u, auth, target)
}
