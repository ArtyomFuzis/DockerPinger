package messaging

import (
	"encoding/json"
	"fmt"
	amqp "github.com/rabbitmq/amqp091-go"
	"log"
	"os"
	"pinger/iping"
	"pinger/transfer"
	"strconv"
	"time"
)

func createConnectionChannel() *amqp.Channel {
	port, err := strconv.Atoi(os.Getenv("RABBITMQ_PORT"))
	if err != nil {
		port = 5672
	}
	url := fmt.Sprintf("amqp://%s:%s@%s:%d/",
		os.Getenv("RABBITMQ_USER"),
		os.Getenv("RABBITMQ_PASSWORD"),
		os.Getenv("RABBITMQ_HOST"),
		port,
	)
	tryCnt, err := strconv.Atoi(os.Getenv("RETRY_ATTEMPTS_RABBITMQ"))
	if err != nil {
		tryCnt = 10
	}
	var res *amqp.Connection
	for i := 1; i < tryCnt; i++ {
		con, err := amqp.Dial(url)
		if err != nil {
			retryTimes, err := strconv.Atoi(os.Getenv("RETRY_TIME_RABBITMQ"))
			if err != nil {
				retryTimes = 5
			}
			log.Printf("Retry connecting to Rabbit in %ds\n", retryTimes)
			time.Sleep(time.Duration(retryTimes) * time.Second)
		} else {
			res = con
			break
		}
	}
	if res == nil {
		log.Fatal("Attempts count exceeded. Failed to connect to database")
	}
	ch, err := res.Channel()
	if err != nil {
		log.Fatal("Failed to open channel:", err)
	}
	return ch
}
func createQueues() (*amqp.Queue, *amqp.Queue, error) {
	queue1Name := os.Getenv("RABBITMQ_QUEUE_ADD_SERVICE")
	if queue1Name == "" {
		queue1Name = "add_service"
	}
	q1, err := chn.QueueDeclare(
		queue1Name,
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		return nil, nil, err
	}
	queue2Name := os.Getenv("RABBITMQ_QUEUE_ADD_PING")
	if queue2Name == "" {
		queue2Name = "add_ping"
	}
	q2, err := chn.QueueDeclare(
		queue2Name,
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		return nil, nil, err
	}
	return &q1, &q2, err
}
func SendToAddPing(address string, date time.Time, state bool) {
	body, err := json.Marshal(transfer.MessageAddPing{Address: address, Date: date, State: state})
	if err != nil {
		log.Println("Unable to json encode message body:", err)
		return
	}
	err = chn.Publish(
		"",
		queueToAddPing.Name,
		false,
		false,
		amqp.Publishing{
			ContentType: "application/json",
			Body:        body,
		})
	if err != nil {
		log.Println("failed to publish a message:", err)
	}
}
func startConsumeAddService() {
	messages, err := chn.Consume(
		queueToAddService.Name,
		"",
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		log.Fatalln("failed to register a consumer:", err)
	}

	for message := range messages {
		var res transfer.MessageAddService
		err := json.Unmarshal(message.Body, &res)
		if err != nil {
			log.Println("Failed to unmarshal message:", err)
			continue
		}
		if res.Action == transfer.Add {
			pinger.AddService(res.Address)
		} else if res.Action == transfer.Delete {
			pinger.DeleteService(res.Address)
		}
	}
}

var queueToAddService *amqp.Queue
var queueToAddPing *amqp.Queue
var chn *amqp.Channel
var pinger iping.PingerInterface

func ServeRabbit(cmdPinger iping.PingerInterface) {
	pinger = cmdPinger
	chn = createConnectionChannel()
	defer func(con *amqp.Channel) {
		err := con.Close()
		if err != nil {
			log.Fatal("Failed to close channel:", err)
		}
	}(chn)
	var err error
	queueToAddService, queueToAddPing, err = createQueues()
	if err != nil {
		log.Fatal("Failed to create queues:", err)
	}
	startConsumeAddService()
}
