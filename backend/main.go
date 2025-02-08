package main

import (
	messaging "backend/amqp"
	"backend/httpserver"
	"backend/logging"
	"log"
)

func main() {
	logging.ConfigureLogger()
	log.Println("Starting...")
	go httpserver.Serve()
	messaging.ServeRabbit()
}
