package main

import (
	"TrainerConnect/cmd/internal/user"
	"log"
	"net/http"
	"time"
)

func main() {
	router := http.NewServeMux()
	handler := user.NewHandler()
	handler.Register(router)
	start(router)

}

func start(router http.Handler) {
	server := &http.Server{
		Addr:         ":1234",
		Handler:      router,
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}

	log.Fatal(server.ListenAndServe())
}
