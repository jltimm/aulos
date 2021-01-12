package auth

import (
	"context"

	"../secrets"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/clientcredentials"
)

// GetConfig creates and returns the config
func GetConfig() *clientcredentials.Config {
	apiURL := "https://accounts.spotify.com/api/token"
	config := &clientcredentials.Config{
		ClientID:     secrets.GetClientID(),
		ClientSecret: secrets.GetClientSecret(),
		TokenURL:     apiURL,
	}
	return config
}

// GetToken returns the token. Should not see expiration errors
func GetToken() *oauth2.Token {
	token, err := GetConfig().Token(context.Background())
	if err != nil {
		panic(err)
	}
	return token
}
