package mqtt

import (
	"github.com/philipparndt/mqtt-gateway/client"
	"github.com/philipparndt/mqtt-gateway/config"
)

func Start(config config.MQTTConfig, onMessage client.OnMessageListener) {
	client.Start(config, "denon")
	client.Subscribe(config.Topic+"/ports/+/poe/set", onMessage)
}
