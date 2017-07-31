package proxy

import (
	"github.com/fiorix/go-smpp/smpp"
	"github.com/fiorix/go-smpp/smpp/pdu/pdutext"
	"github.com/fiorix/go-smpp/smpp/pdu/pdufield"
	"go.uber.org/ratelimit"
	log "github.com/sirupsen/logrus"
	"context"
)

type RateLimiter struct {
	rl ratelimit.Limiter
}

func NewRateLimiter () RateLimiter {
	rl := RateLimiter{}
	rl.rl = ratelimit.New(1)
	return rl
}

func (rl RateLimiter) Wait (ctx context.Context) error {
	rl.rl.Take()
	return nil
}

type SMPP struct {
	tx *smpp.Transceiver
	ConnStatus smpp.ConnStatusID
}

func NewSMPP () (*SMPP) {
	s := &SMPP{}
	tx := &smpp.Transceiver{
		Addr:    "localhost:2775",
		User:    "smppclient1",
		Passwd:  "pwd1",
		SystemType: "smpp",
		RateLimiter: NewRateLimiter(),
	}
	connStatus := tx.Bind()
	s.ConnStatus = smpp.Disconnected
	go func() {
		for c := range connStatus {
			log.Println("SMPP connection status:", c.Status())
			s.ConnStatus = c.Status()
		}
	}()
	s.tx = tx
	return s
}

func (s *SMPP) Send(message *Message) (error) {
	_, err := s.tx.Submit(&smpp.ShortMessage{
		Src:      message.Src,
		Dst:      message.Dst,
		Text:     pdutext.Raw(message.Content),
		Register: pdufield.NoDeliveryReceipt,
		SourceAddrTON: 0,
		SourceAddrNPI: 1,
		DestAddrTON: 1,
		DestAddrNPI: 1,
	})
	if err != nil {
		return err
	}
	return nil
}

