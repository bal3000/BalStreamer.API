package configuration

import (
	"encoding/json"
	"fmt"
	"os"
)

// Configuration is the model to store the config data
type Configuration struct {
	ConnectionString string `json:"connectionString"`
}

// ReadConfig reads the json file and popluates the Configuration struct
func ReadConfig() Configuration {
	file, _ := os.Open("config.json")
	defer file.Close()

	decoder := json.NewDecoder(file)
	configuration := Configuration{}
	err := decoder.Decode(&configuration)

	if err != nil {
		fmt.Println(err)
	}
	return configuration
}
