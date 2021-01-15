package auth

import (
	"context"

	"../secrets"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/clientcredentials"
)

var config *clientcredentials.Config = &clientcredentials.Config{
	ClientID:     secrets.GetClientID(),
	ClientSecret: secrets.GetClientSecret(),
	TokenURL:     "https://accounts.spotify.com/api/token",
}

// GetConfig returns the config
func GetConfig() *clientcredentials.Config {
	return config
}

// GetToken returns the token. Should not see expiration errors
func GetToken() *oauth2.Token {
	token, err := config.Token(context.Background())
	if err != nil {
		panic(err)
	}
	return token
}
