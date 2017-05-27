package auth

import (
	"fmt"
	"net/http"
	"net/url"

	"github.com/garyburd/go-oauth/oauth"
	"github.com/skratchdot/open-golang/open"
)

// Client :
type Client struct {
	*oauth.Client
}

// NewClient :
func NewClient(token, secret string) *Client {
	return &Client{
		Client: &oauth.Client{
			Credentials:                   oauth.Credentials{Token: token, Secret: secret},
			TemporaryCredentialRequestURI: "https://www.hatena.com/oauth/initiate",
			ResourceOwnerAuthorizationURI: "https://www.hatena.ne.jp/oauth/authorize",
			TokenRequestURI:               "https://www.hatena.com/oauth/token",
		},
	}
}

// Auth :
func (c *Client) Auth(client *http.Client) (*oauth.Credentials, error) {
	u := url.Values{}
	u.Add("scope", "read_private,write_private")

	tempCred, err := c.RequestTemporaryCredentials(client, "oob", u)
	if err != nil {
		return nil, err
	}

	err = open.Start(c.AuthorizationURL(tempCred, nil))
	if err != nil {
		return nil, err
	}

	fmt.Printf("1. Go to %s\n2. Authorize the application\n3. Enter verification code:\n", "*")
	var code string
	for {
		fmt.Scanln(&code)
		if code != "" {
			break
		}
	}

	tokenCred, _, err := c.RequestToken(client, tempCred, code)
	if err != nil {
		return nil, err
	}
	return tokenCred, nil
}
