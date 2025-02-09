package main

import (
	"log"
	messaging "pinger/amqp"
	"pinger/cmd"
	"pinger/iping"
	"pinger/logging"
)

func main() {
	logging.ConfigureLogger()
	log.Println("Starting...")
	var pinger iping.PingerInterface = &cmd.Pinger{}
	go messaging.ServeRabbit(pinger)
	pinger.DoPinging()
}
