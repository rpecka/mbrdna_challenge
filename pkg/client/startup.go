package client

import (
	"fmt"
	"github.com/AlecAivazis/survey"
)

type startupInfo struct {
	ClientID string
	ClientSecret string
}

func RunStartupSurvey() (*startupInfo, error) {
	var answerIdx int
	err := survey.AskOne(&survey.Select{
		Message: "Which MBCV client credentials would you like to use?",
		Options: []string{"Use Russell's", "Use your own"},
	}, &answerIdx, survey.WithValidator(survey.Required))
	if err != nil {
		return nil, err
	}

	if answerIdx == 0 {
		fmt.Println("Using Russell's credentials...")
		return &startupInfo{
			ClientID: "64ce6056-79b9-40ab-80e0-71d3c805c575",
			ClientSecret: "kjuqQzTqgqpnYgmEblpnbEuyfwErKGZqDxKkQjrCPdXlQggvnkgFYKXNRwtaHPLy",
		}, nil
	}
	
	var info startupInfo
	err = survey.Ask([]*survey.Question{
		{
			Name: "ClientID",
			Prompt: &survey.Input{
				Message: "Enter your MBCV client ID",
			},
				},
		{
			Name: "ClientSecret",
			Prompt: &survey.Input{
				Message: "Enter your MBCV client secret",
			},
				},
	}, &info, survey.WithValidator(survey.Required))
	return &info, err
}
