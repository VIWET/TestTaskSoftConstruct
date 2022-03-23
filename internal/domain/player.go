package domain

import "github.com/gorilla/websocket"

type Player struct {
	UUID string `json:"uuid"`
	Send chan []byte
	Room *Room
	Conn *websocket.Conn
}

func (p *Player) Read() {
	for {
		_, msg, err := p.Conn.ReadMessage()
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
		err := p.Conn.WriteMessage(websocket.TextMessage, msg)
		if err != nil {
			break
		}
	}
	p.Conn.Close()
}
