package server

import (
	"context"
	"database/sql"
	"fmt"
	"net/http"
	"time"

	"github.com/VIWET/TestTaskSoftConstruct/internal/domain"
	"github.com/VIWET/TestTaskSoftConstruct/internal/repository"
	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
)

type server struct {
	httpServer       *http.Server
	config           *Config
	router           *mux.Router
	logger           *logrus.Logger
	db               *sql.DB
	playerRepository repository.PlayerRepository
	gameRepository   repository.GameRepository
	rooms            map[*domain.Room]bool
	createChan       chan *domain.Room
	deleteChan       chan *domain.Room
}

func New(config *Config) *server {
	return &server{
		config:     config,
		router:     mux.NewRouter(),
		logger:     logrus.New(),
		rooms:      make(map[*domain.Room]bool),
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
	if err := s.configureLogger(); err != nil {
		return err
	}
	s.logger.Info(fmt.Sprintf("logger configured on level: %s", s.config.LogLevel))

	if err := s.configureDatabase(); err != nil {
		return err
	}
	s.logger.Info(fmt.Sprintf("database on %s:%s", s.config.DatabaseConfig.Host, s.config.DatabaseConfig.Port))

	s.playerRepository = repository.NewPlayerRepository(s.db)
	s.gameRepository = repository.NewGameRepository(s.db)

	s.setRoutes()

	s.configureServer()

	go s.ManageRooms()

	s.logger.Info(fmt.Sprintf("serving at http://localhost%s/", s.config.Addr))

	return s.httpServer.ListenAndServe()
}

func (s *server) configureServer() {
	s.httpServer = &http.Server{
		Addr:           s.config.Addr,
		Handler:        s.router,
		MaxHeaderBytes: 1 << 20,
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
	}
}

func (s *server) configureLogger() error {
	level, err := logrus.ParseLevel(s.config.LogLevel)
	if err != nil {
		return err
	}

	s.logger.SetLevel(level)

	return nil
}

func (s *server) configureDatabase() error {
	db, err := sql.Open("mysql", s.config.DatabaseConfig.GetConnectionString())
	if err != nil {
		s.logger.Fatal(err)
		return err
	}

	if err := db.Ping(); err != nil {
		s.logger.Fatal(err)
		return err
	}

	s.db = db

	return nil
}

func (s *server) setRoutes() {
	s.router.Handle("/", s.Index()).Methods("GET")
	s.router.Handle("/login/{userId}", s.Login()).Methods("GET")
	s.router.Handle("/logout", s.Middleware(s.Logout())).Methods("GET")
	s.router.Handle("/room", s.Middleware(s.CreateRoom())).Methods("POST")
	s.router.Handle("/room/{uuid}", s.Middleware(s.ConnectRoom()))
}

func (s *server) Shutdown(ctx context.Context) error {
	err := s.httpServer.Shutdown(ctx)
	if err != nil {
		return err
	}

	err = s.playerRepository.DropInGameStatus()
	if err != nil {
		return err
	}

	err = s.db.Close()
	if err != nil {
		return err
	}

	return nil
}
