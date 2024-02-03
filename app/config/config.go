package config

import (
	"encoding/json"
	"os"
	"regexp"
	"rnd7/denon-mqtt/logger"
)

type MQTTConfig struct {
	URL    string `json:"url"`
	Retain bool   `json:"retain"`
	Topic  string `json:"topic"`
	QoS    int    `json:"qos"`
}

type Port struct {
	Name string `json:"name"`
	Port string `json:"port"`
}

type Config struct {
	MQTT  MQTTConfig `json:"mqtt"`
	Denon struct {
		IP string `json:"ip"`
	} `json:"denon"`
	LogLevel string `json:"loglevel,omitempty"`
}

func replaceEnvVariables(input []byte) []byte {
	envVariableRegex := regexp.MustCompile(`\${([^}]+)}`)

	return envVariableRegex.ReplaceAllFunc(input, func(match []byte) []byte {
		envVarName := match[2 : len(match)-1] // Extract the variable name without "${}".
		return []byte(os.Getenv(string(envVarName)))
	})
}

func LoadConfig(file string) (Config, error) {
	data, err := os.ReadFile(file)
	if err != nil {
		logger.Error("Error reading config file", err)
		return Config{}, err
	}

	data = replaceEnvVariables(data)

	// Create a Config object
	var config Config

	// Unmarshal the JSON data into the Config object
	err = json.Unmarshal(data, &config)
	if err != nil {
		logger.Error("Unmarshalling JSON:", err)
		return Config{}, err
	}

	if config.LogLevel == "" {
		config.LogLevel = "info"
	}

	return config, nil
}
