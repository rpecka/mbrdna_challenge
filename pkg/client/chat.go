package client

import (
	"fmt"
	"github.com/AlecAivazis/survey"
)

const (
	help = "help"
	quit = "quit"
)

func StartChat(vechicleID string, client *Client) {
	fmt.Println("You are now connected to the bot...")
	showHelp()
	for {
		var text string
		err := survey.AskOne(&survey.Input{}, &text)
		if err != nil {
			break
		}
		switch text {
		case "":
			continue
		case help:
			showHelp()
			continue
		case quit:
			return
		}

		fmt.Println()
		response, err := client.SendMessage(text, vechicleID)
		if err != nil {
			fmt.Printf("Error: %v\n\n", err)
			continue
		}
		fmt.Printf("%v\n\n", response.Text)
	}
}

func showHelp() {
	fmt.Println("Usage:\n\tType a message and press return to chat with the bot\n\tType `quit` or CTRL-C to exit" +
		"\n\tType `help` to show this message")
}
