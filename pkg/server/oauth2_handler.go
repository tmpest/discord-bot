package server

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"time"

	"github.com/bradfitz/gomemcache/memcache"
)

type oAuth2RedirectHandler struct{}

func (handler oAuth2RedirectHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Recieved a request")
	// Capture the code and state from the request
	authCode, ok := getQueryParamFromRequest("code", r)
	if !ok {
		return
	}
	state, ok := getQueryParamFromRequest("state", r)
	if !ok {
		return
	}

	fmt.Println("Making a request to Discord")
	response, error := http.PostForm(discordOAuth2TokenEndpoint, url.Values{
		"client_id":     {os.Getenv("TMPEST_BOT_CLIENT_ID")},
		"client_secret": {os.Getenv("TMPEST_BOT_CLIENT_SECRET")},
		"grant_type":    {"authorization_code"},
		"code":          {authCode[0]},
		"redirect_uri":  {redirectURI},
		"scope":         {"connections"},
	})
	if error != nil {
		fmt.Println("There was a problem making the request to Discord for the Token", error)
		return
	}

	defer response.Body.Close()

	if response.StatusCode != 200 {
		fmt.Printf("Recieved a non-success based status code from Discord!\nStatus Code Received: %+v\n%+v\n", response.StatusCode, response)
		fmt.Println("Error Response Body:")
		body, error := ioutil.ReadAll(response.Body)
		if error == nil {
			fmt.Printf("%+v\n", string(body))
		}
		return
	}

	// Extract the token info from the successful response
	responseBody := &AuthResponseBody{}
	responseBodyJSON := make([]byte, 0)
	response.Body.Read(responseBodyJSON)
	error = json.Unmarshal(responseBodyJSON, responseBody)
	if error != nil {
		fmt.Println("There was a problem parsing the json response body from Discord", error)
		return
	}

	fmt.Println("Caching the token information")
	// Write token info to Redis
	cache := memcache.New("memcached-11217.c10.us-east-1-2.ec2.cloud.redislabs.com:11217")
	expiresAt := time.Duration(time.Second.Seconds() * float64(responseBody.ExpiresAt))
	cachePayload, error := json.Marshal(&TokenInformation{
		AccessToken:  responseBody.AccessToken,
		RefreshToken: responseBody.RefreshToken,
		ExpiresAt:    time.Now().Add(expiresAt),
	})
	if error != nil {
		fmt.Println("There was a problem serializing the cache payload", error)
		return
	}
	item := &memcache.Item{Key: state[0], Value: cachePayload}
	error = cache.Set(item)
	if error != nil {
		fmt.Println("There was a problem setting the Redis Cache", error)
		return
	}
	fmt.Printf("Success! Cached token information for account: %+v\n", state[0])
}

func getQueryParamFromRequest(paramName string, r *http.Request) ([]string, bool) {
	value, ok := r.URL.Query()[paramName]
	if !ok {
		fmt.Printf("Invalid Request Received! '%+v' query param is required!\n", paramName)
	}
	return value, ok
}
