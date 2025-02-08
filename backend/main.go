package main

import (
	"backend/httpserver"
	"backend/logging"
	"log"
)

func main() {
	logging.ConfigureLogger()
	log.Println("Starting...")
	httpserver.Serve()
}
