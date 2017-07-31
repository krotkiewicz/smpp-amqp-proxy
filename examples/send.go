package main

import (
	"log"

	"github.com/streadway/amqp"
	"encoding/json"
)

func failOnError(err error, msg string) {
	if err != nil {
		log.Fatalf("%s: %s", msg, err)
	}
}

type Message struct {
	Content string `json:"content"`
	Src string `json:"src"`
	Dst string `json:"dst"`
}

func main() {
	conn, err := amqp.Dial("amqp://guest:guest@localhost:5672//")
	log.Println("a")
	failOnError(err, "Failed to connect to RabbitMQ")
	defer conn.Close()

	ch, err := conn.Channel()
	failOnError(err, "Failed to open a channel")
	defer ch.Close()

	q, err := ch.QueueDeclare(
		"in", // name
		true,   // durable
		false,   // delete when unused
		false,   // exclusive
		false,   // no-wait
		nil,     // arguments
	)
	failOnError(err, "Failed to declare a queue")

	for i := 0; i < 10; i++ {
		msg := Message{
			Content: "dupa",
			Src: "1",
			Dst: "2",
		}
		body, err := json.Marshal(msg)
		failOnError(err, "Failed to publish a message")
		err = ch.Publish(
			"",     // exchange
			q.Name, // routing key
			false,  // mandatory
			false,  // immediate
			amqp.Publishing{
				ContentType: "text/plain",
				Body:        []byte(body),
			})
		log.Printf(" [x] Sent %s", body)

	}
}
