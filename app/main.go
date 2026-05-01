package main

import (
	"expvar"
	"net/http"
	_ "net/http/pprof"
	"os"
	"os/signal"
	"syscall"

	"github.com/philipparndt/denon-to-mqtt-gw/config"
	"github.com/philipparndt/denon-to-mqtt-gw/denon"
	"github.com/philipparndt/denon-to-mqtt-gw/mqtt"
	"github.com/philipparndt/go-logger"
)

func main() {
	logger.Init("info", logger.Logger())
	initPprof()

	if len(os.Args) < 2 {
		logger.Error("No config file specified")
		os.Exit(1)
	}

	configFile := os.Args[1]
	logger.Info("Config file", "path", configFile)
	cfg, err := config.LoadConfig(configFile)
	if err != nil {
		logger.Error("Failed loading config", "error", err)
		return
	}

	logger.SetLevel(cfg.LogLevel)
	mqtt.Start(cfg.MQTT)

	denon.Start(cfg.Denon.IP + ":23")

	expvar.Publish("denonState", expvar.Func(func() any {
		return denon.GetState()
	}))

	logger.Info("Application is now ready. Press Ctrl+C to quit.")

	quitChannel := make(chan os.Signal, 1)
	signal.Notify(quitChannel, syscall.SIGINT, syscall.SIGTERM)
	<-quitChannel

	logger.Info("Received quit signal")
}

func initPprof() {
	go func() {
		http.ListenAndServe(":6060", nil)
	}()
}
