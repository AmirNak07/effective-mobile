package main

import (
	httpTrasport "effective-mobile/internal/http"
	"log"
	"net/http"
)

func main() {
	handler := httpTrasport.NewHandler()

	router := httpTrasport.NewRouter(handler)

	server := http.Server{
		Addr:    ":8080",
		Handler: router,
	}

	log.Fatal(server.ListenAndServe())
}
