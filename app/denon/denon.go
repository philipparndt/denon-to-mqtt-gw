package denon

import (
	"fmt"
	"github.com/philipparndt/denon-to-mqtt-gw/mqtt"
	"github.com/philipparndt/go-logger"
	"net"
	"os"
	"strconv"
	"strings"
	"sync"
)

var state Message = Message{
	Power:     "",
	Volume:    0,
	VolumeMax: 0,
}

const PWR = "PW"
const MAIN_VOLUME = "MV"
const MAIN_VOLUME_MAX = "MVMAX"

func RequestValue(command string) {
	_, err := connection.Write([]byte(command + "?\r"))
	if err != nil {
		logger.Error("Write to server failed:", err.Error())
		os.Exit(1)
	}
}

var connection *net.TCPConn

func initConnection(servAddr string) {
	tcpAddr, err := net.ResolveTCPAddr("tcp", servAddr)
	if err != nil {
		println("ResolveTCPAddr failed:", err.Error())
		os.Exit(1)
	}

	connection, err = net.DialTCP("tcp", nil, tcpAddr)
	if err != nil {
		println("Dial failed:", err.Error())
		os.Exit(1)
	}
	connectionWg.Done()
}

func parseStringToFloat(input string) (float64, error) {
	// Remove non-numeric characters from the string
	numericPart := strings.TrimLeftFunc(input, func(r rune) bool {
		return !('0' <= r && r <= '9' || r == '.')
	})

	// If there's a single digit after "MV", consider it as a decimal point
	if strings.HasPrefix(input, "MV") && len(numericPart) > 2 {
		numericPart = numericPart[0:2] + "." + numericPart[2:]
	}

	// Convert the remaining string to a float
	result, err := strconv.ParseFloat(numericPart, 32)
	if err != nil {
		return 0, err
	}

	return result, nil
}

var connectionWg sync.WaitGroup

func Start(servAddr string) {
	connectionWg.Add(1)
	go run(servAddr)
	connectionWg.Wait()
}

func run(servAddr string) {
	initConnection(servAddr)
	defer func(connection *net.TCPConn) {
		err := connection.Close()
		if err != nil {
			println("Close failed:", err.Error())
			os.Exit(1)
		}
	}(connection)

	for _, command := range []string{PWR, MAIN_VOLUME_MAX, MAIN_VOLUME} {
		RequestValue(command)
	}

	reply := make([]byte, 1024)
	for {
		_, err := connection.Read(reply)
		if err != nil {
			logger.Error("Read from server failed:", err.Error())
			os.Exit(1)
		}

		notify := false

		split := strings.Split(string(reply), "\r")
		for _, s := range split {
			trimmed := strings.TrimSpace(s)
			if trimmed == "PWON" {
				state.Power = "ON"
				logger.Info("Power is on")
				notify = true
			} else if trimmed == "PWSTANDBY" {
				state.Power = "STANDBY"
				logger.Info("Power is in standby")
				notify = true
			} else if strings.HasPrefix(trimmed, "MVMAX") {
				float, err := parseStringToFloat(trimmed)
				if err == nil {
					logger.Info(fmt.Sprintf("Max Volume is %f", float))
					state.VolumeMax = float
					notify = true
				} else {
					logger.Error("Error parsing volume:", err.Error())
				}
			} else if strings.HasPrefix(trimmed, "MV") {
				float, err := parseStringToFloat(trimmed)
				if err == nil {
					logger.Info(fmt.Sprintf("Volume is %f", float))
					state.Volume = float
					notify = true
				} else {
					logger.Error("Error parsing volume:", err.Error())
				}
			} else {
				// logger.Info("unhandled reply from server=", trimmed, len(trimmed))
			}
			//if trimmed != "" {
			//	logger.Info("reply from server=", trimmed, len(trimmed))
			//}
		}

		if notify {
			mqtt.PublishJSON("state", state)
		}
	}
}
