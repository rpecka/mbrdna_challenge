package server

import (
	"../requests"
	"encoding/json"
	"fmt"
	"github.com/soundhound/houndify-sdk-go"
	"io/ioutil"
	"log"
	"net/http"
)

func CreateRoutes(mux *http.ServeMux, houndClient *houndify.Client) {
	mux.HandleFunc("/chat", func(writer http.ResponseWriter, request *http.Request) {
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
		houndResponse, err := houndClient.TextSearch(houndify.TextRequest{
			Query: chatRequest.Text,
		})
		if err != nil {
			log.Print(fmt.Errorf("could not execute text search: %v", err))
			writeError("failed to query hound server", http.StatusBadGateway, writer)
			return
		}
		readableResponse, err := houndify.ParseWrittenResponse(houndResponse)
		if err != nil {
			log.Print(fmt.Errorf("could not parse response from Hound: %v", err))
			writeError("failed to parse response from hound server", http.StatusBadGateway, writer)
			return
		}
		writer.Write([]byte(readableResponse))
	})
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


