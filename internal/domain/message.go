package domain

import (
	"encoding/json"
)

type Message struct {
	Sender  *Player `json:"sender"`
	Content string  `json:"content"`
}

func Unmarshal(msg []byte) (*Message, error) {
	m := &Message{}
	err := json.Unmarshal(msg, m)
	if err != nil {
		return nil, err
	}

	return m, nil
}

func Marshal(msg *Message) ([]byte, error) {
	m, err := json.Marshal(msg)
	if err != nil {
		return nil, err
	}

	return m, nil
}
