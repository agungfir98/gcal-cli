package auth

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"sync"
	"time"

	"github.com/agungfir98/gcal-cli/server"
	"github.com/agungfir98/gcal-cli/utils/browser"
	pathutils "github.com/agungfir98/gcal-cli/utils/path_utils"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/calendar/v3"
)

type Client struct {
	codeCh    chan string
	config    *oauth2.Config
	token     *oauth2.Token
	tokenFile string
}

func NewClient() *Client {
	cred := GetCredential()
	cfg, err := google.ConfigFromJSON(cred, calendar.CalendarReadonlyScope)
	if err != nil {
		log.Fatalf("unable to parse client secret file to config: %v\n", err)
	}

	c := &Client{}
	c.codeCh = make(chan string)
	c.config = cfg
	c.config.RedirectURL = "http://localhost:8080/callback"
	c.tokenFile = pathutils.GetTokenFile()

	return c
}

func (c *Client) GetClient() *http.Client {
	token, err := c.tokenFromFile()
	if err != nil {
		c.getTokenFromWeb()
		c.saveToken()
		return c.client()
	}

	c.token = token
	err = c.refreshToken()
	if err != nil {
		c.getTokenFromWeb()
		c.saveToken()
		return c.client()
	}

	return c.client()
}

func (c *Client) tokenFromFile() (*oauth2.Token, error) {
	f, err := os.Open(c.tokenFile)
	if err != nil {
		return nil, fmt.Errorf("error occured: %v\n", err)
	}

	defer f.Close()
	token := &oauth2.Token{}
	err = json.NewDecoder(f).Decode(token)
	if err != nil {
		return nil, fmt.Errorf("error decoding token: %v\n", err)
	}

	return token, nil
}

func (c *Client) getTokenFromWeb() {
	authUrl := c.config.AuthCodeURL("state-token", oauth2.AccessTypeOffline)
	var code string

	browser, err := browser.GetBrowserOpener()
	if err != nil {
		log.Fatalln(err)
	}

	err = browser.Open(authUrl)
	if err != nil {
		log.Printf("unable to open the browser, you may do it manually.\n visit: %v\n", authUrl)

		go func() {
			fmt.Print("find 'code' query in the url and paste the code here: ")
			if _, err := fmt.Scan(&code); err != nil {
				log.Fatalf("Unable to read authorization code: %v", err)
			}
			c.codeCh <- code
		}()

		select {
		case <-c.codeCh:
			fmt.Println("converting auth code into token")
		case <-time.After(3 * time.Minute):
			fmt.Println("\nTimeout no input received")
			os.Exit(1)
		}

		token, err := c.config.Exchange(context.TODO(), code)

		if err != nil {
			log.Fatalf("unable to retreive token from web: %v\n", err)
		}

		c.token = token
		return
	}

	wg := sync.WaitGroup{}
	wg.Add(1)

	server := server.NewServer(":8080", c.codeCh)

	server.Start(&wg)

	code = <-c.codeCh

	if err := server.Shutdown(); err != nil {
		fmt.Println("Error shutting down server:", err)
	}

	wg.Wait()

	token, err := c.config.Exchange(context.Background(), code)
	if err != nil {
		log.Fatalf("unable to retreive token from web: %v\n", err)
	}

	c.token = token
	return
}

func (c *Client) saveToken() {
	fmt.Printf("\rSaving token file to: %s\n", c.tokenFile)

	if c.token == nil {
		log.Fatalf("no token")
	}

	f, err := os.Create(c.tokenFile)
	if err != nil {
		log.Fatalf("Unable to cache oauth token: %v", err)
	}
	defer f.Close()
	err = json.NewEncoder(f).Encode(c.token)
	if err != nil {
		log.Fatalf("failed to encode token: %v\n", err)
	}
}

func (c *Client) client() *http.Client {
	return c.config.Client(context.Background(), c.token)
}

func (c *Client) refreshToken() error {
	tokenSource := c.config.TokenSource(context.Background(), c.token)
	token, err := tokenSource.Token()
	if err != nil {
		return fmt.Errorf("unable to refresh token: %v", err)
	}
	c.token = token

	return nil

}

func GetCredential() []byte {
	path := pathutils.GetCredentialsFile()
	b, err := os.ReadFile(path)
	if err != nil {
		log.Fatalf("unable to read credentials file: %v\n", err)
	}
	return b
}
