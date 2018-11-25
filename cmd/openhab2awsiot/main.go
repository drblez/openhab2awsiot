package main

import (
	"github.com/Sirupsen/logrus"
	"go.uber.org/dig"
	"openhab2awsiot/config"
	"openhab2awsiot/logger"
	"openhab2awsiot/mqtt_service"
	"openhab2awsiot/transformer/openhab2awsiot"
	"os"
	"os/signal"
	"syscall"
)

func buildContainer() *dig.Container {
	c := dig.New()
	c.Provide(config.Init)
	c.Provide(logger.Init)
	c.Provide(mqtt_service.Init)
	c.Provide(openhab2awsiot.Init)
	return c
}

func do(conf *config.Config, log *logrus.Logger, mqtt *mqtt_service.MQTTService) {
	log.Debugf("Start")
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	err := mqtt.Start()
	if err != nil {
		log.Errorf("Subscribe error: %+v", err)
		return
	}
	<-c
}

func main() {
	container := buildContainer()
	err := container.Invoke(do)
	if err != nil {
		panic(err)
	}
}
