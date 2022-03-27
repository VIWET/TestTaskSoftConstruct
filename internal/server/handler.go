package server

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"
	"time"

	"github.com/VIWET/TestTaskSoftConstruct/internal/domain"
	"github.com/gorilla/context"
	"github.com/gorilla/mux"
)

type errorResponse struct {
	StatusCode int    `json:"status_code"`
	Message    string `json:"message"`
}

func respond(w http.ResponseWriter, r *http.Request, code int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	if data != nil {
		json.NewEncoder(w).Encode(data)
	}
}

func errorRespond(w http.ResponseWriter, r *http.Request, code int, err error) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	json.NewEncoder(w).Encode(&errorResponse{
		StatusCode: code,
		Message:    err.Error(),
	})
}

func (s *server) CreateRoom() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")

		game := &domain.Game{}
		err := json.NewDecoder(r.Body).Decode(&game)
		if err != nil {
			s.logger.Error(err)
			errorRespond(w, r, http.StatusUnprocessableEntity, err)
			return
		}

		game, err = s.gameRepository.GetGame(game.ID)
		if err != nil {
			s.logger.Error(err)
			errorRespond(w, r, http.StatusBadRequest, err)
			return
		}

		room := domain.NewRoom(*game)
		s.createChan <- room

		respond(w, r, http.StatusOK, room)
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
			errorRespond(w, r, http.StatusInternalServerError, err)
			return
		}

		games, err := s.gameRepository.GetAllGames()
		if err != nil {
			s.logger.Error(err)
			errorRespond(w, r, http.StatusInternalServerError, err)
			return
		}

		respond(w, r, http.StatusOK, Responce{
			Games:   games,
			Rooms:   rooms,
			Players: players,
		})
	}
}

func (s *server) ConnectRoom() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		roomUUID, ok := mux.Vars(r)["uuid"]
		if !ok {
			errorRespond(w, r, http.StatusBadRequest, errors.New("no roomId"))
			return
		}

		id := context.Get(r, "userId")
		if id == nil {
			errorRespond(w, r, http.StatusInternalServerError, errors.New("error"))
			return
		}

		var selectedRoom *domain.Room
		for room := range s.rooms {
			if room.UUID == roomUUID {
				selectedRoom = room
			}
		}

		if selectedRoom == nil {
			errorRespond(w, r, http.StatusBadRequest, errors.New("error: there is no room "+roomUUID))
			return
		}

		if len(selectedRoom.Players) == 4 {
			respond(w, r, http.StatusOK, struct {
				status_code int
				message     string
			}{
				status_code: http.StatusOK,
				message:     "room is full",
			})
			return
		}

		userId, err := strconv.Atoi(id.(string))
		if err != nil {
			errorRespond(w, r, http.StatusInternalServerError, err)
			return
		}

		p, err := s.playerRepository.GetPlayer(userId)
		if err != nil {
			errorRespond(w, r, http.StatusInternalServerError, err)
			return
		}

		err = s.playerRepository.SetInGameStatus(p.ID, 1)
		if err != nil {
			s.logger.Error(err)
			errorRespond(w, r, http.StatusInternalServerError, err)
			return
		}
		defer s.playerRepository.SetInGameStatus(p.ID, 0)

		selectedRoom.ServeHTTP(w, r, p)
	}
}

func (s *server) Login() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		userID, ok := mux.Vars(r)["userId"]
		if !ok {
			errorRespond(w, r, http.StatusBadRequest, errors.New("no userId"))
			return
		}

		s.logger.Info("LOGIN", userID)

		id, err := strconv.Atoi(userID)
		if err != nil {
			s.logger.Error(err)
			errorRespond(w, r, http.StatusInternalServerError, err)
			return
		}

		p, err := s.playerRepository.GetPlayer(id)
		if err != nil {
			s.logger.Error(err)
			errorRespond(w, r, http.StatusInternalServerError, err)
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
		respond(w, r, http.StatusOK, userID)
	}
}

func (s *server) Logout() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id := context.Get(r, "userId")
		if id == nil {
			errorRespond(w, r, http.StatusBadRequest, errors.New("no userId"))
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
