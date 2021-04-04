package client

import (
	"../requests"
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

type Client struct {
	ServerURL string
	authToken string
	refreshToken string
	httpClient *http.Client
}

func (c *Client) SendMessage(message, vehicleID string) (*requests.ChatRespnse, error) {
	body, err := json.Marshal(requests.ChatRequest{Text: message})
	if err != nil {
		return nil, fmt.Errorf("failed to encode message: %v", err)
	}
	req, err := http.NewRequest(http.MethodPost, c.ServerURL, bytes.NewBuffer(body))
	if err != nil {
		return nil, fmt.Errorf("failed to create a new request: %v", err)
	}

	req.Header.Set(requests.MBCVVehicleIDKey, vehicleID)
	req.Header.Set(requests.MBCVAuthTokenKey, c.authToken)

	if c.httpClient == nil {
		c.httpClient = &http.Client{}
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to execute http request: %v", err)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("request failed: %v", resp.StatusCode)
	}

	body, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read request body: %v", err)
	}

	chatResponse := requests.ChatRespnse{}
	err = json.Unmarshal(body, &chatResponse)
	if err != nil {
		return nil, fmt.Errorf("failed to parse json from request body: %v", err)
	}
	return &chatResponse, nil
}