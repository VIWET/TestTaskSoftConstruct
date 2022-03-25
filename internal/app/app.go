package app

import (
	"context"
	"io/ioutil"
	"os"
	"os/signal"
	"syscall"

	"github.com/VIWET/TestTaskSoftConstruct/internal/server"
	"github.com/sirupsen/logrus"
	"gopkg.in/yaml.v3"
)

func Run(configPath string) {
	config := server.NewConfig()

	configFile, err := ioutil.ReadFile(configPath)
	if err != nil {
		logrus.Fatal(err)
	}

	if err := yaml.Unmarshal(configFile, &config); err != nil {
		logrus.Fatal(err)
	}

	s := server.New(config)

	go s.Run()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGTERM, syscall.SIGINT)
	<-quit

	if err := s.Shutdown(context.Background()); err != nil {
		logrus.Errorf("error occured on server shutting down: %s", err.Error())
	}
}
