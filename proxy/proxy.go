package proxy

import (
	log "github.com/sirupsen/logrus"
	"github.com/fiorix/go-smpp/smpp"
	"time"
)

type Proxy struct {
	smpp *SMPP
	amqp *AMQP
}

func NewProxy () (*Proxy){
	return &Proxy{
		amqp : NewAMQP(),
		smpp: NewSMPP(),
	}
}

func (proxy *Proxy) Serve () () {
	msgs := proxy.amqp.Consume("aaa")
	for msg := range msgs {
		for proxy.smpp.ConnStatus != smpp.Connected {
			time.Sleep(time.Second)
		}
		err := proxy.smpp.Send(msg)

		if err == nil {
			log.Println("SMPP sent")
		} else {
			log.Println(err)
			proxy.amqp.Push(msg)
		}
		proxy.amqp.Ack(msg)
	}
}

func (proxy *Proxy) onSend(msg Message) (error) {
	return nil
}
