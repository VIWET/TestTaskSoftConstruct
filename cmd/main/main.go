package main

import (
	"flag"

	"github.com/VIWET/TestTaskSoftConstruct/internal/app"
)

var (
	configPath string
)

func init() {
	flag.StringVar(&configPath, "config-path", "config/config.yaml", "path to server config file")
}

func main() {
	app.Run(configPath)
}
