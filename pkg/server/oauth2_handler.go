package server

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"os"
	"time"

	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
)

type oAuth2RedirectHandler struct{}

func (handler oAuth2RedirectHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Recieved a request, using version: v0.1.9")
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

	body := json.NewDecoder(response.Body)
	if response.StatusCode != 200 {
		fmt.Printf("Recieved a non-success based status code from Discord!\nStatus Code Received: %+v\n%+v\n", response.StatusCode, response)
		fmt.Println("Error Response Body:")
		var errorBody DiscordErrorResponse
		error = body.Decode(&errorBody)
		if error == nil {
			fmt.Printf("%+v\n", errorBody)
		} else {
			fmt.Println("Couldn't decode the error body!")
			fmt.Println(error)
		}
		return
	}

	// Extract the token info from the successful response
	var responseBody AuthResponseBody
	body.Decode(&responseBody)
	if error != nil {
		fmt.Println("There was a problem parsing the json response body from Discord", error)
		return
	}

	fmt.Println("Caching the token information")
	// Write token info to Data Store
	expiresAt := time.Duration(time.Second.Seconds() * float64(responseBody.ExpiresAt))
	payload, error := json.Marshal(&TokenInformation{
		AccessToken:  responseBody.AccessToken,
		RefreshToken: responseBody.RefreshToken,
		ExpiresAt:    time.Now().Add(expiresAt),
	})
	if error != nil {
		fmt.Println("There was a problem serializing the token information payload", error)
		return
	}
	fmt.Println("Payload:")
	fmt.Println(string(payload))

	session, error := session.NewSession()
	if error != nil {
		fmt.Printf("Unable to create new AWS Session\n%+v\n", error)
		return
	}
	dynamoDBClient := dynamodb.New(session)

	var accountIDAttributeValue, tokenInformationAttributeValue *dynamodb.AttributeValue
	accountIDAttributeValue = accountIDAttributeValue.SetS(state[0])
	tokenInformationAttributeValue = tokenInformationAttributeValue.SetS(string(payload))

	tableName := "discord-token"

	input := dynamodb.PutItemInput{
		Item:      map[string]*dynamodb.AttributeValue{"accountID": accountIDAttributeValue, "tokenInfo": tokenInformationAttributeValue},
		TableName: &tableName,
	}
	request, _ := dynamoDBClient.PutItemRequest(&input)
	error = request.Send()
	if error != nil {
		fmt.Println("There was a problem writing to Dynamo DB", error)
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
