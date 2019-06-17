package configuration

import (
	"encoding/json"
	"fmt"
	"os"
)

var (
	DefaultDBUrl   = "mongodb://127.0.0.1"
	Port           = ":5000"
	EventBroker    = "rabbitmq"
	EventBrokerURL = "amqp://guest:guest@localhost:5672"
	Exchange       = "events"
	Queue          = "event.created"
)

type Configuration struct {
	DatabaseURL    string `json:"databaseurl"`
	Port           string `json:"port"`
	EventBroker    string `json:"eventbroker"`
	EventBrokerURL string `json:"eventbrokerurl"`
	Exchange       string `json:"exchange"`
	Queue          string `json:"queue"`
}

//ReadConfig : Reads the configuration json and returns a json object
func ReadConfig(filename string) (Configuration, error) {
	conf := Configuration{
		DefaultDBUrl,
		Port,
		EventBroker,
		EventBrokerURL,
		Exchange,
		Queue,
	}

	file, err := os.Open(filename)
	if err != nil {
		fmt.Println("Config file not found, moving on with defaults")
		return conf, err
	}

	json.NewDecoder(file).Decode(&conf)
	return conf, err
}
