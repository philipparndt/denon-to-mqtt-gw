package mqtt

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"os"
	"rnd7/denon-mqtt/config"
	"rnd7/denon-mqtt/logger"
	"sync"
	"time"

	PAHO "github.com/eclipse/paho.mqtt.golang"
)

var messagesPublishedCtr int

var client PAHO.Client
var baseTopic string

var connectionWg sync.WaitGroup

type OnMessageListener func(string, []byte)

func Start(config config.MQTTConfig, onMessage OnMessageListener) {
	connectionWg.Add(1)
	go connect(config, onMessage)
	connectionWg.Wait()
}

func generateRandomClientID(length int) string {
	charset := "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	seededRand := rand.New(rand.NewSource(time.Now().UnixNano()))
	result := make([]byte, length)
	for i := range result {
		result[i] = charset[seededRand.Intn(len(charset))]
	}
	return "denon_mqtt_" + string(result)
}

func LogMessagesPublished() {
	for {
		time.Sleep(time.Hour)
		logger.Info(fmt.Sprintf("Messages published (last hour): %d", messagesPublishedCtr))
		messagesPublishedCtr = 0
	}
}

func connect(config config.MQTTConfig, onMessage OnMessageListener) {
	baseTopic = config.Topic
	statusTopic := baseTopic + "/bridge/state"
	clientID := generateRandomClientID(10)
	logger.Debug("Generated client ID:", clientID)

	opts := PAHO.NewClientOptions().
		AddBroker(config.URL).
		SetClientID(clientID).
		SetWill(statusTopic, "offline", 1, true)

	client = PAHO.NewClient(opts)
	if token := client.Connect(); token.Wait() && token.Error() != nil {
		logger.Error("Error connecting to MQTT broker:", token.Error())
		os.Exit(1)
	}
	defer client.Disconnect(250)

	client.Subscribe(
		baseTopic+"/ports/+/poe/set",
		2,
		func(_ PAHO.Client, message PAHO.Message) {
			onMessage(message.Topic(), message.Payload())
		},
	)

	PublishAbsolute(statusTopic, "online", true)

	logger.Info("Connected to MQTT broker", config.URL)
	connectionWg.Done()
	go LogMessagesPublished()

	// Keep the connection active until the application is terminated
	select {}
}

func PublishAbsolute(topic string, message string, retained bool) {
	token := client.Publish(topic, 2, retained, message)
	token.Wait()

	messagesPublishedCtr++
	logger.Debug("Published message", topic, message)

	if token.Error() != nil {
		logger.Error("Error publishing message", token.Error())
	}
}

func PublishJSON(topic string, data any) {
	jsonData, err := json.Marshal(data)
	if err != nil {
		logger.Error("Error marshaling to JSON", err)
	} else {
		PublishAbsolute(baseTopic+"/"+topic, string(jsonData), false)
	}
}
