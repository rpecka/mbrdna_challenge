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

func CreateRoutes(mux *http.ServeMux, houndClient *houndify.Client, mbcvClient *mbcv.ServerClient) {
	mux.HandleFunc("/chat", func(writer http.ResponseWriter, request *http.Request) {
		authToken, vehicleID, err := extractHeaders(request)
		if err != nil {
			writeError("missing auth token or vehicle ID", http.StatusBadRequest, writer)
			return
		}
		ctx := makeContext(request.Context(), authToken, vehicleID)
		if request.Method != http.MethodPost {
			writeError("only POST is supported", http.StatusBadRequest, writer)
			return
		}

		body, err := ioutil.ReadAll(request.Body)
		if err != nil {
			writeError("failed to read request data", http.StatusBadRequest, writer)
			return
		}

		chatRequest := requests.ChatRequest{}
		err = json.Unmarshal(body, &chatRequest)
		if err != nil {
			writeError("request could not be parsed", http.StatusBadRequest, writer)
			return
		}
		houndResponseString, err := houndClient.TextSearch(houndify.TextRequest{
			Query: chatRequest.Text,
		})
		if err != nil {
			log.Print(fmt.Errorf("could not execute text search: %v", err))
			writeError("failed to query hound server", http.StatusBadGateway, writer)
			return
		}
		houndResponse, err := requests.ParseHoundResponse(houndResponseString)
		if err != nil {
			log.Print(fmt.Errorf("could not parse response from Hound: %v", err))
			writeError("failed to parse response from hound server", http.StatusBadGateway, writer)
			return
		}

		err = runCommandForIntent(houndResponse.Intent, mbcvClient, ctx)
		if err != nil {
			log.Print(fmt.Errorf("failed to execute command for intent: %v", err))
			writeError("failed to execute connected vehicle command", http.StatusBadGateway, writer)
			return
		}

		writer.Write([]byte(houndResponse.WrittenResponse))
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

func runCommandForIntent(intent string, mbcvClient *mbcv.ServerClient, ctx context.Context) error {
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

func writeError(message string, status int, writer http.ResponseWriter) {
	writer.WriteHeader(status)
	response, err := json.Marshal(requests.ErrorResponse{Error: message})
	if err != nil {
		log.Print(fmt.Errorf("failed to serialize resonse: %v", err))
		return
	}
	writer.Write(response)
}


