package httpserver

import (
	"log"
	"net/http"
	"os"
	"time"
)

func Serve() {
	serv := http.Server{
		Addr:              ":" + os.Getenv("HTTP_PORT"),
		ReadTimeout:       10 * time.Second,
		ReadHeaderTimeout: 5 * time.Second,
		WriteTimeout:      10 * time.Second,
		IdleTimeout:       50 * time.Second,
		MaxHeaderBytes:    2 << 20,
	}
	http.HandleFunc("/info", getServicesInfo)
	http.HandleFunc("/add", addService)
	log.Fatal(serv.ListenAndServe())
}
