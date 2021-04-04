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

	//c := client.Client{
	//	ServerURL:    "http://localhost",
	//}

	mbcvClient := mbcv.Client{
		ClientID:     info.ClientID,
		ClientSecret: info.ClientSecret,
	}
	err = mbcvClient.Authenticate()
	if err != nil {
		panic(fmt.Errorf("failed to authenticate with mbcv: %v", err))
	}
}
