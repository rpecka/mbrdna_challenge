package main

import (
	"log"
	"net/http"

	"github.com/soundhound/houndify-sdk-go"

	"github.com/rpecka/mbrdna_challenge/pkg/mbcv"
	"github.com/rpecka/mbrdna_challenge/pkg/server"
)

func main() {
	houndClient := &houndify.Client{
		ClientID:          "aX1HeFqKSZoNJ0Wh-kZ6tA==",
		ClientKey:         "sCIOCsrBADASfXWEY_6-xiI4qm4VJhNf77LDH4mPqzCm4w3QZacJo4-mHZhNYGOTdbXadttX7TklNl2NGs9WhQ==",
	}

	mbcvClient := &mbcv.AuthenticatedClient{}

	mux := http.NewServeMux()

	handler := server.BasicIntentHandler{}

	server.CreateRoutes(mux, houndClient, mbcvClient, handler)

	s := &http.Server{
		Addr: ":80",
		Handler: mux,
	}

	log.Fatal(s.ListenAndServe())
}
