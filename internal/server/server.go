package server

import (
	"database/sql"
	"fmt"
	"net/http"

	"github.com/VIWET/TestTaskSoftConstruct/internal/domain"
	"github.com/VIWET/TestTaskSoftConstruct/internal/repository"
	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
)

type server struct {
	config           *Config
	router           *mux.Router
	logger           *logrus.Logger
	db               *sql.DB
	playerRepository repository.PlayerRepository
	rooms            map[*domain.Room]bool
	games            []*domain.Game
	createChan       chan *domain.Room
	deleteChan       chan *domain.Room
}

func New(config *Config) *server {
	return &server{
		config: config,
		router: mux.NewRouter(),
		logger: logrus.New(),
		rooms:  make(map[*domain.Room]bool),
		games: []*domain.Game{
			{
				ID:    1,
				Title: "1",
			},
			{
				ID:    2,
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
	if err := s.configureLogger(); err != nil {
		return err
	}
	s.logger.Info(fmt.Sprintf("logger configured on level: %s", s.config.LogLevel))

	if err := s.configureDatabase(); err != nil {
		return err
	}
	s.logger.Info(fmt.Sprintf("database on %s:%s", s.config.DatabaseConfig.Host, s.config.DatabaseConfig.Port))

	s.playerRepository = repository.NewPlayerRepository(s.db)

	s.router.HandleFunc("/", s.Index()).Methods("GET")
	s.router.HandleFunc("/room", s.CreateRoom()).Methods("POST")
	s.router.HandleFunc("/room/{uuid}", s.ConnectRoom())

	go s.ManageRooms()

	s.logger.Info(fmt.Sprintf("serving at http://localhost%s/", s.config.Addr))
	return http.ListenAndServe(s.config.Addr, s.router)
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
