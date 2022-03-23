package server

import (
	"net/http"

	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
)

type server struct {
	router *mux.Router
	logger *logrus.Logger
}

func (s *server) Run() error {

	return http.ListenAndServe(":8080", s.router)
}
