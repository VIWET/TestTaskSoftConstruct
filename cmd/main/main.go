package main

import (
	"github.com/VIWET/TestTaskSoftConstruct/internal/app"
	"github.com/sirupsen/logrus"
)

func main() {
	if err := app.Run(); err != nil {
		logrus.Error(err)
		return
	}
}
