package rabbitmq

import (
	"encoding/json"
)

func MappingJsonToRabbitMQMessage(data interface{}) ([]byte, error) {
	jsonBytes, err := json.Marshal(data)
	if err != nil {
		return nil, err
	}

	return jsonBytes, nil
}
