package main

import (
	"flag"

	"github.com/VIWET/TestTaskSoftConstruct/internal/app"
	"github.com/sirupsen/logrus"
)

var (
	configPath string
)

func init() {
	flag.StringVar(&configPath, "config-path", "config/config.yaml", "path to server config file")
}

func main() {
	if err := app.Run(configPath); err != nil {
		logrus.Error(err)
		return
	}
}
