package server

import (
	"encoding/json"
	"net/http"

	"github.com/VIWET/TestTaskSoftConstruct/internal/domain"
	"github.com/gorilla/mux"
)

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
