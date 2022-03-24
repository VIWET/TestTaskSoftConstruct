package server

import (
	"encoding/json"
	"net/http"
	"strconv"
	"time"

	"github.com/VIWET/TestTaskSoftConstruct/internal/domain"
	"github.com/gorilla/context"
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
		Games   []*domain.Game   `json:"games"`
		Rooms   []*domain.Room   `json:"rooms"`
		Players []*domain.Player `json:"players"`
	}

	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		var rooms []*domain.Room
		for room := range s.rooms {
			rooms = append(rooms, room)
		}

		players, err := s.playerRepository.GetAllPlayers()
		if err != nil {
			s.logger.Error(err)
			return
		}

		err = json.NewEncoder(w).Encode(Responce{
			Games:   s.games,
			Rooms:   rooms,
			Players: players,
		})
		if err != nil {
			s.logger.Error(err)
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

		id := context.Get(r, "userId")
		if id == nil {
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

		userId, err := strconv.Atoi(id.(string))
		if err != nil {
			return
		}

		p, err := s.playerRepository.GetPlayer(userId)
		if err != nil {
			return
		}

		err = s.playerRepository.SetInGameStatus(p.ID, 1)
		if err != nil {
			s.logger.Error(err)
		}
		defer s.playerRepository.SetInGameStatus(p.ID, 0)

		selectedRoom.ServeHTTP(w, r, p)
	}
}

func (s *server) Login() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		userID, ok := mux.Vars(r)["userId"]
		if !ok {
			return
		}

		id, err := strconv.Atoi(userID)
		if err != nil {
			s.logger.Error(err)
			return
		}

		p, err := s.playerRepository.GetPlayer(id)
		if err != nil {
			return
		}

		if p == nil {
			return
		}

		c := http.Cookie{
			Name:     "UserID",
			Value:    userID,
			Path:     "/",
			HttpOnly: true,
		}

		http.SetCookie(w, &c)
	}
}

func (s *server) Logout() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id := context.Get(r, "userId")
		if id == nil {
			return
		}

		c := http.Cookie{
			Name:     "UserID",
			Value:    "",
			Path:     "/",
			HttpOnly: true,
			Expires:  time.Unix(0, 0),
		}

		http.SetCookie(w, &c)
	}
}
