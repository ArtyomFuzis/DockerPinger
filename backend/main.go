package main

import (
	"backend/database"
	"backend/logging"
	"log"
	"time"
)

func main() {
	logging.ConfigureLogger()
	repo := database.GetPingRepository()
	log.Println("Successful start")
	for {
		time.Sleep(5 * time.Second)
		log.Println(repo.GetLastPing("addr").ID, repo.GetLastPing("addr").Date)
	}
}
