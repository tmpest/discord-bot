package server

import "time"

// TokenRequestBody is the OAuth2 request body recieved from Discord via the redirect URI
type TokenRequestBody struct {
	ClientID     string `json:"client_id"`
	ClientSecret string `json:"client_secret"`
	GrantType    string `json:"grant_type"`
	Code         string `json:"code"`
	RedirectURI  string `json:"redirect_uri"`
	Scope        string `json:"scope"`
}

// AuthResponseBody this is the body of the response we get back from Discord
type AuthResponseBody struct {
	AccessToken  string `json:"access_token"`
	ExpiresAt    int32  `json:"expires_in"`
	RefreshToken string `json:"refresh_token"`
	Scope        string `json:"scope"`
	TokenType    string `json:"token_type"`
}

// TokenInformation is the cached information
type TokenInformation struct {
	AccessToken  string    `json:"access_token"`
	ExpiresAt    time.Time `json:"expires_in"`
	RefreshToken string    `json:"refresh_token"`
}
