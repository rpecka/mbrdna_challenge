package main

import (
	"../../pkg/server"
	"github.com/soundhound/houndify-sdk-go"
	"log"
	"net/http"
)

func main() {
	client := &houndify.Client{
		ClientID:          "aX1HeFqKSZoNJ0Wh-kZ6tA==",
		ClientKey:         "sCIOCsrBADASfXWEY_6-xiI4qm4VJhNf77LDH4mPqzCm4w3QZacJo4-mHZhNYGOTdbXadttX7TklNl2NGs9WhQ==",
	}

	mux := http.NewServeMux()

	server.CreateRoutes(mux, client)

	s := &http.Server{
		Addr: ":80",
		Handler: mux,
	}

	log.Fatal(s.ListenAndServe())
}
