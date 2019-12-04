package server

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
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
	// Build request for token information from Discord
	body, error := json.Marshal(TokenRequestBody{
		ClientID:     os.Getenv("TMPEST_BOT_CLIENT_ID"),
		ClientSecret: os.Getenv("TMPEST_BOT_CLIENT_SECRET"),
		GrantType:    "authorization_code",
		Code:         authCode[0],
		RedirectURI:  redirectURI,
		Scope:        "connections",
	})

	if error != nil {
		fmt.Println("There was a problem parsing the token info!", error)
		return
	}
	fmt.Printf("Resquest Body to Discord:\n%+v\n", body)

	tokenRequest, error := http.NewRequest(http.MethodPost, discordOAuth2TokenEndpoint, bytes.NewReader(body))
	if error != nil {
		fmt.Println("There was a problem creating the request", error)
		return
	}
	tokenRequest.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	fmt.Printf("Resquest to Discord:\n%+v\n", tokenRequest)

	// Make a client and make the request to discord for a token
	client := http.Client{}
	response, error := client.Do(tokenRequest)
	if error != nil {
		fmt.Println("There was a problem making the request to Discord for the Token", error)
		return
	}
	if response.StatusCode != 200 {
		fmt.Printf("Recieved a non-success based status code from Discord!\nStatus Code Received: %+v\n%+v\n", response.StatusCode, response)
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
