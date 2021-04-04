package mbcv

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
)

const (
	Lock = "LOCK"
	Unlock = "UNLOCK"
)

const (
	apiURL = "https://api.mercedes-benz.com/experimental/connectedvehicle/v2/vehicles/"
)

type ServerClient struct {
	ClientID string
	ClientSecret string
	httpClient *http.Client
}

func (c *ServerClient) SendDoorCommand(command, vehicleID, mbToken string) (*CommandResponse, error) {
	commandReq := commandRequest{Command: command}
	body, err := json.Marshal(commandReq)
	if err != nil {
		return nil, err
	}

	log.Print(string(body))

	req, err := http.NewRequest(http.MethodPost, apiURL + vehicleID + "/doors", bytes.NewBuffer(body))
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("authorization", "Bearer " + mbToken)

	if c.httpClient == nil {
		c.httpClient = &http.Client{}
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, errors.New("failed to successfully run request: " + err.Error())
	}

	body, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, errors.New("failed to read body: " + err.Error())
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("error response: %v", string(body))
	}

	cmdResponse := CommandResponse{}
	err = json.Unmarshal(body, &cmdResponse)
	if err != nil {
		return nil, fmt.Errorf("failed to parse command response: %v", string(body))
	}
	return &cmdResponse, nil
}
