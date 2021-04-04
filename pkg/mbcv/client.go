package mbcv

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"

	"github.com/AlecAivazis/survey"
	"github.com/rpecka/mbrdna_challenge/pkg/mbcv/requests"
)

const (
	oAuthURL = "https://id.mercedes-benz.com/as/"
	oAuthRedirectFormat = oAuthURL + "authorization.oauth2?response_type=code&client_id=%v&redirect_uri=&scope=mb:vehicle:status:general%%20mb:user:pool:reader%%20offline_access&state=%v"
	oAuthTokenURL = oAuthURL + "token.oauth2"
)

type Client struct {
	AuthenticatedClient
	ClientID string
	ClientSecret string
}

func (c *Client) Authenticate() (string, error) {
	fmt.Printf("Please open the following link and copy the URL you are redirected to:\n%v\n", c.makeOAuthURL())
	var urlString string
	err := survey.AskOne(&survey.Input{
		Message: "Please paste the callback URL that you were redirected to",
	}, &urlString, survey.WithValidator(survey.Required))
	if err != nil {
		return "", err
	}
	redirectURL, err := url.Parse(urlString)
	if err != nil {
		return "", fmt.Errorf("there was a problem parsing the redirect url: %v", err)
	}
	code := redirectURL.Query().Get("code")
	// TODO: check this to prevent XSF
	//redirectURL.Query().Get("state")
	if code == "" {
		return "", fmt.Errorf("redirect url did not contain a `code` query parameter")
	}

	tokenResponse, err := c.requestAuthToken(code)
	if err != nil {
		return "", fmt.Errorf("failed to get auth token from code: %v", err)
	}
	return tokenResponse.AccessToken, nil
}

func (c Client) makeOAuthURL() string {
	// TODO: Generate random value for state and check it later to prevent XSF
	return fmt.Sprintf(oAuthRedirectFormat, c.ClientID, "123e4567-e89b-12d3-a456-426655440000")
}

func (c Client) requestAuthToken(authCode string) (*requests.AuthTokenResponse, error) {
	data := url.Values{
		"grant_type": {"authorization_code"},
		"code": {authCode},
		"redirect_uri": {""},
	}

	req, err := http.NewRequest(http.MethodPost, oAuthTokenURL, strings.NewReader(data.Encode()))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	req.SetBasicAuth(c.ClientID, c.ClientSecret)

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, err
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("non-200 status code when requesting auth token: %v, %v", resp.StatusCode, string(body))
	}

	var tokenResponse requests.AuthTokenResponse
	err = json.Unmarshal(body, &tokenResponse)
	if err != nil {
		return nil, err
	}
	return &tokenResponse, nil
}
