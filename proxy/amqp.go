package proxy

import (
	"github.com/streadway/amqp"
	log "github.com/sirupsen/logrus"
	"sync"
	"time"
)

type AMQP struct {
	conn *amqp.Connection
	channel *amqp.Channel
	consumers map[string] chan *Message
	toAck map[*Message] *amqp.Delivery

	sync.Mutex
}

func (a *AMQP) loop() () {
	retry := 1

	for {
		conn, err := amqp.Dial("amqp://guest:guest@localhost:5672//")
		if err != nil {
			time.Sleep(time.Duration(retry) * time.Second)
			log.Printf("AMQP error %s, retying in %ds", err, retry)
			if retry < 32 {
				retry *= 2
			}
		} else {
			retry = 1
			log.Println("AMQP connected")
			a.conn = conn
			break
		}
	}

	connError := make(chan *amqp.Error)
	a.conn.NotifyClose(connError)

	go func () {
		<- connError
		a.conn.Close()
		a.conn = nil
		a.loop()
	}()

	a.channel = a.getChannel()

	for k := range a.consumers {
		a.Consume(k)
	}
}

func NewAMQP () (* AMQP){
	a := &AMQP{}
	a.consumers = make(map[string] chan *Message)
	a.toAck = make(map[*Message] *amqp.Delivery)
	a.loop()
	return a
}

func (a *AMQP) getChannel() (*amqp.Channel) {
	ch, err := a.conn.Channel()

	ch.QueueDeclare(
		"in",
		true,
		false,
		false,
		false,
		nil,
	)
	ch.QueueDeclare(
		"out",
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		return nil
	}

	err = ch.Qos(
		10,
		0,
		false,
	)
	if err != nil {
		return nil
	}
	return ch
}

func (a *AMQP) Consume(consumer string) (<-chan *Message) {
	consumerChannel := a.consumers[consumer]
	if consumerChannel == nil {
		consumerChannel = make(chan *Message)
		a.consumers[consumer] = consumerChannel
	}

	go func() {
		deliveries, err := a.channel.Consume(
			"in",
			consumer,
			false,
			false,
			false,
			true,
			nil,
		)
		if err != nil {
			return
		}

		for d := range deliveries {
			msg := Message{}
			log.Printf("Reciving message %s", d.Body)
			msg.fromJSON(d.Body)
			a.toAck[&msg] = &d
			consumerChannel <- &msg
		}
	}()
	return consumerChannel
}

func (a *AMQP) Ack(message *Message) () {
	delivery, ok := a.toAck[message]
	if ok {
		delete(a.toAck, message)
		delivery.Ack(false)
	}
}

func (a *AMQP) Push(message *Message) (error) {
	body, err := message.toJSON()
	if err != nil {
		return err
	}
	err = a.channel.Publish(
		"",
		"in",
		false,
		false,
		amqp.Publishing{
			ContentType: "text/plain",
			Body:        body,
		},
	)
	return err
}

