package auth

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"

	pathutils "github.com/agungfir98/gcal-cli/utils/path_utils"
	"golang.org/x/oauth2"
)

func GetClient(config *oauth2.Config) *http.Client {
	tokenFile := pathutils.GetTokenFile()
	token, err := TokenFromFile(tokenFile)
	if err != nil {
		token = GetTokenFromWeb(config)
		fmt.Printf("\rSaving token file to: %s\n", tokenFile)
		saveToken(tokenFile, token)
	}

	tokenSource := config.TokenSource(context.Background(), token)
	refreshedToken, err := tokenSource.Token()
	if err != nil {
		fmt.Printf("unable to refresh the token: %v\n", err)
		token = GetTokenFromWeb(config)
		fmt.Printf("\rSaving token file to: %s\n", tokenFile)
		saveToken(tokenFile, token)
	}

	if refreshedToken.AccessToken != token.AccessToken {
		fmt.Printf("\rToken refreshed, saving token file to: %s\n", tokenFile)
		saveToken(tokenFile, refreshedToken)
	}

	return config.Client(context.Background(), token)
}

func TokenFromFile(file string) (*oauth2.Token, error) {
	f, err := os.Open(file)
	if err != nil {
		return nil, err
	}

	defer f.Close()

	t := &oauth2.Token{}

	err = json.NewDecoder(f).Decode(t)
	return t, err
}

func GetTokenFromWeb(cfg *oauth2.Config) *oauth2.Token {
	authUrl := cfg.AuthCodeURL("state-token", oauth2.AccessTypeOffline)

	fmt.Printf("visit the following to link to obtain the code:\n %v\n", authUrl)

	var authCode string

	fmt.Print("paste the code in the url here: ")
	if _, err := fmt.Scan(&authCode); err != nil {
		log.Fatalf("Unable to read authorization code: %v", err)
	}

	tok, err := cfg.Exchange(context.TODO(), authCode)

	if err != nil {
		log.Fatalf("unable to retreive token from web: %v\n", err)
	}

	return tok
}

func saveToken(path string, token *oauth2.Token) {
	f, err := os.Create(path)
	if err != nil {
		log.Fatalf("Unable to cache oauth token: %v", err)
	}
	defer f.Close()
	json.NewEncoder(f).Encode(token)
}

func GetCredential() []byte {
	path := pathutils.GetCredentialsFile()
	b, err := os.ReadFile(path)
	if err != nil {
		log.Fatalf("unable to read credentials file: %v\n", err)
	}
	return b
}
