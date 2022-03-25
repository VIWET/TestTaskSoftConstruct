package domain

import (
	"net/http"

	"github.com/google/uuid"
	"github.com/gorilla/websocket"
)

const (
	socketBufferSize  = 1024
	messageBufferSize = 256
)

var upgrader = &websocket.Upgrader{
	ReadBufferSize:  socketBufferSize,
	WriteBufferSize: socketBufferSize,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

type Room struct {
	UUID        string           `json:"uuid"`
	Game        Game             `json:"game"`
	Players     map[*Player]bool `json:"-"`
	JoinChan    chan *Player     `json:"-"`
	LeaveChan   chan *Player     `json:"-"`
	MessageChan chan []byte      `json:"-"`
}

func NewRoom(game Game) *Room {
	return &Room{
		UUID:        uuid.NewString(),
		Game:        game,
		Players:     make(map[*Player]bool),
		JoinChan:    make(chan *Player),
		LeaveChan:   make(chan *Player),
		MessageChan: make(chan []byte),
	}
}

func (r *Room) Run(deleteChan chan *Room) {
Loop:
	for {
		select {
		case player := <-r.JoinChan:
			r.Players[player] = true
		case player := <-r.LeaveChan:
			delete(r.Players, player)
			close(player.Send)
			if len(r.Players) == 0 {
				break Loop
			}
		case msg := <-r.MessageChan:
			m, _ := Unmarshal(msg)
			for p := range r.Players {
				if m.Sender != p {
					select {
					case p.Send <- msg:
					default:
						delete(r.Players, p)
						close(p.Send)
					}
				}
			}
		}
	}
	deleteChan <- r
}

func (room *Room) ServeHTTP(w http.ResponseWriter, r *http.Request, p *Player) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		return
	}

	p.Conn = conn
	p.Room = room
	p.Send = make(chan []byte)

	room.JoinChan <- p
	defer func() {
		room.LeaveChan <- p
	}()
	go p.Write()
	p.Read()
}
