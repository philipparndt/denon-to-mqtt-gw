package main

import (
    "github.com/philipparndt/denon-to-mqtt-gw/config"
    "github.com/philipparndt/denon-to-mqtt-gw/denon"
    "github.com/philipparndt/denon-to-mqtt-gw/mqtt"
    "github.com/philipparndt/go-logger"
    "os"
    "os/signal"
    "syscall"
)

func main() {
    if len(os.Args) < 2 {
        logger.Error("No config file specified")
        os.Exit(1)
    }

    configFile := os.Args[1]
    logger.Info("Config file", configFile)
    cfg, err := config.LoadConfig(configFile)
    if err != nil {
        logger.Error("Failed loading config", err)
        return
    }

    logger.SetLevel(cfg.LogLevel)
    mqtt.Start(cfg.MQTT, func(s string, bytes []byte) {
        logger.Info("Received message", s, string(bytes))
    })

    denon.Start(cfg.Denon.IP + ":23")

    logger.Info("Application is now ready. Press Ctrl+C to quit.")

    quitChannel := make(chan os.Signal, 1)
    signal.Notify(quitChannel, syscall.SIGINT, syscall.SIGTERM)
    <-quitChannel

    logger.Info("Received quit signal")
}
