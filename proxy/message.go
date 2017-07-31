package proxy

import (
	"github.com/fiorix/go-smpp/smpp/pdu"
	"encoding/json"
)

type Message struct {
	Content string `json:"content"`
	Src string `json:"src"`
	Dst string `json:"dst"`
	Retries int `json:"retries,omitempty"`
}

func (msg *Message) fromPDU (pdu.Body) {
}

func (msg *Message) fromJSON (data []byte) (error) {
	err := json.Unmarshal(data, msg)
	if err != nil {
		return err
	}
	return nil
}

func (msg *Message) toJSON () ([]byte, error){
	data, err := json.Marshal(msg)
	if err != nil {
		return nil, err
	}
	return data, nil
}

