package server

import (
	"encoding/json"
	"net/http"

	"github.com/VIWET/TestTaskSoftConstruct/internal/domain"
	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
)

type server struct {
	router     *mux.Router
	logger     *logrus.Logger
	rooms      map[*domain.Room]bool
	games      []*domain.Game
	createChan chan *domain.Room
	deleteChan chan *domain.Room
}

func New() *server {
	return &server{
		router: mux.NewRouter(),
		logger: logrus.New(),
		rooms:  make(map[*domain.Room]bool),
		games: []*domain.Game{
			{
				UUID:  "1",
				Title: "1",
			},
			{
				UUID:  "2",
				Title: "2",
			},
		},
		createChan: make(chan *domain.Room),
		deleteChan: make(chan *domain.Room),
	}
}

func (s *server) ManageRooms() {
	for {
		select {
		case room := <-s.createChan:
			s.rooms[room] = true
			go room.Run(s.deleteChan)
		case room := <-s.deleteChan:
			delete(s.rooms, room)
		}
	}
}

func (s *server) Run() error {
	s.logger.Info("starting chat server on port 8080")
	s.router.HandleFunc("/", s.Index()).Methods("GET")
	s.router.HandleFunc("/room", s.CreateRoom()).Methods("POST")
	s.router.HandleFunc("/room/{uuid}", s.ConnectRoom())

	go s.ManageRooms()

	return http.ListenAndServe(":8080", s.router)
}

func (s *server) CreateRoom() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")

		var game domain.Game
		err := json.NewDecoder(r.Body).Decode(&game)
		if err != nil {
			return
		}

		room := domain.NewRoom(game)
		s.createChan <- room

		err = json.NewEncoder(w).Encode(room)
		if err != nil {
			return
		}
	}
}

func (s *server) Index() http.HandlerFunc {
	type Responce struct {
		Games []*domain.Game `json:"games"`
		Rooms []*domain.Room `json:"rooms"`
	}

	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		var rooms []*domain.Room
		for room := range s.rooms {
			rooms = append(rooms, room)
		}

		err := json.NewEncoder(w).Encode(Responce{
			Games: s.games,
			Rooms: rooms,
		})
		if err != nil {
			return
		}
	}
}

func (s *server) ConnectRoom() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		roomUUID, ok := mux.Vars(r)["uuid"]
		if !ok {
			return
		}

		var selectedRoom *domain.Room
		for room := range s.rooms {
			if room.UUID == roomUUID {
				selectedRoom = room
			}
		}

		if selectedRoom == nil {
			return
		}

		if len(selectedRoom.Players) == 4 {
			return
		}

		selectedRoom.ServeHTTP(w, r)
	}
}
