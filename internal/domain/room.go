package domain

import (
	"log"
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

func (r *Room) Run() {
	log.Printf("running chat room %v", r.UUID)
	for {
		select {
		case player := <-r.JoinChan:
			log.Printf("new player in room %v", r.UUID)
			r.Players[player] = true
		case player := <-r.LeaveChan:
			log.Printf("chatter player room %v", r.UUID)
			delete(r.Players, player)
			close(player.Send)
		case msg := <-r.MessageChan:
			for p := range r.Players {
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

func (room *Room) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		return
	}

	p := &Player{
		Conn: conn,
		Send: make(chan []byte, messageBufferSize),
		Room: room,
	}

	room.JoinChan <- p
	defer func() {
		room.LeaveChan <- p
	}()
	go p.Write()
	p.Read()
}
