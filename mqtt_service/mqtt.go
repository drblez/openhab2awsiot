package mqtt_service

import (
	"fmt"
	"github.com/Sirupsen/logrus"
	MQTT "github.com/eclipse/paho.mqtt.golang"
	"github.com/joomcode/errorx"
	"net"
	"openhab2awsiot/config"
	"openhab2awsiot/transformer"
)

var (
	Errors          = errorx.NewNamespace("mqtt_service")
	ConnectError    = Errors.NewType("connect_error")
	PublishError    = Errors.NewType("publish_error")
	SubscribeError  = Errors.NewType("subscribe_error")
	ParametersError = Errors.NewType("parameters_error")
)

type MQTTService struct {
	opts        *MQTT.ClientOptions
	client      MQTT.Client
	log         *logrus.Logger
	timeout     int
	topic       string
	transformer transformer.Transformer
}

func (mqtt *MQTTService) connect() error {
	if mqtt.client == nil {
		mqtt.client = MQTT.NewClient(mqtt.opts)
	}
	if mqtt.client.IsConnected() {
		return nil
	} else {
		if token := mqtt.client.Connect(); token.Wait() && token.Error() != nil {
			return ConnectError.WrapWithNoMessage(token.Error())
		}
		return nil
	}
}

func (mqtt *MQTTService) onConnectionLostHandler(client MQTT.Client, err error) {
	mqtt.log.Debugf("Connection to MQTT lost")
	if err := mqtt.connect(); err != nil {
		mqtt.log.Errorf("Reconnect error: %+v", err)
	}
}

func (mqtt *MQTTService) messageHandler(client MQTT.Client, message MQTT.Message) {
	mqtt.log.Debugf("Do: %s %s", message.Topic(), string(message.Payload()))
	msgFrom := &transformer.Message{
		Topic:   message.Topic(),
		Payload: message.Payload(),
	}
	msgTo, err := mqtt.transformer.Transform(msgFrom)
	if err != nil {
		mqtt.log.Errorf("Transformer error: %+v", err)
		return
	}
	err = mqtt.publish(msgTo.Topic, msgTo.Payload)
	if err != nil {
		mqtt.log.Errorf("Publish error: %+v", err)
		return
	}
}

func (mqtt *MQTTService) onConnectHandler(client MQTT.Client) {
	mqtt.log.Debugf("Connected to MQTT")
	if token := mqtt.client.Subscribe(mqtt.topic, 0, nil); token.Wait() && token.Error() != nil {
		mqtt.log.Errorf("Subscribe error: %+v", token.Error())
	}
}

func Init(config *config.Config, log *logrus.Logger, transformer transformer.Transformer) (*MQTTService, error) {
	mqtt := new(MQTTService)
	addr := net.JoinHostPort(config.MQTT.Host, config.MQTT.Port)
	mqtt.opts = MQTT.NewClientOptions().AddBroker(fmt.Sprintf("tcp://%s", addr))
	mqtt.timeout = config.MQTT.Timeout
	mqtt.log = log
	mqtt.topic = "#"
	if transformer == nil {
		return nil, ParametersError.New("transformer is nil")
	}
	mqtt.transformer = transformer
	mqtt.opts.OnConnectionLost = mqtt.onConnectionLostHandler
	mqtt.opts.OnConnect = mqtt.onConnectHandler
	mqtt.opts.DefaultPublishHandler = mqtt.messageHandler
	return mqtt, nil
}

func (mqtt *MQTTService) Start() error {
	if err := mqtt.connect(); err != nil {
		return err
	}
	return nil
}

func (mqtt *MQTTService) publish(topic string, payload []byte) error {
	if err := mqtt.connect(); err != nil {
		return err
	}
	if token := mqtt.client.Publish(topic, 0, false, payload); token.Wait() && token.Error() != nil {
		return PublishError.WrapWithNoMessage(token.Error())
	}
	return nil
}
