package mbcv

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"

	"github.com/rpecka/mbrdna_challenge/pkg/mbcv/requests"
)

const (
	apiURL = "https://api.mercedes-benz.com/experimental/connectedvehicle/v2/vehicles/"
)

type AuthenticatedClient struct {
	httpClient http.Client
}

func (c AuthenticatedClient) GetVehicles(token string) (*requests.GetVehiclesResponse, error) {
	req, err := c.makeRequest(http.MethodGet, nil, token)
	if err != nil {
		return nil, err
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, err
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, errors.New("failed to read body: " + err.Error())
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("error response: %v", string(body))
	}

	var vehiclesResponse requests.GetVehiclesResponse
	err = json.Unmarshal(body, &vehiclesResponse)
	if err != nil {
		return nil, fmt.Errorf("failed to parse get vehicles response: %v", err)
	}
	return &vehiclesResponse, nil
}

func (c AuthenticatedClient) GetLocation(vehicleID, token string) (*requests.GetLocationResponse, error) {
	req, err := c.makeRequestURI(http.MethodGet, vehicleID+"/location", nil, token)
	if err != nil {
		return nil, err
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, err
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, errors.New("failed to read body: " + err.Error())
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("error response: %v", string(body))
	}

	var locationResponse requests.GetLocationResponse
	err = json.Unmarshal(body, &locationResponse)
	if err != nil {
		return nil, fmt.Errorf("failed to parse get location response: %v", err)
	}
	return &locationResponse, nil
}

func (c AuthenticatedClient) SendDoorCommand(command requests.DoorCommand, vehicleID string, token string) (*requests.CommandResponse, error) {
	commandReq := requests.CommandRequest{Command: string(command)}
	body, err := json.Marshal(commandReq)
	if err != nil {
		return nil, err
	}

	req, err := c.makeRequestURI(http.MethodPost, vehicleID+"/doors", bytes.NewBuffer(body), token)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "application/json")

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

	cmdResponse := requests.CommandResponse{}
	err = json.Unmarshal(body, &cmdResponse)
	if err != nil {
		return nil, fmt.Errorf("failed to parse command response: %v", string(body))
	}
	return &cmdResponse, nil
}

func (c AuthenticatedClient) makeRequestURI(method string, uri string, body io.Reader, token string) (*http.Request, error) {
	req, err := http.NewRequest(method, apiURL+uri, body)
	if err != nil {
		return nil, err
	}
	req.Header.Set("authorization", "Bearer "+token)
	return req, err
}

func (c AuthenticatedClient) makeRequest(method string, body io.Reader, token string) (*http.Request, error) {
	return c.makeRequestURI(method, "", body, token)
}
