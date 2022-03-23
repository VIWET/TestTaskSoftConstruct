package server

import (
	"encoding/json"
	"net/http"

	"github.com/VIWET/TestTaskSoftConstruct/internal/domain"
	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
)

type server struct {
	router *mux.Router
	logger *logrus.Logger
	rooms  map[*domain.Room]bool
	games  []*domain.Game
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
	}
}

func (s *server) Run() error {
	s.logger.Info("starting chat server on port 8080")
	s.router.HandleFunc("/", s.Index()).Methods("GET")
	s.router.HandleFunc("/room", s.CreateRoom()).Methods("POST")
	s.router.Handle("/room/{uuid}", s.ConnectRoom())

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
		s.rooms[room] = true
		go room.Run()

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
		roomUUID := mux.Vars(r)["uuid"]

		var selectedRoom *domain.Room
		for room := range s.rooms {
			if room.UUID == roomUUID {
				selectedRoom = room
			}
		}

		selectedRoom.ServeHTTP(w, r)
	}
}
