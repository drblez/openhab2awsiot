package openhab2awsiot

import (
	"encoding/json"
	"fmt"
	"github.com/joomcode/errorx"
	"openhab2awsiot/transformer"
	"strconv"
	"strings"
)

var (
	Error        = errorx.NewNamespace("openhab2awsiot_error")
	TopicError   = Error.NewType("topic_error")
	ConvertError = Error.NewType("convert_error")
)

type OpenHab2AWSIoT struct{}

func Init() transformer.Transformer {
	return &OpenHab2AWSIoT{}
}

type shadow struct {
	State struct {
		Reported struct {
			State   interface{} `json:"state,omitempty"`
			Command interface{} `json:"command,omitempty"`
		} `json:"reported"`
	} `json:"state"`
}

func (OpenHab2AWSIoT) Transform(from *transformer.Message) (*transformer.Message, error) {
	path := strings.Split(from.Topic, "/")
	if len(path) != 3 {
		return nil, TopicError.New("Invalid topic: %+v", path)
	}
	if path[0] != "openhab" {
		return nil, TopicError.New("Invalid first topic element: %+v", path)
	}
	if path[2] != "state" && path[2] != "command" {
		return nil, TopicError.New("Invalid third topic element: %+v", path)
	}
	to := &transformer.Message{
		Topic: fmt.Sprintf("$aws/things/%s/shadow/update", path[1]),
	}
	s := shadow{}
	num, err := strconv.ParseInt(string(from.Payload), 10, 64)
	if err != nil {
		if path[2] == "command" {
			s.State.Reported.Command = string(from.Payload)
		} else {
			s.State.Reported.State = string(from.Payload)
		}
	} else {
		if path[2] == "command" {
			s.State.Reported.Command = num
		} else {
			s.State.Reported.State = num
		}
	}
	to.Payload, err = json.Marshal(s)
	if err != nil {
		return nil, ConvertError.New("Marshall error: %+v", err)
	}
	return to, nil
}
