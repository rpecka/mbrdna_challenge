package server

import (
	"../mbcv"
	"../requests"
	"context"
	"encoding/json"
	"fmt"
	"github.com/pkg/errors"
	"github.com/soundhound/houndify-sdk-go"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
)

func CreateRoutes(mux *http.ServeMux, houndClient *houndify.Client, mbcvClient *mbcv.AuthenticatedClient) {
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

		if houndResponse.Intent != "" {
			err = runCommandForIntent(houndResponse.Intent, mbcvClient, ctx)
		}
		if err != nil {
			log.Print(fmt.Errorf("failed to execute command for intent: %v", err))
			response = requests.ChatRespnse{Text: "There was a problem executing your request on your vehicle"}
			return
		}

		response = requests.ChatRespnse{Text: houndResponse.WrittenResponse}
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

func runCommandForIntent(intent string, mbcvClient *mbcv.AuthenticatedClient, ctx context.Context) error {
	components := strings.Split(intent, ".")
	if len(components) < 2 {
		return fmt.Errorf("too few components in %v", intent)
	}
	if len(components) > 2 {
		return fmt.Errorf("too many components in %v", intent)
	}
	if components[0] != "CAR" {
		return fmt.Errorf("unrecognized intent domain: %v", components[0])
	}
	_, err := mbcvClient.SendDoorCommand(strings.ToUpper(components[1]), vehicleID(ctx), authToken(ctx))
	if err != nil {
		return err
	}
	return nil
}


