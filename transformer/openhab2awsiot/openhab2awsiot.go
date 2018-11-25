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

func (OpenHab2AWSIoT) Transform(from *transformer.Message) (*transformer.Message, error) {
	path := strings.Split(from.Topic, "/")
	if len(path) != 4 {
		return nil, TopicError.New("Invalid topic: %+v", path)
	}
	if path[1] != "openhab" {
		return nil, TopicError.New("Invalid first topic element: %+v", path)
	}
	if path[3] != "state" && path[3] != "command" {
		return nil, TopicError.New("Invalid third topic element: %+v", path)
	}
	to := &transformer.Message{
		Topic: fmt.Sprintf("/awsiot/%s", path[2]),
	}
	payload := make(map[string]interface{})
	payload[path[3]] = nil
	num, err := strconv.ParseInt(string(from.Payload), 10, 64)
	if err != nil {
		payload[path[3]] = string(from.Payload)
	} else {
		payload[path[3]] = num
	}
	to.Payload, err = json.Marshal(payload)
	if err != nil {
		return nil, ConvertError.New("Marshall error: %+v", err)
	}
	return to, nil
}
