package server

import (
	"../mbcv"
	mbcvr "../mbcv/requests"
	"../requests"
	"context"
	"encoding/json"
	"fmt"
	"github.com/pkg/errors"
	"github.com/soundhound/houndify-sdk-go"
	"io/ioutil"
	"log"
	"net/http"
)

func CreateRoutes(mux *http.ServeMux, houndClient *houndify.Client, mbcvClient *mbcv.AuthenticatedClient, handler IntentHandler) {
	mux.HandleFunc("/chat", func(writer http.ResponseWriter, request *http.Request) {
		response := requests.ChatRespnse{Text: "There was a problem while processing your request"}
		defer func() {
			response, err := json.Marshal(response)
			if err != nil {
				log.Print(fmt.Errorf("failed to serialize resonse: %v", err))
				return
			}
			writer.Write(response)
		}()
		authToken, vehicleID, err := extractHeaders(request)
		if err != nil {
			log.Print("missing auth token or vehicle ID")
			return
		}
		ctx := makeContext(request.Context(), authToken, vehicleID)
		if request.Method != http.MethodPost {
			log.Print("attempt to do something other than POST to /chat")
			return
		}

		body, err := ioutil.ReadAll(request.Body)
		if err != nil {
			log.Print(fmt.Errorf("failed to read request data: %v", err))
			return
		}

		chatRequest := requests.ChatRequest{}
		err = json.Unmarshal(body, &chatRequest)
		if err != nil {
			log.Print(fmt.Errorf("could not parse request json: %v", err))
			return
		}
		houndResponseString, err := houndClient.TextSearch(houndify.TextRequest{
			Query: chatRequest.Text,
		})
		if err != nil {
			log.Print(fmt.Errorf("could not execute text search: %v", err))
			return
		}
		houndResponse, err := requests.ParseHoundResponse(houndResponseString)
		if err != nil {
			log.Print(fmt.Errorf("could not parse response from Hound: %v", err))
			return
		}

		if houndResponse.Intent == "" {
			response = requests.ChatRespnse{Text: houndResponse.WrittenResponse}
			return
		}

		commandResponse, err := runCommandForIntent(houndResponse.Intent, handler, mbcvClient, ctx)
		if err != nil {
			log.Print(fmt.Errorf("failed to execute command for intent: %v", err))
			response = requests.ChatRespnse{Text: "There was a problem executing your request on your vehicle"}
			return
		}
		if commandResponse != "" {
			response = requests.ChatRespnse{Text: commandResponse}
		} else {
			response = requests.ChatRespnse{Text: houndResponse.WrittenResponse}
		}
	})
}

func extractHeaders(request *http.Request) (string, string, error) {
	authToken := request.Header.Get(requests.MBCVAuthTokenKey)
	vehicleID := request.Header.Get(requests.MBCVVehicleIDKey)
	if authToken == "" || vehicleID == "" {
		return "", "", errors.New("missing vehicle ID or MBCV auth token")
	}
	return authToken, vehicleID, nil
}

func runCommandForIntent(intent string, handler IntentHandler, mbcvClient *mbcv.AuthenticatedClient, ctx context.Context) (string, error) {
	command := handler.commandForIntent(intent)
	if command == nil {
		return "", fmt.Errorf("could not determine command for intent: %v", intent)
	}
	switch *command {
	case unlockCar:
		_, err := mbcvClient.SendDoorCommand(mbcvr.UNLOCK, vehicleID(ctx), authToken(ctx))
		return "", err
	case lockCar:
		_, err := mbcvClient.SendDoorCommand(mbcvr.LOCK, vehicleID(ctx), authToken(ctx))
		return "", err
	case findCar:
		location, err := mbcvClient.GetLocation(vehicleID(ctx), authToken(ctx))
		if err != nil {
			return "", err
		}
		return fmt.Sprintf("Located your vehicle at %v,%v", location.Longitude.Value, location.Latitude.Value), nil
	default:
		return "", fmt.Errorf("unable to handle command: %v", command)
	}
}


