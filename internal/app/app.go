package app

import "github.com/VIWET/TestTaskSoftConstruct/internal/server"

func Run() error {
	s := server.New()
	return s.Run()
}
