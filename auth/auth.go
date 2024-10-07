package auth

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/agungfir98/gcal-cli/utils/browser"
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
		fmt.Print("\rToken refreshed")
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

	browser, err := browser.GetBrowserOpener()
	if err != nil {
		log.Fatalln(err)
	}

	err = browser.Open(authUrl)
	if err != nil {
		log.Printf("unable to open the browser, you may do it manually.\n visit: %v", authUrl)
	}

	var authCode string

	inputCh := make(chan string)

	go func() {
		fmt.Print("paste the code in the url here: ")
		if _, err := fmt.Scan(&authCode); err != nil {
			log.Fatalf("Unable to read authorization code: %v", err)
		}
		inputCh <- authCode
	}()

	select {
	case <-inputCh:
		fmt.Println("converting auth code into token")
	case <-time.After(3 * time.Minute):
		fmt.Println("\nTimeout no input received")
		os.Exit(1)
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
