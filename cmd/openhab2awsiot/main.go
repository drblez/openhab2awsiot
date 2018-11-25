package main

import (
	"github.com/Sirupsen/logrus"
	"go.uber.org/dig"
	"openhab2awsiot/config"
	"openhab2awsiot/logger"
)

func buildContainer() *dig.Container {
	c := dig.New()
	c.Provide(config.InitConfig)
	c.Provide(logger.InitLogger)
	return c
}

func do(conf *config.Config, log *logrus.Logger) {
	log.Debugf("Start")
}

func main() {
	container := buildContainer()
	err := container.Invoke(do)
	if err != nil {
		panic(err)
	}
}
