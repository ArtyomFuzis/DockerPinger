package database

import (
	"fmt"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"log"
	"os"
	"strconv"
	"sync"
	"time"
)

func connectPingDatabase() *gorm.DB {
	var res *gorm.DB
	port, err := strconv.Atoi(os.Getenv("PingDB_PORT"))
	if err != nil {
		port = 5432
	}
	dsn := fmt.Sprintf(
		"host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		os.Getenv("PingDB_HOST"),
		port,
		os.Getenv("PingDB_USER"),
		os.Getenv("PingDB_PASSWORD"),
		os.Getenv("PingDB_DATABASE"),
	)
	tryCnt, err := strconv.Atoi(os.Getenv("RETRY_ATTEMPTS_DB"))
	if err != nil {
		tryCnt = 10
	}
	for i := 1; i < tryCnt; i++ {
		con, err := gorm.Open(postgres.New(postgres.Config{DSN: dsn}))
		if err != nil {
			//log.Println("Database connection failed: ", err)
			retryTimes, err := strconv.Atoi(os.Getenv("RETRY_TIME_DB"))
			if err != nil {
				retryTimes = 5
			}
			log.Printf("Retry connecting to Database in %ds\n", retryTimes)
			time.Sleep(time.Duration(retryTimes) * time.Second)
		} else {
			res = con
			break
		}
	}
	if res == nil {
		log.Fatal("Attempts count exceeded. Failed to connect to database")
	}
	err = res.AutoMigrate(&Ping{})
	if err != nil {
		log.Fatal("Unable to auto migrate Ping table:", err)
	}
	err = res.AutoMigrate(&PingedServices{})
	if err != nil {
		log.Fatal("Unable to auto migrate PingedServices table:", err)
	}
	return res
}

var conn *gorm.DB
var mut sync.Mutex

func GetPingRepository() *PingRepository {
	mut.Lock()
	if conn == nil {
		conn = connectPingDatabase()
	}
	mut.Unlock()
	return &PingRepository{pingConnection: conn}
}
