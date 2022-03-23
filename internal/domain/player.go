package domain

import (
	"github.com/gorilla/websocket"
)

type Player struct {
	UUID string          `json:"uuid"`
	Name string          `json:"name"`
	Send chan []byte     `json:"-"`
	Room *Room           `json:"-"`
	Conn *websocket.Conn `json:"-"`
}

func (p *Player) Read() {
	for {
		m := Message{}
		err := p.Conn.ReadJSON(&m)
		if err != nil {
			break
		}
		m.Sender = p
		msg, err := Marshal(&m)
		if err == nil {
			p.Room.MessageChan <- msg
		} else {
			break
		}
	}
	p.Conn.Close()
}

func (p *Player) Write() {
	for msg := range p.Send {
		m, err := Unmarshal(msg)
		if err != nil {
			break
		}
		err = p.Conn.WriteJSON(m)
		if err != nil {
			break
		}
	}
	p.Conn.Close()
}
