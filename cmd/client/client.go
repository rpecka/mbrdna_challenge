package main

import (
	"../../pkg/client"
	"../../pkg/mbcv"
	"fmt"
)

func main() {
	info, err := client.RunStartupSurvey()
	if err != nil {
		panic(err)
	}

	mbcvClient := mbcv.Client{
		ClientID:     info.ClientID,
		ClientSecret: info.ClientSecret,
	}
	authToken, err := mbcvClient.Authenticate()
	if err != nil {
		panic(fmt.Errorf("failed to authenticate with mbcv: %v", err))
	}

	vehicles, err := mbcvClient.GetVehicles(authToken)
	if err != nil {
		panic(fmt.Errorf("failed to get the user's vehicles: %v", err))
	}

	if len(*vehicles) < 1 {
		panic("the user does not have any vehicles")
	}

	fmt.Printf("Using first found vehicle: %v\n", (*vehicles)[0].ID)

	c := &client.Client{
		ServerURL:    "http://localhost/",
		AuthToken: authToken,
	}

	client.StartChat((*vehicles)[0].ID, c)
}
