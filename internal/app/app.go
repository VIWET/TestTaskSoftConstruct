package app

import (
	"io/ioutil"

	"github.com/VIWET/TestTaskSoftConstruct/internal/server"
	"github.com/sirupsen/logrus"
	"gopkg.in/yaml.v3"
)

func Run(configPath string) error {
	config := server.NewConfig()

	configFile, err := ioutil.ReadFile(configPath)
	if err != nil {
		logrus.Fatal(err)
	}

	if err := yaml.Unmarshal(configFile, &config); err != nil {
		logrus.Fatal(err)
	}

	s := server.New(config)
	return s.Run()
}
